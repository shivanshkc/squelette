package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/shivanshkc/squelette/pkg/utils/httputils"
)

// Middleware implements all the REST middleware methods.
type Middleware struct{}

func (m *Middleware) Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer recoverer(w, r)
		// Next middleware or handler.
		next.ServeHTTP(w, r)
	}
}

// recoverer is supposed to be called with a defer statement to recover a panic.
func recoverer(w http.ResponseWriter, r *http.Request) {
	// Recover the panic.
	errAny := recover()
	if errAny == nil {
		return
	}

	slog.ErrorContext(r.Context(), "panic occurred during request execution: %v", errAny)
	// Convert to error for handling.
	err, ok := errAny.(error)
	if !ok {
		err = fmt.Errorf("unexpected error: %v", errAny)
	}

	// Response.
	httputils.WriteErr(w, err)
}

// CORS middlewares handled the CORS issues.
func (m *Middleware) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
