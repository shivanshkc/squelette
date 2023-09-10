package http

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/shivanshkc/squelette/pkg/utils/ctxutils"
)

//nolint:funlen // This function may get shorter when we write more tests in the future that demand re-usability.
func TestAccessLogger(t *testing.T) {
	t.Parallel()

	// Create a middleware instance with a mock logger.
	writer := &bytes.Buffer{}
	mockMW := middlewareWithMockLogger(writer)

	// Data to verify against.
	expectedTraceID := "example-trace-id"
	expectedResponseStatus := http.StatusBadRequest
	expectedLogCount := 2

	// Mock HTTP request and response-writer.
	req, res := httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder()
	req.Header.Set("x-trace-id", expectedTraceID)

	// Echo instance that uses the mock request and response writer.
	echoInstance := echo.New()
	eCtx := echoInstance.NewContext(req, res)

	// Create an instance of access-logger middleware that passes control to a mock handler.
	accessLoggerMW := mockMW.AccessLogger(func(c echo.Context) error {
		return c.NoContent(expectedResponseStatus) //nolint:wrapcheck
	})

	// Expect no error.
	if err := accessLoggerMW(eCtx); err != nil {
		t.Errorf("expected no error but got: %+v", err)
		return
	}

	// Fetch context-info to verify if it was set correctly by the middleware.
	ctxInfo := ctxutils.GetRequestCtxInfo(req.Context())

	// Trace ID must be what's present in the request headers.
	if ctxInfo.TraceID != expectedTraceID {
		t.Errorf("expected trace ID to be: %s but got: %s", expectedTraceID, ctxInfo.TraceID)
		return
	}

	// Request ID must be initialized.
	if ctxInfo.RequestID == "" {
		t.Errorf("expected request ID to be initialized but got empty")
		return
	}

	// Expect the correct response code.
	if res.Code != expectedResponseStatus {
		t.Errorf("expected response code to be %d but got: %d", http.StatusOK, res.Code)
		return
	}

	// Count number of log statements printed.
	var actualLogCount int
	for scanner := bufio.NewScanner(writer); scanner.Scan(); {
		actualLogCount++
	}

	// Expect the correct number of log statements.
	if actualLogCount != expectedLogCount {
		t.Errorf("expected log count to be %d but got: %d", expectedLogCount, actualLogCount)
		return
	}
}
