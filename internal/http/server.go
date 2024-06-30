package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/shivanshkc/squelette/pkg/config"
	"github.com/shivanshkc/squelette/pkg/utils/httputils"
	"github.com/shivanshkc/squelette/pkg/utils/signals"
)

// Server is the HTTP server of this application.
type Server struct {
	Config     *config.Config
	Middleware *Middleware
	httpServer *http.Server
}

// Start sets up all the dependencies and routes on the server, and calls ListenAndServe on it.
func (s *Server) Start() {
	// All routes will be attached to this multiplexer.
	mux := http.NewServeMux()
	// Register the REST methods.
	s.registerRoutes(mux)

	// Create the HTTP server.
	s.httpServer = &http.Server{
		Addr:              s.Config.HTTPServer.Addr,
		ReadHeaderTimeout: time.Minute,
		Handler:           mux,
	}

	// Gracefully shut down upon interruption.
	signals.OnSignal(func(_ os.Signal) {
		slog.Info("interruption detected, gracefully shutting down the server")
		// Graceful shutdown.
		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			slog.Error("failed to gracefully shutdown the server", "err", err)
		}
	})

	slog.Info("starting http server", "name", s.Config.Application.Name, "addr", s.Config.HTTPServer.Addr)
	// Start the HTTP server.
	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("error in ListenAndServe call", "err", err)
		panic(err)
	}
}

// registerRoutes attaches middleware and REST methods to the server.
func (s *Server) registerRoutes(mux *http.ServeMux) {
	// The child mux will allow us to add global middlewares to all routes together.
	childMux := http.NewServeMux()

	// Sample REST method.
	childMux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "example log")
		httputils.Write(w, http.StatusOK, nil, map[string]any{"code": "OK"})
	})

	// More API routes here...

	// Enable profiling if configured.
	if s.Config.Application.PProf {
		s.enableProfiling(childMux)
	}

	// Attach global middleware.
	mux.Handle("/", s.Middleware.Recovery(s.Middleware.CORS(s.Middleware.AccessLogger(
		childMux.ServeHTTP,
	))))
}

// enableProfiling enables profiling and registers pprof REST endpoints.
func (s *Server) enableProfiling(mux *http.ServeMux) {
	// Enable block profiling.
	runtime.SetBlockProfileRate(1)
	// Enable mutex profiling.
	runtime.SetMutexProfileFraction(1)

	// Create and setup the multiplexer.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	slog.Info("pprof endpoints available at: /debug/pprof")
}
