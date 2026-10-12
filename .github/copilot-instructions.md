# GitHub Copilot Instructions

Refer to [AGENTS.md](../AGENTS.md) for full architecture details.

## Core Rules
1. Never commit or push without explicit approval from the developer.
2. Architecture follows Orchid Starter Clean Architecture (`modules/<domain>/{delivery,usecase,repository,domain}`).
3. Framework is `github.com/go-chi/chi/v5` with `github.com/olympus-forge/logos` and `github.com/olympus-forge/bunker`.
4. Usecases must not import `net/http` or HTTP-specific types; keep business logic isolated.
5. All security standards (Dual-Secret JWT, Refresh Token Rotation, Double Submit CSRF) must be preserved.
6. Verify all code changes with `go vet ./...` and `go test -count=1 ./...`.
