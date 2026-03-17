# Squelette

## Introduction

Squelette (French for "Skeleton") is a Go web service template that provides a clean, production-ready foundation for
building HTTP APIs. It includes structured logging, middleware, error handling, and configuration management with
minimal dependencies.

## Getting Started

1. Clone the repository:
    ```sh
    git clone https://github.com/shivanshkc/squelette.git
    cd squelette
    ```

2. Replace `squelette` with your desired project name in all files and directories.

3. Rename the `cmd/squelette` folder to your desired project name.

4. Create a config file by running:
    ```sh
    cp config/config.example.json config/config.json
    ```

5. Run using:
    ```sh
    make run
    ```

## Makefile Commands

The `Makefile` includes several commands to streamline common tasks:

- `make build`: Build the project.
- `make run`: Tidy dependencies, build, and run the project.
- `make image`: Build the container image of the project.
- `make container`: Run an application container.
- `make test`: Run tests for the project.
- `make lint`: Run linters to check code quality.
- `make tidy`: Run `go mod tidy`.

## Adding an API

New REST endpoints are added in the `addRoutes()` method of the Handler struct (`internal/rest/rest.go`).
Here's a basic example:

```go
mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
    httputils.WriteJson(w, http.StatusOK, nil, map[string]any{"code": "OK"})
})
```

## Project Structure

```
cmd/
└── squelette/
    └── main.go           # Entry point, HTTP server setup, graceful shutdown
config/
└── config.example.json   # Example configuration file
internal/
├── config/               # Configuration loading
├── logger/               # Structured logging with context support
└── rest/                 # HTTP handler, routing, and middleware
pkg/
└── utils/
    └── httputils/        # HTTP response helpers and error types
```

## Adding Complex Handlers

For larger handlers, create methods on the `Handler` struct in `internal/rest/rest.go`:

1. **Add a handler method**:
```go
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
    // Your business logic here
    httputils.WriteJson(w, http.StatusOK, nil, users)
}
```

2. **Register in `addRoutes()`** (`internal/rest/rest.go`):
```go
mux.HandleFunc("GET /api/users", h.GetUsers)
```

3. **For new domains, create separate files** under `internal/rest/`:
```go
// internal/rest/users.go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) { ... }
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) { ... }
```

## Error Handling

Use the built-in error utilities in `pkg/utils/httputils` for consistent responses:

```go
httputils.WriteError(w, httputils.NotFound().WithReasonStr("User not found"))
httputils.WriteError(w, httputils.BadRequest().WithReasonStr("Invalid input"))

// Available error types: BadRequest, Unauthorized, PaymentRequired, Forbidden,
// NotFound, RequestTimeout, Conflict, PreconditionFailed, InternalServerError,
// ServiceUnavailable
```

## Middleware

Middleware is defined in `internal/rest/middleware.go`. The following middleware is applied by default (in `addMiddleware()`):

- **Recovery**: Recovers from panics and returns a 500 response.
- **Access Logger**: Logs incoming requests and outgoing responses with correlation IDs.
- **CORS**: Handles cross-origin requests based on configured allowed origins.
- **Body Size Limit**: Limits request body size (default 16 KB).

To add new middleware, create a function in `internal/rest/middleware.go`:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic
        next.ServeHTTP(w, r)
    })
}
```

Then wire it into `addMiddleware()` in `internal/rest/rest.go`:

```go
func (h *Handler) addMiddleware(conf config.Config) {
    next := bodySizeLimitMiddleware(h.underlying, maxBodyReadBytes)
    next = corsMiddleware(next, conf.HttpServer.AllowedOrigins, conf.HttpServer.CorsMaxAgeSec)
    next = accessLoggerMiddleware(next)
    next = authMiddleware(next)
    next = recoveryMiddleware(next)

    h.underlying = next
}
```

## Response Helpers

Use `httputils.WriteJson()` for consistent JSON responses:

```go
// Success response
httputils.WriteJson(w, http.StatusOK, nil, data)

// With custom headers
headers := map[string]string{"X-Total-Count": "100"}
httputils.WriteJson(w, http.StatusOK, headers, data)
```
