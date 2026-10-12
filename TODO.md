# TODO — Sisa Perbaikan untuk Dipertimbangkan

Hasil check ulang 12 Okt 2026: semua P0 + sebagian P1/P2 sudah DONE.
File ini hanya berisi yang **belum** + catatan kecil. Centang `[x]` jika selesai.

> Aturan tiap selesai: `go vet ./... && go test -count=1 ./...` harus PASS (sesuai `AGENTS.md`).

---

## P1 — Penting (kerjakan dulu)

- [ ] **P1-3 Rate-limit login & refresh (anti brute-force)**
  - Lokasi: `POST /api/auth/login`, `POST /api/auth/refresh`, `internal/bootstrap/server/chi_server.go`
  - Usulan: chi `throttle` atau `golang.org/x/time/rate` per IP + per email, mis. 10 req/menit login, lockout sementara + log gagal.
  - Done jika: login ke-11 dalam 1 menit → 429, test baru `TestLoginRateLimit` PASS.

- [ ] **P1-5 Ganti in-memory ke Postgres + GC session**
  - Lokasi: `modules/auth/repository/session_repository.go`, `modules/user/repository/user_repository.go`, ERD di `README.md`
  - Usulan: implement `SessionRepositoryInterface` + `UserRepositoryInterface` versi Postgres, migrasi sesuai ERD, index `users(email)`, `sessions(user_id, expires_at)`, cron hapus `expires_at < now()`.
  - Done jika: restart tidak hilangkan data, test E2E tetap PASS dengan DB.

## P2 — Kualitas & maintainability

- [ ] **P2-2 Access-token revocation saat logout (opsional)**
  - Lokasi: `modules/auth/usecase/auth_usecase.go:Logout`
  - Usulan: denylist JTI singkat di Redis/memory sampai `exp` (window 15 mnt), atau perkecil TTL ke 5 mnt.
  - Done jika: token lama ditolak setelah logout dalam window TTL.

- [ ] **P2-4 Propagasi RequestID ke log/usecase**
  - Lokasi: `internal/bootstrap/server/chi_server.go`, `internal/common/request_context.go`, `internal/middleware/`
  - Usulan: suntik `X-Request-Id` chi ke `context.Context`, pakai di semua `logger.Info` + usecase.
  - Done jika: tiap baris log request memuat `request_id` yang sama dari middleware → handler.

- [ ] **P2-5 Dockerfile + CI**
  - Lokasi: root, `.github/`
  - Usulan: `Dockerfile` multi-stage, `docker-compose.yml` (api + postgres), workflow `vet + test` tiap push.
  - Done jika: `docker build` sukses, CI hijau.

---

## Minor — Cepat, dampak rapi

- [ ] **go mod tidy (validator direct)**
  - `validator/v10` masih `// indirect` di `go.mod` padahal dipakai langsung di `pkg/validator/`.
  - Jalankan: `go mod tidy && go vet ./...`

- [ ] **git add file baru (efek fix .gitignore `api` → `/api`)**
  - Masih `??`: `cmd/`, `modules/auth/delivery/api/`, `modules/user/delivery/api/`, `pkg/validator/`, `REVIEW.md`, file ini.
  - Jalankan: `git add cmd modules/auth/delivery/api modules/user/delivery/api pkg/validator REVIEW.md TODO.md .gitignore config/config.go go.mod go.sum internal/ modules/auth/usecase modules/user/usecase tests/` lalu review `git status` sebelum commit (jangan commit/push tanpa izin eksplisit).

---

## Sudah DONE (arsip, jangan dikerjakan lagi)

- [x] P0-1 force role `user` saat register publik
- [x] P0-2 `RequireRole("admin")` untuk `GET /api/users`
- [x] P0-3 validator `email/min/max` di Register/Login/Update
- [x] P1-1 refresh cookie-only (hapus fallback body)
- [x] P1-2 `ExpiresIn` dari config TTL
- [x] P1-4 validasi secret min.32 di production
- [x] P2-1 `NewChiApplication()` fresh per container
- [x] P2-3 `409 Conflict` untuk email duplikat
- [x] Fix `.gitignore` `api` → `/api`
