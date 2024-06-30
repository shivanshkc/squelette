package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shivanshkc/squelette/pkg/logger"
)

// AccessLogger middleware handles access logging.
func (m *Middleware) AccessLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		// This will be used to calculate the total request execution time.
		start := time.Now()
		// Shorthand for the request.
		req := eCtx.Request()

		newCtx := logger.AddContextValue(req.Context(), "request_id", uuid.NewString())
		// Update the request with the new context.
		*req = *req.WithContext(newCtx)

		// Embedding the writer into the custom-writer to persist status-code for logging.
		cWriter := &responseWriterWithCode{ResponseWriter: eCtx.Response()}
		// Update the underlying response writer.
		eCtx.SetResponse(echo.NewResponse(cWriter, eCtx.Echo()))

		// Request entry log.
		slog.InfoContext(req.Context(), "request received", "url", req.URL.String())
		// Release control to the next middleware or handler.
		err := next(eCtx)
		// Request exit log.
		slog.InfoContext(req.Context(), "request completed", "latency", time.Since(start))

		return err
	}
}

// responseWriterWithCode is a wrapper for http.ResponseWriter for persisting statusCode.
type responseWriterWithCode struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriterWithCode) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
