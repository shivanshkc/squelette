package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/shivanshkc/template-microservice-go/pkg/config"
	"github.com/shivanshkc/template-microservice-go/pkg/logger"
	"github.com/shivanshkc/template-microservice-go/pkg/utils/errutils"
	"github.com/shivanshkc/template-microservice-go/pkg/utils/signals"
)

// Server is the HTTP server of this application.
type Server struct {
	Config     *config.Config
	Logger     *logger.Logger
	Middleware *Middleware

	echoInstance *echo.Echo
}

// Start sets up all the dependencies and routes on the server, and calls ListenAndServe on it.
func (s *Server) Start() {
	// Create echo instance.
	s.echoInstance = echo.New()
	s.echoInstance.HideBanner = true            // No banner.
	s.echoInstance.Logger.SetOutput(io.Discard) // No internal logging.
	// Add a custom HTTP error handler to the echo instance.
	s.echoInstance.HTTPErrorHandler = func(err error, eCtx echo.Context) {
		var errHTTP *errutils.HTTPError

		// Determine the type of error by checking with echo's builtin errors.
		if errors.Is(err, echo.ErrNotFound) {
			errHTTP = errutils.NotFound().WithReasonErr(err)
		} else {
			errHTTP = errutils.ToHTTPError(err)
		}

		_ = eCtx.JSON(errHTTP.StatusCode, errHTTP)
	}

	// Register the REST methods.
	s.registerRoutes()

	// Create the HTTP server.
	server := &http.Server{
		Addr:              s.Config.HTTPServer.Addr,
		ReadHeaderTimeout: time.Minute,
	}

	// Attach this http server to echo.
	// This is required, otherwise echoInstance.Close will not close the server.
	s.echoInstance.Server = server

	// Gracefully shut down upon interruption.
	signals.OnSignal(func(_ os.Signal) {
		s.Logger.Info().Msg("interruption detected, gracefully shutting down the server")

		// Graceful shutdown.
		if err := server.Shutdown(context.Background()); err != nil {
			s.Logger.Error().Err(fmt.Errorf("failed to gracefully shutdown the server: %w", err)).Send()
		}
	})

	s.Logger.Info().Msg(fmt.Sprintf("http server running at: %s", s.Config.HTTPServer.Addr))
	// Start the HTTP server.
	if err := s.echoInstance.StartServer(server); !errors.Is(err, http.ErrServerClosed) {
		s.Logger.Fatal().Err(fmt.Errorf("error in echoInstance.StartServer call: %w", err)).Send()
	}
}

// registerRoutes attaches middleware and REST methods to the server.
func (s *Server) registerRoutes() {
	// Setup global middleware.
	s.echoInstance.Use(s.Middleware.Recovery)     // For panic recovery.
	s.echoInstance.Use(s.Middleware.CORS)         // For CORS.
	s.echoInstance.Use(s.Middleware.Secure)       // Protection against XSS attack, content type sniffing etc
	s.echoInstance.Use(s.Middleware.AccessLogger) // For access logging.

	// Sample REST method.
	s.echoInstance.GET("/api", func(c echo.Context) error {
		s.Logger.ForContext(c.Request().Context()).Info().Msg("sample log statement")
		return c.JSON(http.StatusOK, map[string]interface{}{"code": "OK"}) //nolint:wrapcheck
	})

	// More methods can be defined here.
}
