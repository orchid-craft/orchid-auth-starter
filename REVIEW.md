# Ulasan Proyek `orchid-auth-starter` — Modular Monolith Auth + CSRF Protection

Tanggal ulasan: 12 Oktober 2026
Metode: baca seluruh struktur kode (`security/`, `modules/auth`, `modules/user`, `internal/middleware`, `pkg/cookie`, `internal/bootstrap`), `go vet ./...` dan `go test -count=1 ./...` (semua PASS).

---

## 1. Ringkasan Eksekutif

Proyek ini adalah starter backend Go dengan arsitektur Modular Monolith + Clean Architecture yang menerapkan standar keamanan modern: Dual-Secret JWT, Refresh Token Rotation (RTR) + Reuse Detection, cookie hardened, dan Double-Submit CSRF + HMAC.

Status: **template pembelajaran / portfolio yang sangat bagus, belum production-grade.** Layak dipertahankan dan dilanjutkan setelah 3 perbaikan prioritas P0.

---

## 2. Yang Sudah Benar (Pertahankan)

### 2.1 Arsitektur rapi dan disiplin
- Pemisahan `delivery/ -> usecase/ -> repository/ + domain/` dipatuhi.
- `modules/auth/usecase/auth_usecase.go` tidak mengimpor `net/http` — murni business logic.
- Komunikasi `auth -> user` hanya lewat interface (`UserUsecaseInterface`).
- Klaim user disuntik lewat `security/context.go`, bukan impor antar modul.
- DI terpusat di `internal/bootstrap/container/application_container.go`.

### 2.2 Keamanan — konsep kelas industri
- Dual-secret isolation benar:
  - `ValidateAccessToken` pakai `AccessTokenSecret`.
  - `ValidateRefreshToken` pakai `RefreshTokenSecret`.
  - Cek `SigningMethodHMAC` anti serangan `alg:none`.
- CSRF Double-Submit + HMAC-SHA256 + `subtle.ConstantTimeCompare` di `security/csrf.go` — benar.
- Atribut cookie tepat di `pkg/cookie/cookie.go`:
  - `access_token`: `HttpOnly=true`, `SameSite=Lax`, `Path=/`
  - `refresh_token`: `HttpOnly=true`, `SameSite=Strict`, `Path=/api/auth`
  - `csrf_token`: `HttpOnly=false`, `SameSite=Lax`
- Refresh token disimpan sebagai hash SHA-256, bukan plaintext.
- RTR + reuse detection dengan `RevokeAllForUser` sesuai RFC 6819 di `auth_usecase.go:139-140`.
- Password memakai `bcrypt.DefaultCost`.

### 2.3 Testing E2E bagus
- `tests/auth_flow_test.go` mencakup: register, login, cek atribut cookie, CSRF ditolak 403 tanpa token / dengan token palsu, CSRF lolos 200 dengan token valid, rotasi, reuse attack 403, logout `MaxAge=-1`.

### 2.4 Dokumentasi luar biasa
- `README.md` 639 baris: sequence diagram, ERD, UML, flowchart JWT, tabel cookie, contoh curl. Sangat bagus untuk onboarding dan portfolio.

---

## 3. Temuan & Kelemahan

### 3.1 P0 — Kritis, perbaiki sebelum production

- [x] **P0-1 Privilege escalation saat register** (SELESAI)
  - Lokasi: `modules/user/usecase/user_usecase.go`, `modules/auth/delivery/models/request/auth_request.go`, `modules/auth/usecase/auth_usecase.go`
  - Masalah: `req.Role` diterima mentah. Siapapun bisa `POST /api/auth/register {"role":"admin"}` menjadi admin.
  - Perbaikan: force `role="user"` saat register publik. DTO pendaftaran publik tidak mengekspos field role. Whitelist role admin.

- [x] **P0-2 `GET /api/users` tanpa RBAC** (SELESAI)
  - Lokasi: `modules/user/delivery/api/rest/handler.go`, `internal/middleware/auth_middleware.go:RequireRole`
  - Masalah: semua user terautentikasi bisa list semua user, padahal komentarnya “hanya admin”.
  - Perbaikan: tambah middleware `RequireRole("admin")` membatasi akses endpoint hanya untuk role admin.

- [x] **P0-3 Validasi input minim** (SELESAI)
  - Lokasi: `pkg/validator/validator.go`, `handler_v1.go` (auth + user), `user_request.go`, `auth_request.go`
  - Masalah: tidak ada cek format email, panjang / kompleksitas password, `display_name` kosong.
  - Perbaikan: pakai `go-playground/validator/v10` dengan pesan error informatif berbahasa Indonesia, validasi di delivery sebelum ke usecase.

### 3.2 P1 — Penting untuk production

- [x] **P1-1 Refresh via JSON body melemahkan isolasi** (SELESAI)
  - Lokasi: `modules/auth/delivery/api/rest/v1/handler_v1.go:RefreshToken`
  - Masalah: fallback `{"refresh_token":...}` memperluas attack surface, mengalahkan tujuan `Path=/api/auth` + `HttpOnly`.
  - Perbaikan: cookie-only untuk refresh. Hapus fallback body.

- [x] **P1-2 `ExpiresIn` hardcoded** (SELESAI)
  - Lokasi: `modules/auth/usecase/auth_usecase.go`, `internal/bootstrap/container/application_container.go`
  - Masalah: `ExpiresIn: 15*60` tidak ikut `cfg.AccessTokenTTL()`. Jika TTL diubah, response berbohong.
  - Perbaikan: oper TTL dari config ke usecase dan hitung `int64(u.accessTokenTTL.Seconds())`.

- [ ] **P1-3 Tanpa rate-limit / brute-force protection**
  - Lokasi: `POST /api/auth/login`, `POST /api/auth/refresh`
  - Perbaikan: tambah rate limiter (chi `throttle` atau `golang.org/x/time/rate`), lockout sementara, log percobaan gagal.

- [x] **P1-4 Secret default lemah tanpa validasi** (SELESAI)
  - Lokasi: `config/model.go`, `config/config.go`
  - Masalah: default secret di kode, tidak ada cek panjang min. 32 saat boot.
  - Perbaikan: `Fatal` jika `IsProduction() && len(secret) < 32`, wajibkan env di production.

- [ ] **P1-5 In-memory storage only**
  - Lokasi: `modules/auth/repository/session_repository.go`, `modules/user/repository/user_repository.go`
  - Masalah: hilang saat restart, tidak ada GC `expires_at`, memori tumbuh terus.
  - Perbaikan: ganti ke Postgres sesuai ERD di README + job cleanup expired session.

### 3.3 P2 — Kualitas & maintainability

- [x] **P2-1 Singleton chi router tidak testable** (SELESAI)
  - Lokasi: `internal/bootstrap/server/chi_server.go:NewChiApplication`, `internal/bootstrap/container/container.go`
  - Masalah: `sync.Once` global membuat `NewContainer()` berbagi router, rute menumpuk antar test.
  - Perbaikan: buat `NewChiApplication()` fresh per container, singleton hanya untuk `GetChiApplication()`.

- [ ] **P2-2 Access token tidak bisa direvoke saat logout**
  - Lokasi: `auth_usecase.go:Logout`
  - Masalah: wajar untuk JWT stateless, tapi window 15 menit masih bisa dipakai jika dicuri.
  - Perbaikan opsional: denylist JTI singkat di Redis / perkecil TTL ke 5 menit.

- [x] **P2-3 Error mapping belum konsisten** (SELESAI)
  - Contoh: `Register` selalu 400 termasuk `ErrUserAlreadyExists` (seharusnya 409), `ListUsers` tanpa bedakan not-found vs internal.
  - Perbaikan: mapping domain error -> HTTP status (409 Conflict, 404 Not Found, 401 Unauthorized, 403 Forbidden).

- [ ] **P2-4 Observability minim**
  - `RequestID` middleware dipasang tapi tidak dipropagasi ke log usecase. Belum ada metrics / tracing.
  - Perbaikan: suntik `request_id` ke context + structured log.

- [ ] **P2-5 Belum ada Docker / CI**
  - Lokasi: `.github/` hanya ada instruksi copilot, `scripts/test_api.sh` ada tapi tidak otomatis.
  - Perbaikan: tambah `Dockerfile`, `docker-compose` (api+postgres), GH Actions `vet + test`.

---

## 4. Rating

| Dimensi | Nilai | Catatan |
|---|---|---|
| Arsitektur & Struktur | 9.0 / 10 | Modular, DI, boundary bersih |
| Implementasi Keamanan | 8.5 / 10 | Konsep benar, ada 2 lubang P0 |
| Kualitas Kode | 8.0 / 10 | Enak dibaca, error `bunker` konsisten |
| Testing | 7.0 / 10 | E2E bagus, coverage rendah, minim unit usecase |
| Dokumentasi | 9.5 / 10 | README sangat lengkap |
| Production-Readiness | 6.0 / 10 | No DB, no Docker/CI, no rate-limit, no validasi |
| **OVERALL** | **8.1 / 10** | Good, di atas rata-rata starter Go |

Skala: 9+ excellent, 8 good, 7 cukup, <6 kurang.

---

## 5. Opini

Ini starter auth Go yang matang secara konsep — jauh di atas tutorial pada umumnya yang hanya simpan JWT di localStorage. Pilihan `SameSite=Strict + Path isolation + RTR + HMAC CSRF` menunjukkan pemahaman threat model XSS / CSRF / replay yang benar.

Kekurangan utama bukan di pemahaman security, melainkan di hardening production: validasi, RBAC, persistensi, dan rate-limit. Semua bisa diperbaiki inkremental tanpa bongkar arsitektur.

---

## 6. Kesimpulan & Langkah Berikutnya

1. Kerjakan P0-1 sampai P0-3 dulu — tanpa itu jangan expose ke internet.
2. Lalu P1-1, P1-2, P1-4 (murah, dampak besar).
3. Baru P1-5 (Postgres) + P2-1 (router fresh) + Docker/CI.
4. Setelah P0 beres, nilai proyek naik ke ~9.0 dan layak jadi template internal / showcase portfolio.

Checklist di seksi 3 bisa dipakai sebagai to-do perbaikan berulang. Tandai `[x]` setiap item selesai dan jalankan `go vet ./... && go test -count=1 ./...` setelah setiap perubahan sesuai aturan `AGENTS.md`.
