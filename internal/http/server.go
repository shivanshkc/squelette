package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/shivanshkc/squelette/pkg/config"
	"github.com/shivanshkc/squelette/pkg/utils/errutils"
	"github.com/shivanshkc/squelette/pkg/utils/signals"
)

// Server is the HTTP server of this application.
type Server struct {
	Config     *config.Config
	Middleware *Middleware

	echoInst *echo.Echo
}

// Start sets up all the dependencies and routes on the server, and calls ListenAndServe on it.
func (s *Server) Start() {
	// Create echo instance.
	s.echoInst = echo.New()
	s.echoInst.HideBanner = true            // No banner.
	s.echoInst.Logger.SetOutput(io.Discard) // No internal logging.

	// Add a custom HTTP error handler to the echo instance.
	s.echoInst.HTTPErrorHandler = s.errorHandler
	// Register the REST methods.
	s.registerRoutes()

	// Create the HTTP server.
	server := &http.Server{
		Addr:              s.Config.HTTPServer.Addr,
		ReadHeaderTimeout: time.Minute,
	}

	// Attach this http server to echo.
	// This is required for methods like echoInst.Close to work.
	s.echoInst.Server = server

	// Gracefully shut down upon interruption.
	signals.OnSignal(func(_ os.Signal) {
		slog.Info("interruption detected, gracefully shutting down the server")
		// Graceful shutdown.
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error("failed to gracefully shutdown the server", "err", err)
		}
	})

	slog.Info("http server started", "name", s.Config.Application.Name, "addr", s.Config.HTTPServer.Addr)
	// Start the HTTP server.
	if err := s.echoInst.StartServer(server); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("error in echoInstance.StartServer call", "err", err)
		panic(err)
	}
}

// registerRoutes attaches middleware and REST methods to the server.
func (s *Server) registerRoutes() {
	// Setup global middleware.
	s.echoInst.Use(s.Middleware.Recovery)     // For panic recovery.
	s.echoInst.Use(s.Middleware.CORS)         // For CORS.
	s.echoInst.Use(s.Middleware.Secure)       // Protection against XSS attack, content type sniffing etc
	s.echoInst.Use(s.Middleware.AccessLogger) // For access logging.

	// Sample REST method.
	s.echoInst.GET("/api", func(c echo.Context) error {
		slog.InfoContext(c.Request().Context(), "example log")
		return c.JSON(http.StatusOK, map[string]any{"code": "OK"})
	})

	fmt.Println(s.Config)
	// Enable profiling if configured.
	if s.Config.Application.PProf {
		s.enableProfiling()
	}
}

// enableProfiling enables profiling and registers pprof REST endpoints.
func (s *Server) enableProfiling() {
	// Enable block profiling.
	runtime.SetBlockProfileRate(1)
	// Enable mutex profiling.
	runtime.SetMutexProfileFraction(1)

	// Create and setup the multiplexer.
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Attach the multiplexer to echo.
	s.echoInst.GET("/debug/pprof/*", echo.WrapHandler(mux))
	slog.Info("pprof endpoints available at: /debug/pprof")
}

// errorHandler handles all echo HTTP errors.
func (s *Server) errorHandler(err error, eCtx echo.Context) {
	ctx := eCtx.Request().Context()
	// Convert to HTTP error to send back the response.
	errHTTP := errutils.ToHTTPError(err)

	// Log HTTP errors.
	switch errHTTP.StatusCode / 100 {
	case 4: //nolint:gomnd // Represents 4xx behaviour.
		slog.InfoContext(ctx, "bad request", err, errHTTP)
	case 5: //nolint:gomnd // Represents 5xx behaviour.
		slog.ErrorContext(ctx, "internal error", err, errHTTP)
	default:
		slog.ErrorContext(ctx, "unexpected error", err, errHTTP)
	}

	// Response.
	_ = eCtx.JSON(errHTTP.StatusCode, errHTTP)
}
