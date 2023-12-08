package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shivanshkc/squelette/pkg/logger"
)

// Middleware implements all the REST middleware methods.
type Middleware struct {
	Logger *logger.Logger
}

// Recovery is a panic recovery middleware.
func (m *Middleware) Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           func(c echo.Context) bool { return false },
		StackSize:         middleware.DefaultRecoverConfig.StackSize,
		DisableStackAll:   false,
		DisablePrintStack: false,
		// This allows the usage of our custom logger.
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log := m.Logger.WithContext(c.Request().Context())
			log.Error().Err(err).Bytes("stack", stack).Msg("")
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
