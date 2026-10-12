# CLAUDE.md - Instructions for Claude Code & AI Assistants

See [AGENTS.md](AGENTS.md) for detailed architecture documentation.

## Critical Safety Rules
- **DO NOT COMMIT OR PUSH** without explicit user permission/acceptance.
- Always run `go vet ./...` and `go test -count=1 ./...` before finishing tasks.

## Build and Test Commands
- **Build API**: `go build ./cmd/api`
- **Run tests**: `go test -v ./...`
- **Run tests (no cache)**: `go test -count=1 ./...`
- **Static analysis**: `go vet ./...`
- **Dependencies**: `go mod tidy`

## Architecture: Orchid Starter Clean Architecture
- **Web Router**: `github.com/go-chi/chi/v5`
- **Logger**: `github.com/olympus-forge/logos`
- **Errors**: `github.com/olympus-forge/bunker`
- **Layers per Module (`modules/<domain>`)**:
  - `delivery/api/rest/v1`: Handlers and HTTP DTOs
  - `usecase`: Business logic (pure Go, decoupled from HTTP)
  - `repository`: Data access interfaces and implementations
  - `domain/models` & `domain/errors.go`: Entities and errors
- **Bootstrap & DI**: `internal/bootstrap/container` and `internal/bootstrap/server`
- **Security**:
  - `security/`: Dual-Secret JWT, HMAC-SHA256 Double-Submit CSRF, SHA-256 token hashing
  - `pkg/cookie/`: HttpOnly, SameSite, Secure cookie abstraction
  - `internal/middleware/`: Auth and CSRF middlewares
