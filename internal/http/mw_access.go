package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/shivanshkc/squelette/pkg/utils/ctxutils"
)

// AccessLogger middleware handles access logging.
func (m *Middleware) AccessLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		// This will be used to calculate the total request execution time.
		start := time.Now()
		// Shorthand for the underlying request.
		req := eCtx.Request()

		// Setup the request's context.
		ctxutils.SetRequestCtxInfo(req)
		// Fetch the logger for the updated request context.
		log := m.Logger.WithContext(req.Context())

		// Embedding the writer into the custom-writer to persist status-code for logging.
		cWriter := &responseWriterWithCode{ResponseWriter: eCtx.Response()}
		// Update the underlying response writer.
		eCtx.SetResponse(echo.NewResponse(cWriter, eCtx.Echo()))

		// Request entry log.
		log.Info().Str("method", req.Method).Str("url", req.URL.String()).
			Msg("request received")

		// Release control to the next middleware or handler.
		err := next(eCtx)

		// Request exit log.
		log.Info().Int("code", cWriter.statusCode).Int64("latency", int64(time.Since(start))).
			Msg("request completed")

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
