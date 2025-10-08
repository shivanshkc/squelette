# Squelette

## Introduction

Squelette (French for "Skeleton") is a Go web service template that provides a clean, production-ready foundation for building HTTP APIs. It includes structured logging, middleware, error handling, and configuration management with minimal dependencies.

## Getting Started

1. Clone the repository:
    ```sh
    git clone https://github.com/shivanshkc/squelette.git
    cd squelette
    ```

2. Replace `squelette` with your desired project name in all files and directories.

3. Rename the `cmd/squelette` folder to your desired project name.

4. Delete the `CHANGELOG.md` file.

5. Create a configs file by running:
    ```sh
    cp config/config.sample.json config/config.json
    ```

6. Run using:
    ```sh
    make run
    ```

## Makefile Commands

The `Makefile` includes several commands to streamline common tasks:

- `make build`: Build the project.
- `make run`: Compile and run the project.
- `make image`: Build the container image of the project.
- `make container`: Run an application container.
- `make test`: Run tests for the project.
- `make lint`: Run linters to check code quality.

## Adding an API

### Quick Start
New REST endpoints are added in the `mux()` method of the Server struct (`internal/http/server.go`). Here's a basic example:

```go
mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
    httputils.Write(w, http.StatusOK, nil, map[string]any{"users": []string{}})
})
```

### Project Structure
The project follows a layered architecture:

```
internal/
├── handlers/          # Business logic for API endpoints
├── http/             # HTTP server setup and routing
├── middleware/       # Request middleware (CORS, logging, recovery)
└── utils/
    ├── httputils/    # HTTP response helpers
    └── errutils/     # Error handling utilities
```

### Adding Complex Handlers

For larger handlers, create methods on the `Handler` struct:

1. **Add to existing handler** (`internal/handlers/handler.go`):
```go
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
    // Your business logic here
    httputils.Write(w, http.StatusOK, nil, users)
}
```

2. **Register in mux** (`internal/http/server.go`):
```go
// The Server struct now has a Handler field
mux.HandleFunc("/api/users", s.Handler.GetUsers)
```

3. **For new domains, create separate files**:
```go
// internal/handlers/user_handler.go
type UserHandler struct{}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) { ... }
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) { ... }
```

### Error Handling

Use the built-in error utilities for consistent responses:

```go
// Return standard HTTP errors
httputils.WriteErr(w, errutils.NotFound().WithReasonStr("User not found"))
httputils.WriteErr(w, errutils.BadRequest().WithReasonStr("Invalid input"))

// Available error types: BadRequest, Unauthorized, Forbidden, NotFound, 
// Conflict, InternalServerError, ServiceUnavailable
```

### Adding Middleware

Create new middleware in `internal/middleware/mw.go`:

```go
func (m Middleware) Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic
        next.ServeHTTP(w, r)
    })
}
```

Apply middleware in the `mux()` method:
```go
return mw.CORS(mw.Auth(mw.AccessLogger(mw.Recovery(mux))))
```

### Response Helpers

Use `httputils.Write()` for consistent JSON responses:

```go
// Success response
httputils.Write(w, http.StatusOK, nil, data)

// With custom headers
headers := map[string]string{"X-Total-Count": "100"}
httputils.Write(w, http.StatusOK, headers, data)
```
