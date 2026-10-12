# Gemini & Antigravity Agent Guidelines

See [AGENTS.md](AGENTS.md) for the complete Single Source of Truth architecture and security reference.

## Critical Instructions
- **Do NOT execute `git commit` or `git push`** without explicit user permission.
- Always execute `go vet ./...` and `go test -count=1 ./...` to verify any modifications.
- Preserve Orchid Starter Clean Architecture structure:
  - `internal/bootstrap/container`
  - `internal/bootstrap/server`
  - `modules/<name>/{delivery,usecase,repository,domain}`
  - `security/` & `pkg/cookie/`
- Use `github.com/go-chi/chi/v5`, `github.com/olympus-forge/logos`, and `github.com/olympus-forge/bunker`.
