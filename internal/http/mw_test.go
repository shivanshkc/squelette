package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shivanshkc/squelette/pkg/config"
	"github.com/shivanshkc/squelette/pkg/logger"
	"github.com/shivanshkc/squelette/pkg/utils/errutils"
)

func TestMiddleware_Recovery(t *testing.T) {
	t.Parallel()

	// Create a middleware instance with a mock logger.
	writer := &bytes.Buffer{}
	mockMW := middlewareWithMockLogger(writer)

	// This error should be present in the response.
	expectedResponse := errutils.InternalServerError().WithReasonStr("test panic")
	// Create a mock echo context for HTTP request execution.
	recorder := httptest.NewRecorder()
	mockEchoCtx := mockEchoContext(nil, recorder)

	// Create an instance of the recovery MW that passes control to a mock handler.
	recoveryMW := mockMW.Recovery(func(c echo.Context) error { panic(expectedResponse.Error()) })

	// Expect no error.
	if err := recoveryMW(mockEchoCtx); err != nil {
		t.Errorf("expected no error but got: %v", err)
		return
	}

	// Expect the correct status code.
	if recorder.Code != expectedResponse.StatusCode {
		t.Errorf("expected status code: %d but got: %d", expectedResponse.StatusCode, recorder.Code)
		return
	}

	// Decode the response body for verification.
	responseBody := map[string]any{}
	if err := json.NewDecoder(recorder.Body).Decode(&responseBody); err != nil {
		t.Errorf("failed to decode response body: %v", err)
		return
	}

	// Expect the correct status.
	if responseBody["status"] != expectedResponse.Status {
		t.Errorf("expected status to be: %s but got: %s", expectedResponse.Status, responseBody["status"])
		return
	}

	// Expect the correct reason.
	if responseBody["reason"] != expectedResponse.Reason {
		t.Errorf("expected reason to be: %s but got: %s", expectedResponse.Reason, responseBody["reason"])
		return
	}
}

func TestMiddleware_CORS(t *testing.T) {
	t.Parallel()

	// Create a middleware instance with a mock logger.
	mockMW := middlewareWithMockLogger(io.Discard)

	// HTTP request with the OPTIONS verb to test CORS.
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Add("origin", "example-origin")
	// Recorder for verifying the response.
	rec := httptest.NewRecorder()

	// Create a mock echo context for HTTP request execution.
	mockEchoCtx := mockEchoContext(req, rec)

	// Create an instance of the CORS MW that passes control to a mock handler.
	corsMW := mockMW.CORS(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Expect no error.
	if err := corsMW(mockEchoCtx); err != nil {
		t.Errorf("expected no error but got: %v", err)
		return
	}

	// Check the value of the allow-origin header.
	allowOriginHeader := rec.Header().Get(echo.HeaderAccessControlAllowOrigin)
	if allowOriginHeader != "*" {
		t.Errorf("expected the allow origin header to be * but got: %s", allowOriginHeader)
		return
	}
}

func TestMiddleware_Secure(t *testing.T) {
	t.Parallel()

	// Create a middleware instance with a mock logger.
	mockMW := middlewareWithMockLogger(io.Discard)

	// Recorder for verifying the response.
	rec := httptest.NewRecorder()
	// Create a mock echo context for HTTP request execution.
	mockEchoCtx := mockEchoContext(nil, rec)

	// Create an instance of the CORS MW that passes control to a mock handler.
	secureMW := mockMW.Secure(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Expect no error.
	if err := secureMW(mockEchoCtx); err != nil {
		t.Errorf("expected no error but got: %v", err)
		return
	}

	xssActual, xssExpected := rec.Header().Get(echo.HeaderXXSSProtection), middleware.DefaultSecureConfig.XSSProtection
	// Check if the XSS-Protection header is set.
	if xssActual != xssExpected {
		t.Errorf("expected %s header value to be %s but got: %s",
			echo.HeaderXXSSProtection, xssExpected, xssActual)
		return
	}
}

// middlewareWithMockLogger returns a middleware instance that uses a mock logger.
// The provided writer is used as the underlying io.Writer for the logger instance.
func middlewareWithMockLogger(writer io.Writer) *Middleware {
	// Setup basic dependencies.
	conf := config.LoadMock()
	logger.Init(writer, conf.Logger.Level, conf.Logger.Pretty)
	return &Middleware{}
}

// mockEchoContext returns an echo context that uses a mock http request and response instance.
// It also returns the *httptest.ResponseRecorder type as it cannot be obtained from the echo context.
func mockEchoContext(req *http.Request, rec *httptest.ResponseRecorder) echo.Context {
	// Use a simple request if provided is nil.
	if req == nil {
		req = httptest.NewRequest(http.MethodGet, "/", nil)
	}
	// Use a simple recorder if provided is nil.
	if rec == nil {
		rec = httptest.NewRecorder()
	}

	// Echo instance that uses the mock request and response writer.
	echoInstance := echo.New()
	// Add a custom HTTP error handler to the echo instance.
	echoInstance.HTTPErrorHandler = func(err error, c echo.Context) {
		errHTTP := errutils.ToHTTPError(err)
		_ = c.JSON(errHTTP.StatusCode, errHTTP)
	}
	// This context is required by the middleware.
	return echoInstance.NewContext(req, rec)
}
