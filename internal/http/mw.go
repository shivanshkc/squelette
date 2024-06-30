package http

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Middleware implements all the REST middleware methods.
type Middleware struct{}

// Recovery is a panic recovery middleware.
func (m *Middleware) Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           func(c echo.Context) bool { return false },
		StackSize:         middleware.DefaultRecoverConfig.StackSize,
		DisableStackAll:   false,
		DisablePrintStack: false,
		// This allows the usage of our custom logger.
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			slog.ErrorContext(c.Request().Context(), "", "stack", stack)
			return err
		},
	})(next)
}

// CORS is a Cross-Origin Resource Sharing (CORS) middleware.
func (m *Middleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          func(c echo.Context) bool { return false },
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"*"},
	})(next)
}

// Secure defends against cross-site scripting (XSS) attack, content type sniffing, clickjacking,
// insecure connection and other code injection attacks.
func (m *Middleware) Secure(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Secure()(next)
}
