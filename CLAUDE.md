# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Run the server
just run
# or: go run cmd/server/main.go

# Generate TLS certificates (requires mkcert)
just gen-certs "localhost 127.0.0.1"
```

## Configuration

Config files are TOML format in `configs/server/`, selected via `PROFILE` env var (defaults to "local").

### TOML Structure

```toml
[cors]
allowed_origins="http://localhost:5173"  # comma-separated for multiple origins

[tls.server]
cert_file="certs/server.crt"  # omit both to disable TLS
key_file="certs/server.key"

[log]
level="DEBUG"  # DEBUG, INFO, WARN, ERROR
```

Config is parsed into `config.Server` struct via Viper with mapstructure tags. See `pkg/config/parser.go` for the full schema.

## Architecture

This is a Go HTTP server starter template using the standard library's `net/http` with TLS support.

### Package Structure

- **cmd/server**: Application entry point. Calls `server.Run()` with context, stdout, and env getter.
- **pkg/server**: HTTP server setup and routing. Uses functional options pattern (`WithPort`, `WithHandler`, `WithTLSConfig`) for server configuration. Router uses default middleware stack (request logger + CORS).
- **pkg/config**: Viper-based configuration parsing. Reads TOML files based on profile name.
- **pkg/httplib**: HTTP utilities including:
  - Middleware composition via `CreateStack()` - applies middlewares top-down (first arg is innermost)
  - JSON response helpers with typed `SuccessResponse[T]` and `ErrorResponse`
  - Response codes enum (Success, InternalError, ValidationError, etc.)
  - TLS configuration from cert files

### Middleware Pattern

Routes are registered with `router.registerRoute(pattern, handler, ...middlewares)`. Per-route middlewares are composed on top of the default stack (request logging + CORS).

## Standard Go Project Layout

This project follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) conventions:

### Core Directories

- **cmd/**: Main applications. Each subdirectory is an executable. Keep code minimal here—import from `pkg/` and `internal/`.
- **pkg/**: Library code safe for external import. Signals reusability to other projects.
- **internal/**: Private code the Go compiler prevents external imports of. Use `internal/app/` for app code, `internal/pkg/` for shared internal libraries.
- **vendor/**: Dependencies via `go mod vendor`. Optional with Go modules unless offline builds are needed.

### Service & Configuration

- **api/**: OpenAPI/Swagger specs, JSON schemas, protocol definitions.
- **configs/**: Configuration templates and defaults (TOML, confd, consul-template).
- **web/**: Static assets, server-side templates, SPAs.

### Build & Deploy

- **build/**: Packaging (Docker, deb, rpm) and CI configs.
- **deployments/**: IaC templates (docker-compose, Kubernetes/Helm, Terraform).
- **scripts/**: Build, install, and analysis scripts.

### Testing & Documentation

- **test/**: External test apps and test data.
- **docs/**: Design docs and user guides.
- **examples/**: Sample usage code.
- **tools/**: Supporting utilities that may import from `pkg/` and `internal/`.

### Avoid

- **src/**: Not recommended for Go—creates unnecessary nesting and conflicts with Go workspace conventions.
