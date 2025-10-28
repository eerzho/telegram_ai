# Telegram AI Backend

## Project Overview
Go backend application for a Telegram bot with AI integration (Google Genkit). The application provides an HTTP API for generating responses using artificial intelligence.

## Tech Stack
- **Language**: Go 1.25.1
- **AI Framework**: Firebase Genkit
- **HTTP Server**: Built-in Go HTTP server with custom packages
- **Dependency Injection**: simpledi
- **Validation**: go-playground/validator
- **Configuration**: env variables (godotenv, caarlos0/env)

## Project Structure

```
.
├── cmd/
│   └── app/          # Application entry point
├── config/           # Application configuration
├── internal/
│   ├── adapter/      # Adapters for external services (Genkit)
│   ├── container/    # DI container
│   ├── controller/   # HTTP controllers
│   ├── domain/       # Domain models (Message, User, Dialog)
│   ├── health_check/ # Health check endpoint
│   └── generate_response/ # AI response generation
│       ├── dto.go         # HTTP DTOs (Input, Output)
│       ├── mapper.go      # DTO to Domain mapping
│       ├── usecase.go     # Business logic
│       └── http_v1.go     # HTTP handler
├── pkg/              # Reusable packages
│   ├── httpserver/   # HTTP server
│   ├── logger/       # Logging
│   ├── logging/      # Logging middleware
│   ├── recovery/     # Recovery middleware
│   ├── cors/         # CORS middleware
│   ├── json/         # JSON utilities
│   └── sse/          # Server-Sent Events
└── prompts/          # AI prompts

```

## Core Components

### Dependency Injection
The project uses `simpledi` for dependency management. All services are registered in the container (`internal/container/container.go`).

### HTTP Server
Custom wrapper over Go's standard HTTP server (`pkg/httpserver/httpserver.go`) with graceful shutdown.

### Middleware
- **logging**: Request logging
- **recovery**: Panic recovery
- **cors**: CORS configuration

### Endpoints
- **Health Check**: `/health` - application health status
- **Generate Response**: AI response generation via Genkit
- **SSE**: Server-Sent Events for streaming responses

## Configuration

The application is configured through environment variables. See `.env.example` for configuration example.

## Development

### Running the application
```bash
go run cmd/app/main.go
```

### Linting
```bash
golangci-lint run
```

Linter configuration: `.golangci.yaml`

## Architecture Principles

1. **Clean Architecture**: Layer separation (domain, usecase, adapter, controller)
2. **Domain-Driven Design**: Pure domain models in `internal/domain` package
3. **Dependency Injection**: Using DI container
4. **Interface-based design**: Working through interfaces
5. **DTO to Domain mapping**: Separate mappers for converting HTTP DTOs to domain models
6. **Graceful Shutdown**: Proper server shutdown handling
7. **Structured Logging**: Using slog for logging

## Key Files

- `cmd/app/main.go:18` - Entry point, initialization and server startup
- `internal/container/container.go` - DI container configuration
- `internal/controller/http/router.go` - Routing configuration
- `config/config.go` - Configuration structures
- `internal/domain/message.go` - Domain models (Message, User, Dialog)
- `internal/adapter/genkit/genkit.go` - Google Genkit integration (receives domain.Dialog)
- `internal/generate_response/mapper.go` - DTO to Domain mapping
- `internal/generate_response/usecase.go` - Business logic for response generation

## Git

**Current branch**: main
**Last commit**: af68237 - fix: validate for sender and add mutex for sse

## Development Notes

- Mutex is used for SSE connections for thread-safety
- Sender validation has been added
- Project follows Go naming conventions
- All public types and functions should be documented

## Data Flow

1. **HTTP Request** → `internal/generate_response/http_v1.go`
2. **JSON Decode** → `Input` DTO (validated)
3. **Usecase** → `internal/generate_response/usecase.go:47` validates and sorts messages
4. **Mapper** → `internal/generate_response/mapper.go:5` converts `Input` → `domain.Dialog`
5. **Genkit Adapter** → `internal/adapter/genkit/genkit.go` receives clean domain model
6. **SSE Stream** → Chunks streamed back to client
