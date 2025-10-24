# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go backend service for a Telegram AI application. The project uses Go 1.25.1 and follows clean architecture principles with dependency injection via `github.com/eerzho/simpledi`.

## Development Commands

### Building and Running
```bash
go run cmd/app/main.go
```

### Dependency Management
```bash
go mod tidy        # Clean up dependencies
go mod download    # Download dependencies
```

### Linting
```bash
golangci-lint run  # Run all configured linters
```

The project uses a very strict golangci-lint configuration (`.golangci.yaml`) based on the "golden config" with 100+ enabled linters. Key linters include staticcheck, govet, revive, gosec, and many others.

### Testing
```bash
go test ./...                    # Run all tests
go test -v ./...                 # Run with verbose output
go test -run TestName ./...      # Run specific test
go test -cover ./...             # Run with coverage
```

## Architecture

### Layered Structure

The codebase follows clean architecture with clear separation of concerns:

- **cmd/app**: Application entry point with server lifecycle management
- **config**: Configuration management using environment variables (via caarlos0/env)
- **internal/container**: Dependency injection container setup
- **internal/controller/http**: HTTP handlers and middleware
- **internal/usecase**: Business logic layer with input/output DTOs
- **pkg**: Reusable packages (logger, httpserver, json)

### Dependency Injection

The project uses `github.com/eerzho/simpledi` for dependency injection. All dependencies are registered in `internal/container/container.go` with the following pattern:

```go
{
    Key:  "serviceName",
    Deps: []string{"dependency1", "dependency2"},
    Ctor: func() any {
        // Constructor logic
    },
}
```

Services are retrieved from the container using `c.MustGet("serviceName")` with type assertion.

### HTTP Server

- Uses standard library `net/http` with custom middleware
- Logging middleware: Adds request ID, logs request/response details
- Recovery middleware: Catches panics and returns 500 errors
- Custom responseWriter wrapper to capture status codes and response sizes
- Health check endpoint: `GET /_hc`

### Configuration

Configuration is loaded from environment variables with defaults. The `.env` file is auto-loaded via `github.com/joho/godotenv/autoload`. Key environment variables:

- `APP_NAME`: Application name (default: "telegram-ai")
- `APP_VERSION`: Application version (required)
- `LOGGER_LEVEL`: Log level - debug, info, warn, error (default: "info")
- `LOGGER_FORMAT`: Log format - text, json (default: "json")
- `HTTP_SERVER_HOST`: Server host (default: "")
- `HTTP_SERVER_PORT`: Server port (default: "8080")
- `HTTP_SERVER_READ_TIMEOUT`: Read timeout (default: "10s")
- `HTTP_SERVER_WRITE_TIMEOUT`: Write timeout (default: "10s")
- `HTTP_SERVER_IDLE_TIMEOUT`: Idle timeout (default: "60s")

### Logging

The project uses `log/slog` for structured logging. Logger is configured via `pkg/logger` and includes:
- Structured JSON or text output
- Contextual fields (app_name, app_version)
- Request-scoped logging in HTTP middleware

**Important**: Never use the global `log` package in non-main files (enforced by depguard linter).

### Error Handling

- Use error wrapping with `fmt.Errorf("%s: %w", op, err)` pattern
- Operation names follow the format `package.Function`
- Errors are propagated up the stack with context

## Code Style Guidelines

### Enforced by Linters

- Maximum line length: 120 characters (golines)
- Function length: Max 100 lines, 50 statements (funlen)
- Cyclomatic complexity: Max 30 per function (cyclop)
- No global variables (gochecknoglobals)
- No init functions (gochecknoinits)
- All errors must be checked (errcheck)
- Comments must end with periods (godot)
- Use named imports from this project (local-prefixes: github.com/eerzho/telegram-ai)

### Best Practices

- Use `MustNew` pattern for constructors that panic on error alongside a regular `New` function
- Keep business logic in usecase layer, HTTP details in controller layer
- Use input/output DTOs for usecase boundaries
- Prefer context-aware logging: `logger.InfoContext(ctx, "message")`
- Use `math/rand/v2` instead of `math/rand` (enforced by depguard)

## Testing Guidelines

- Test files should be in `_test.go` files
- Test files have relaxed linter rules for: bodyclose, dupl, errcheck, funlen, goconst, gosec, noctx
