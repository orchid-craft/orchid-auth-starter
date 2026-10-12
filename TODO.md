# TODO — Sisa Perbaikan untuk Dipertimbangkan

Hasil check ulang 12 Okt 2026: semua P0 + sebagian P1/P2 sudah DONE.
File ini hanya berisi yang **belum** + catatan kecil. Centang `[x]` jika selesai.

> Aturan tiap selesai: `go vet ./... && go test -count=1 ./...` harus PASS (sesuai `AGENTS.md`).

---

## P1 — Penting (kerjakan dulu)

- [x] **P1-3 Rate-limit login & refresh (anti brute-force)** (DONE)
  - Lokasi: `POST /api/auth/login`, `POST /api/auth/refresh`, `internal/middleware/rate_limit_middleware.go`
  - Implementasi: sliding window per IP, membatasi request per window waktu (default 10 req/menit), response 429 Too Many Requests + header `Retry-After`.
  - Test: `TestLoginRateLimit` PASS.

- [ ] **P1-5 Ganti in-memory ke Postgres + GC session**
  - Lokasi: `modules/auth/repository/session_repository.go`, `modules/user/repository/user_repository.go`, ERD di `README.md`
  - Usulan: implement `SessionRepositoryInterface` + `UserRepositoryInterface` versi Postgres, migrasi sesuai ERD, index `users(email)`, `sessions(user_id, expires_at)`, cron hapus `expires_at < now()`.
  - Done jika: restart tidak hilangkan data, test E2E tetap PASS dengan DB.

## P2 — Kualitas & maintainability

- [ ] **P2-2 Access-token revocation saat logout (opsional)**
  - Lokasi: `modules/auth/usecase/auth_usecase.go:Logout`
  - Usulan: denylist JTI singkat di Redis/memory sampai `exp` (window 15 mnt), atau perkecil TTL ke 5 mnt.
  - Done jika: token lama ditolak setelah logout dalam window TTL.

- [x] **P2-4 Propagasi RequestID ke log/usecase** (DONE)
  - Lokasi: `internal/bootstrap/server/chi_server.go`, `internal/common/request_context.go`
  - Implementasi: Chi `middleware.RequestID` disuntikkan ke `context.Context` via `modelCommon.SetRequestContext`, disematkan ke header respons `X-Request-Id`, dan tersedia helper `LogInfoWithContext` / `LogErrorWithContext`.
  - Test: `TestRequestIDContextPropagation` PASS.

- [x] **P2-5 Dockerfile + CI** (DONE)
  - Lokasi: `Dockerfile`, `docker-compose.yml`, `.github/workflows/ci.yml`
  - Implementasi: Multi-stage Dockerfile berbasis alpine minimal, docker-compose untuk service deployment, dan GitHub Actions CI untuk otomatisasi `go vet` dan `go test`.

---

## Minor — Cepat, dampak rapi

- [x] **go mod tidy (validator direct)** (DONE)
- [x] **git add file baru (efek fix .gitignore `api` → `/api`)** (DONE)

---

## Sudah DONE (arsip, jangan dikerjakan lagi)

- [x] P0-1 force role `user` saat register publik
- [x] P0-2 `RequireRole("admin")` untuk `GET /api/users`
- [x] P0-3 validator `email/min/max` di Register/Login/Update
- [x] P1-1 refresh cookie-only (hapus fallback body)
- [x] P1-2 `ExpiresIn` dari config TTL
- [x] P1-3 Rate-limit login & refresh anti brute-force
- [x] P1-4 validasi secret min.32 di production
- [x] P2-1 `NewChiApplication()` fresh per container
- [x] P2-3 `409 Conflict` untuk email duplikat
- [x] P2-4 Propagasi RequestID ke context & logger
- [x] P2-5 Dockerfile multi-stage, docker-compose & GitHub Actions CI
- [x] Fix `.gitignore` `api` → `/api`
