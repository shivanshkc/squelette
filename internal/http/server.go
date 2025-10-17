package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/shivanshkc/squelette/internal/handlers"
	"github.com/shivanshkc/squelette/internal/middleware"
	"github.com/shivanshkc/squelette/internal/utils/httputils"
)

// Server is the HTTP server of this application.
type Server struct {
	addr       string
	httpServer *http.Server
}

// NewServer returns a new Server instance. Use the Start method to start the server.
func NewServer(addr string, handler *handlers.Handler) (*Server, error) {
	httpHandler := func() http.Handler {
		mux := http.NewServeMux()
		mw := middleware.Middleware{}

		// Health check.
		mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
			httputils.Write(w, http.StatusOK, nil, map[string]any{"code": "OK"})
		})

		return mw.CORS(mw.AccessLogger(mw.Recovery(mux)))
	}()

	httpServer := &http.Server{Addr: addr, ReadHeaderTimeout: time.Minute, Handler: httpHandler}
	return &Server{addr: addr, httpServer: httpServer}, nil
}

// Start sets up all the dependencies and routes on the server, and calls ListenAndServe on it.
func (s *Server) Start(ctx context.Context) error {
	// Create the HTTP server.

	// Channel to notify this thread of shut down completion.
	shutdownCompleteChan := make(chan struct{})

	// Cleanup goroutine.
	go func() {
		<-ctx.Done()

		// Shutdown server when the context expires.
		if err := s.httpServer.Shutdown(ctx); err != nil {
			slog.Error("error while shutting down http server: " + err.Error())
		}
		slog.Info("http server shut down complete")

		// Notify the main thread.
		shutdownCompleteChan <- struct{}{}
		close(shutdownCompleteChan)
	}()

	// Blocking call. Will release upon shut down (which happens upon context expiry).
	slog.Info("starting http server", "addr", s.addr)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error in ListenAndServe call: %w", err)
	}

	<-shutdownCompleteChan
	return nil
}
