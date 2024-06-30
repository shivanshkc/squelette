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

	"github.com/gorilla/mux"

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
	// Create the HTTP server.
	s.httpServer = &http.Server{
		Addr:              s.Config.HTTPServer.Addr,
		ReadHeaderTimeout: time.Minute,
		Handler:           s.getHandler(),
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
func (s *Server) getHandler() http.Handler {
	router := mux.NewRouter()

	// Attach middleware.
	router.Use(s.Middleware.Recovery)
	router.Use(s.Middleware.CORS)
	router.Use(s.Middleware.AccessLogger)

	// Sample REST method.
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "hello world")
		httputils.Write(w, http.StatusOK, nil, map[string]any{"code": "OK"})
	})

	// More API routes here...

	// Enable profiling if configured.
	if s.Config.Application.PProf {
		router.Handle("/", s.getProfilingHandler())
	}

	return router
}

// getProfilingHandler returns the handler that serves profiling data.
func (s *Server) getProfilingHandler() http.Handler {
	// Enable block profiling.
	runtime.SetBlockProfileRate(1)
	// Enable mutex profiling.
	runtime.SetMutexProfileFraction(1)

	router := mux.NewRouter()
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	slog.Info("pprof endpoints available at: /debug/pprof")
	return router
}
