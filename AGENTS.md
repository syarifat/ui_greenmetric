# AGENTS.md — UI GreenMetric Backend (Goravel)

> **Update terakhir:** 2026-10-07
> 📄 **Analisa lengkap:** [`.opencode/PROJECT_ANALYSIS.md`](.opencode/PROJECT_ANALYSIS.md) — baca file ini untuk detail arsitektur, diagram relasi DB, daftar endpoint, dan temuan isu.

## Apa project ini

Backend API untuk sistem penilaian kampus berkelanjutan **UI GreenMetric**, dibangun dengan **Goravel v1.18.0** (framework Go laravel-style), driver HTTP **Gin**, ORM **GORM**, database **MySQL**.

**Backend-only.** Frontend (Next.js, folder `frontend/`) dikerjakan tim terpisah dan sengaja di-`.gitignore`.

## Fakta penting (selalu ingat)

- **Bahasa/framework:** Go 1.25 + Goravel v1.18. Konvensi = gaya Laravel (facades, artisan, migrations, seeders, controllers/services/requests).
- **DB aktif: MySQL** di `127.0.0.1:8889`, db `ui_greenmetric` (`root/root`). Driver Postgres juga ter-autoload ( tinggal ganti `DB_CONNECTION`).
- **Port server: 3030** (via `.env` `APP_PORT`; default config 3000). Base URL API: `http://127.0.0.1:3030/api/v1`.
- **Auth: JWT** guard `user`, TTL 60 menit. Semua endpoint `/api/v1/*` (kecuali `POST /auth/login`) butuh header `Authorization: Bearer <token>`.
- **RBAC:** `SUPER_ADMIN` (akses penuh), `ADMIN_KAMPUS` (kelola user/assessment kampus sendiri), `OPERATOR_<KODE_KATEGORI>` (SI/EC/WS/WR/TR/ED/GD — hanya kategori sesuai role).
- **Core bisnis = Scoring Engine** di `app/services/scoring_service.go`: 60 indikator (20 `NUMERIC_FORMULA` dengan rumus eksplisit + 40 `SINGLE_CHOICE` via tier matching) → `earned_points = point_multiplier × max_points` → rekap ke `campus_assessments.overall_score`.
- **CORS terbuka penuh** dan di-handle otomatis driver Gin (bukan via middleware app) — aman untuk frontend lintas-domain.
- **Build sehat:** `go build ./...` lulus. Test minim (hanya 1 example test).

## Struktur singkat

```
main.go / bootstrap/      bootstrap + 26 ServiceProvider
config/                   app, database, jwt, auth, cors, http, filesystems, ai, cache, queue, logging, hashing, session, mail, telemetry, grpc
routes/web.go             23 endpoint (lihat tabel di .opencode/PROJECT_ANALYSIS.md)
routes/grpc.go            kosong (tidak dipakai)
app/http/controllers/     Auth, Campus, User, AdminIndicator, Assessment, Evidence, Dashboard, Log
app/http/middleware/      jwt, rbac, not_found
app/http/requests/        SaveAnswerRequest
app/models/               9 model (Campus, User, Category, Indicator, IndicatorField, IndicatorScoringTier, CampusAssessment, AssessmentAnswer, AssessmentEvidence)
app/services/             AuthService, ScoringService
database/migrations/      10 migration
database/seeders/         kategori, 60 indikator + fields, scoring tier, user default
```

## Perintah umum

```bash
go build ./...            # verifikasi build
./artisan migrate         # jalankan migration
./artisan db:seed         # seed master data + user default
go run .                  # jalankan server (port 3030)
go test ./tests/...       # test (minim)
```

**User default (password: `secretpassword`):** `superadmin@gmail.com`, `adminkampus@gmail.com`, `operatorsi@gmail.com`

## Isu yang perlu diketahui

1. **Port mismatch Docker** — `.env` `APP_PORT=3030` vs `docker-compose.yml` `"3000:3000"` → service tidak terjangkau di Docker.
2. **Endpoint `/logs` baca path salah** — `log_controller.go:33` membaca `storage/logs/goravel.log` (mode single), padahal logging pakai `daily` (`goravel-YYYY-MM-DD.log`).
3. **Secret default endpoint `/logs`** — `super-secret-logs-key` belum di-override `.env`.
4. **`Dockerfile` men-`COPY .env`** → `JWT_SECRET` & kredensial DB terbake ke image.
5. **`APP_DEBUG=true`** di `.env` (harus `false` saat produksi).
6. Belum ada endpoint register / change password / refresh token / verify assessment (`VERIFIED`).

## Konvensi kerja

- Response JSON standar: success `{"status":"success","message":...,"data":{...}}`, error `{"status":"error","code":4xx,"message":...,"errors":{...}}`.
- Pesan error memakai **Bahasa Indonesia**.
- Seeder idempoten (cek duplikat sebelum insert).
- Dokumentasi referensi: `SRS.md`, `SDD.md`, `KONTRAK_API.md`, `PANDUAN_API.md`, `KRITERIA_INDIKATOR_PENILAIAN.md`, `WALKTHROUGH.md`, `TODO.md`.
