package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/shivanshkc/squelette/internal/config"
	"github.com/shivanshkc/squelette/internal/utils/errutils"
	"github.com/shivanshkc/squelette/internal/utils/httputils"
	"github.com/shivanshkc/squelette/internal/utils/signals"
)

// Server is the HTTP server of this application.
type Server struct {
	Config     config.Config
	Middleware Middleware
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
	}).Methods(http.MethodGet)

	// More API routes here...

	// Handle 404.
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httputils.WriteErr(w, errutils.NotFound())
	})

	return router
}
