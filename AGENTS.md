# AI Coding Agent Guidelines & Repository Instructions

Welcome AI Agent (Antigravity, Claude Code, Cursor, Copilot, Windsurf, Aider, Cline, etc.).
This file defines the architecture, project standards, coding conventions, and critical safety rules for the **`orchid-auth-starter`** project.

---

## 🚨 CRITICAL SAFETY RULES (MUST OBEY AT ALL TIMES)

1. **NO COMMITS & NO PUSHES WITHOUT EXPLICIT PERMISSION**:
   - **NEVER** run `git commit` or `git push` unless the user explicitly requests or accepts it in the prompt.
   - Proposing commit messages or showing `git status` / `git diff` is allowed and encouraged, but execution of `commit` or `push` is strictly blocked until confirmed.
2. **VERIFY BEFORE COMPLETION**:
   - Always run `go vet ./...` and `go test -count=1 ./...` after making modifications.
   - Never leave broken builds or failing tests.
3. **PRESERVE DOCUMENTATION INTEGRITY**:
   - Maintain existing comments, docstrings, and architectural rationale unless specifically asked to change them.

---

## 🏗️ ARCHITECTURE: Orchid Starter Clean Architecture

This repository follows the **Orchid Starter** Clean Architecture & Modular pattern using Go standard practices:

```
.
├── cmd/
│   └── api/main.go                          # Entrypoint: Inisialisasi Container & REST Server
├── config/
│   ├── config.go                           # Config loader (caarlos0/env/v11 & godotenv)
│   └── model.go                            # Struct LocalConfig & security params
├── constants/
│   ├── api.go                              # Headers (X-Request-Id, X-CSRF-Token, etc.)
│   └── context_key.go                      # Context keys
├── internal/
│   ├── bootstrap/
│   │   ├── container/
│   │   │   ├── application_container.go    # Direct Injection (DI) untuk Usecase & Repo
│   │   │   └── container.go                # Centralized Container (App, Cfg, DI, Logger)
│   │   └── server/
│   │       ├── chi_server.go               # Inisialisasi chi.NewRouter() & middlewares
│   │       └── restful_server/
│   │           ├── application.go          # Registrasi rute modul (Auth & User)
│   │           ├── restful_server.go       # Server lifecycle, chi.Walk banner & shutdown
│   │           ├── routes.go               # Abstraksi SetupAllRoutes
│   │           └── http_util/              # SuccessResponse, ErrorResponse, HealthCheck
│   ├── common/
│   │   ├── env.go                          # Environment parsing helpers
│   │   ├── request_context.go              # RequestID context helpers
│   │   └── model/                          # Common Response & RequestContext models
│   └── middleware/
│       ├── auth_middleware.go              # Verifikasi Access Token JWT
│       └── csrf_middleware.go              # Proteksi Double-Submit CSRF untuk mutasi data
├── modules/
│   ├── auth/                               # Modul Autentikasi
│   │   ├── delivery/
│   │   │   ├── api/rest/handler.go         # Router mounting (/api/auth)
│   │   │   ├── api/rest/v1/handler_v1.go   # REST Handlers (Register, Login, Refresh, Logout)
│   │   │   └── models/                     # Request & Response DTOs
│   │   ├── domain/                         # Entitas Session & Domain Errors
│   │   ├── repository/                     # Session Repository Interface & Storage
│   │   └── usecase/                        # Business Logic (Dual-Secret, RTR, Reuse Detection)
│   └── user/                               # Modul Pengguna
│       ├── delivery/                       # REST Handlers & DTOs (/api/me, /api/users)
│       ├── domain/                         # User Entities & Domain Errors
│       ├── repository/                     # User Repository Interface & Storage
│       └── usecase/                        # User Business Logic & Password Hashing
├── pkg/
│   └── cookie/                             # Cookie Manager (HttpOnly, SameSite, Secure)
├── security/                               # Primitif Kriptografi (JWT, CSRF, SHA-256)
└── tests/
    └── auth_flow_test.go                   # End-to-end integration tests
```

### Layer Boundary Rules:
- **`delivery/`** hanya bertugas membaca HTTP request, parsing DTO, memanggil usecase, dan mengembalikan response JSON (menggunakan `httpUtil.SuccessResponse` atau `httpUtil.ErrorResponse`).
- **`usecase/`** berisi murni logika bisnis. **DILARANG** mengimpor `net/http` atau dependensi HTTP di usecase.
- **`repository/`** mengabstraksi persistensi data via interface (`repository_interface.go`).
- **`domain/`** murni berisi model entitas dan domain errors.
- **`security/`** murni menangani kriptografi dan token.

---

## 🔒 SECURITY STANDARDS

1. **Dual-Secret JWT Isolation**:
   - `AccessTokenSecret` (umur pendek: 15 menit) dan `RefreshTokenSecret` (umur panjang: 7 hari) harus selalu menggunakan secret key terpisah.
2. **Refresh Token Rotation (RTR) & Reuse Attack Detection**:
   - Setiap penggunaan refresh token lama wajib membatalkan token tersebut dan menerbitkan token baru.
   - Jika refresh token yang sudah di-revoke digunakan kembali (*reuse attempt*), seluruh keluarga token pengguna wajib di-revoke demi keamanan.
3. **Double Submit Cookie CSRF**:
   - Menggunakan token bertanda tangan HMAC-SHA256.
   - Header `X-CSRF-Token` dicocokkan dengan cookie `csrf_token` menggunakan `subtle.ConstantTimeCompare`.
4. **Cookie Security Attributes**:
   - `access_token`: `HttpOnly=true`, `SameSite=Lax`, `Path="/"`.
   - `refresh_token`: `HttpOnly=true`, `SameSite=Strict`, `Path="/api/auth"` (isolasi path).
   - `csrf_token`: `HttpOnly=false` (agar bisa dibaca frontend JavaScript), `SameSite=Lax`.

---

## 📦 DEPENDENCY CONVENTIONS

- **Router**: `github.com/go-chi/chi/v5`
- **Logger**: `github.com/olympus-forge/logos` (panggil via `logos.NewLogger()`)
- **Error Handling**: `github.com/olympus-forge/bunker` (menggunakan `bunker.New(...)` dan `httpUtil.ErrorResponse`)
- **Config**: `github.com/caarlos0/env/v11` + `github.com/joho/godotenv`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **Password Hashing**: `golang.org/x/crypto/bcrypt`

---

## 🛠️ ESSENTIAL WORKFLOW COMMANDS

| Perintah | Deskripsi |
| :--- | :--- |
| `go build ./cmd/api` | Membangun binary server API |
| `go test -v ./...` | Menjalankan seluruh unit & integration test |
| `go test -count=1 ./...` | Menjalankan pengujian tanpa cache |
| `go vet ./...` | Memeriksa static analysis Go |
| `go mod tidy` | Merapikan dependensi go.mod & go.sum |
