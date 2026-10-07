# 📋 Analisa Backend — UI GreenMetric (Goravel)

> **Dibuat:** 2026-10-07 · **Dibuat oleh:** opencode (Atria-Dawn-Preview)
> **Lokasi:** `.../UI GREENMETRIC/be_greenmetric` · **Branch:** `main` · **Module Go:** `ui_greenmetric`
> **Sifat:** Backend-only. Frontend (Next.js, di folder `frontend/`) dikerjakan tim terpisah dan di-`.gitignore`.

---

## 1. Identitas & Tech Stack

| Item | Nilai |
|---|---|
| Bahasa | Go 1.25.0 |
| Framework | Goravel v1.18.0 (Laravel-style untuk Go) |
| HTTP Driver | Gin (`github.com/goravel/gin`) |
| ORM | GORM via Goravel ORM |
| DB aktif | **MySQL** — `.env`: `DB_CONNECTION=mysql`, `127.0.0.1:8889`, db `ui_greenmetric`, user/pass `root/root` |
| Driver DB lain | Postgres juga ter-autoload (siap pakai, ganti `DB_CONNECTION`) |
| Auth | JWT (guard `user`, driver jwt, TTL 60 menit, refresh 2 minggu) |
| Hashing | bcrypt rounds 12 |
| Cache | memory (prefix `{APP_NAME}_cache`) |
| Queue | `sync` (opsional `database` → tabel `jobs`/`failed_jobs`) |
| Logging | `stack` → `daily` → `storage/logs/goravel-YYYY-MM-DD.log` (retention 7 hari, level debug) |
| Storage | disk `public_dir` → `./public` (untuh evidence), `local`/`public` → `storage/app` |
| AI | Provider OpenAI ter-autoload, **belum dipakai** (`AI_PROVIDER=` kosong) |
| Telemetry | OpenTelemetry aktif (exporter OTLP `localhost:4318`) |
| Build | ✅ `go build ./...` **lulus tanpa error** |

**Bootstrap chain:** `main.go` → `bootstrap.Boot()` → 26 ServiceProvider (log, cache, hash, http, session, filesystem, validation, view, route, gin, ai, openai, database, postgres, mysql, auth, crypt, queue, event, grpc, translation, mail, schedule, telemetry, testing) → `routes.Web()` + `routes.Grpc()` → `app.Start()`.

**Tooling:** `./artisan` (CLI Goravel: migrate, seed, jwt:secret), `.air.toml` (hot-reload → `./tmp/main.exe`).

---

## 2. Database — 10 Tabel

Migration terdaftar di `bootstrap/migrations.go`:

| # | File | Tabel | Field kunci |
|---|---|---|---|
| 1 | `20210101000001_create_jobs_table.go` | `jobs`, `failed_jobs` | queue driver database |
| 2 | `20260713154730_create_campuses_table.go` | `campuses` | code (unique), name, institution_type, climate, setting |
| 3 | `20260713154824_create_users_table.go` | `users` | campus_id (FK), name, email (unique), password, role |
| 4 | `20260713155115_create_categories_table.go` | `categories` | code (unique), name, max_points, weight_percentage |
| 5 | `20260713155133_create_indicators_table.go` | `indicators` | category_id (FK), code (unique), title, input_type enum, max_points |
| 6 | `20260713155145_create_indicator_scoring_tiers_table.go` | `indicator_scoring_tiers` | indicator_id (FK), option_label, min/max_value, operator enum, point_multiplier |
| 7 | `20260713155153_create_campus_assessments_table.go` | `campus_assessments` | campus_id (FK), assessment_year, overall_score, status enum |
| 8 | `20260713155204_create_assessment_answers_table.go` | `assessment_answers` | campus_assessment_id (FK), indicator_id (FK), raw_input_data JSON, calculated_value, selected_tier_id, earned_points |
| 9 | `20260713155211_create_assessment_evidences_table.go` | `assessment_evidences` | assessment_answer_id (FK), document_name, description, file_url |
| 10 | `20260719091940_create_indicator_fields_table.go` | `indicator_fields` | indicator_id (FK), key, label, type, options, required |

Semua FK `CascadeOnDelete`.

### Diagram relasi

```
Campus ──1:N── User                    (user.campus_id)
   │
   └──1:N── CampusAssessment ──1:N── AssessmentAnswer ──1:N── AssessmentEvidence
                (campus_assessment_id)  (assessment_answer_id)
                                            │  N:1
                                            ▼
                                         Indicator ──1:N── IndicatorField
                                            │            (indicator_id)
                                            │ 1:N
                                            ▼
                                  IndicatorScoringTier
                                  (indicator_id)

Category ──1:N── Indicator               (indicator.category_id)
```

Tidak ada relasi ManyToMany.

### Seeder (idempoten, di `bootstrap/seeders.go`)

- **7 kategori:** SI (1100, 11%), EC (2000, 20%), WS (1700, 17%), WR (1100, 11%), TR (1700, 17%), ED (1300, 13%), GD (1100, 11%)
- **60 indikator:** SI1–SI8, EC1–EC10, WS1–WS6, WR1–WR6, TR1–TR8, ED1–ED10, GD1–GD12 + `IndicatorField` dinamis per kode (mis. SI1 → `luas_total`, `luas_dasar`; SINGLE_CHOICE → field `option_label` tipe `choice`)
- **Tier penilaian** per indikator (umumnya 5 tier, multiplier 0 / 0.05 / 0.25 / 0.50 / 0.75 / 1.00; operator `<=`, `<`, `>=`, `>`, `BETWEEN`, `CHOICE`)
- **1 kampus + 3 user** (password: `secretpassword`): `superadmin@gmail.com` (SUPER_ADMIN), `adminkampus@gmail.com` (ADMIN_KAMPUS), `operatorsi@gmail.com` (OPERATOR_SI)

---

## 3. Endpoint API — 23 Route

Global middleware: `NotFoundMiddleware` (JSON 404) + `Recover()` (JSON 500). CORS ditangani otomatis driver Gin.

| Method | Path | Handler | Middleware |
|---|---|---|---|
| GET | `/` | view `welcome.tmpl` | global |
| GET | `/public/*` | static file (`./public`) | global |
| GET | `/logs` | `LogController.ViewLogs` | secret query param |
| POST | `/api/v1/auth/login` | `AuthController.Login` | publik |
| POST | `/api/v1/auth/logout` | `AuthController.Logout` | JWT+RBAC |
| GET/POST/PUT/DELETE | `/api/v1/campuses[/{id}]` | `CampusController` | JWT+RBAC (SUPER_ADMIN) |
| GET/POST/PUT/DELETE | `/api/v1/admin/indicators[/{id}]` | `AdminIndicatorController` | JWT+RBAC (SUPER_ADMIN) |
| GET/POST/PUT/DELETE | `/api/v1/users[/{id}]` | `UserController` | JWT+RBAC (ADMIN_KAMPUS / SUPER_ADMIN) |
| GET | `/api/v1/assessments/dashboard` | `DashboardController.Index` | JWT+RBAC |
| GET | `/api/v1/categories` | `AssessmentController.GetAllCategoriesWithIndicators` | JWT+RBAC |
| GET | `/api/v1/categories/{category_code}/indicators` | `AssessmentController.GetIndicatorsByCategory` | JWT+RBAC |
| POST | `/api/v1/assessments/answers` | `AssessmentController.SaveAnswer` | JWT+RBAC |
| POST | `/api/v1/assessments/submit` | `AssessmentController.SubmitAssessment` | JWT+RBAC |
| POST | `/api/v1/evidences/upload` | `EvidenceController.Upload` | JWT+RBAC |
| DELETE | `/api/v1/evidences/{id}` | `EvidenceController.Destroy` | JWT+RBAC |

`routes/grpc.go` **kosong** (tidak ada endpoint gRPC terpakai).

### Layer aplikasi

- **Controllers** (`app/http/controllers/`): `AuthController`, `CampusController`, `UserController`, `AdminIndicatorController`, `AssessmentController`, `EvidenceController`, `DashboardController`, `LogController`.
- **Request class**: `app/http/requests/save_answer_request.go` — `SaveAnswerRequest` (`indicator_code`, `assessment_year`, `raw_input_data`, semua required).
- **Middleware** (`app/http/middleware/`): `jwt_middleware.go` (verifikasi `Authorization: Bearer`), `rbac_middleware.go` (role-based), `not_found_middleware.go` (404 global).
- **Services** (`app/services/`): `AuthService` (login/logout + JWT), `ScoringService` (mesin skoring).

---

## 4. RBAC & Multi-tenancy

- **Role:** `SUPER_ADMIN`, `ADMIN_KAMPUS`, `OPERATOR_<KODE_KATEGORI>` (SI/EC/WS/WR/TR/ED/GD).
- **Aturan** (`rbac_middleware.go`): SUPER_ADMIN akses penuh; ADMIN_KAMPUS kelola user kampusnya sendiri + semua category form; operator hanya boleh akses kategori sesuai suffiks role-nya (dicek vs `{category_code}` di path). Path `/campuses` & `/admin/indicators` eksklusif SUPER_ADMIN. Lainnya → 403.
- **Multi-tenancy** via `campus_id` pada model `User` & `CampusAssessment`.
- **Workflow assessment:** `DRAFT → SUBMITTED` (kunci, tidak bisa ubah/evidence) → `VERIFIED` (enum ada, **endpoint verify belum ada**).

---

## 5. Scoring Engine — Core Bisnis

`app/services/scoring_service.go`:

1. `Calculate(assessmentID, indicatorCode, rawInput)` → ambil/buat `AssessmentAnswer`, persist `raw_input_data` (JSON).
2. `SINGLE_CHOICE` → cocokkan `option_label` ke `IndicatorScoringTier` → `EarnedPoints = PointMultiplier × MaxPoints`.
3. `NUMERIC_FORMULA` → `EvaluateFormula()` lalu threshold matching via operator (`<=`, `<`, `>=`, `>`, `BETWEEN`); fallback multiplier 0 jika tidak ada tier cocok.
4. `UpdateOverallScore()` → rekap total `EarnedPoints` ke `CampusAssessment.OverallScore`.

**Cakupan (diverifikasi):**
- **20 indikator NUMERIC_FORMULA** → **semua** punya rumus eksplisit di `EvaluateFormula` (20 `case`): SI1–SI4, EC1/EC2/EC4/EC5/EC8, WR1/WR5, TR1/TR4/TR5, ED1/ED2/ED3/ED10, GD1/GD8. Validasi div-by-zero ada.
- **40 indikator SINGLE_CHOICE** → tier matching via `option_label`.
- Fallback untuk indikator tak dikenal: field pertama bertipe `float`/`int`, atau key `value`.

✅ **Kesimpulan: scoring engine lengkap untuk seluruh 60 indikator.**

---

## 6. Status & Kualitas

- ✅ **Build bersih** — `go build ./...` lulus.
- ✅ **TODO.md:** Fase 1–5 (setup+RBAC, master data+dashboard, scoring engine, upload+submit, error handling) **semua selesai**.
- ✅ Aplikasi pernah berjalan (log `storage/logs/`: 2026-07-14, 2026-07-19, 2026-08-07).
- ⚠️ **Test nyaris tidak ada** — hanya `tests/feature/example_test.go` (1 assertion `s.True(true)`). Belum ada test controller/scoring/middleware/seeder.
- 📁 Evidence: disimpan ke disk `public_dir` (`./public/evidences/`), dilayani via `GET /public/*`. Validasi PDF/JPG/JPEG/PNG, max 2MB.
- 📚 Dokumentasi di repo: `SRS.md`, `SDD.md`, `KONTRAK_API.md`, `PANDUAN_API.md`, `KRITERIA_INDIKATOR_PENILAIAN.md`, `WALKTHROUGH.md`, `TODO.md`.

---

## 7. ⚠️ Temuan & Isu

| # | Severity | Temuan | Lokasi |
|---|---|---|---|
| 1 | **Tinggi** | **Port mismatch Docker:** `.env` `APP_PORT=3030`, tapi `docker-compose.yml` memetakan `"3000:3000"` → service tidak terjangkau. | `.env`, `docker-compose.yml` |
| 2 | **Tinggi** | **Endpoint `/logs` baca file salah:** membaca `storage/logs/goravel.log` (mode single), padahal logging default `stack → daily` (`goravel-YYYY-MM-DD.log`). Endpoint kemungkinan menampilkan file kosong/tidak ada. | `app/http/controllers/log_controller.go:33` |
| 3 | **Sedang** | **Secret default endpoint `/logs`:** `logs_secret` default `super-secret-logs-key`, `.env` tidak override → kunci bisa ditebak. | `config/app.go` |
| 4 | **Sedang** | **Secret terbake ke image Docker:** `Dockerfile` men-`COPY .env` (berisi `JWT_SECRET` + kredensial DB). | `Dockerfile` |
| 5 | **Sedang** | **`APP_DEBUG=true`** di `.env` local — OK untuk dev, harus `false` di produksi. | `.env` |
| 6 | **Rendah** | Belum ada endpoint register user, change password, atau refresh token (frontend harus login ulang tiap 60 menit). | `routes/web.go` |
| 7 | **Rendah** | Status `VERIFIED` belum punya endpoint verify (untuk SUPER_ADMIN). | `routes/web.go` |
| 8 | **Info** | CORS terbuka penuh (`paths/origins/methods = ["*"]`) — **aman untuk integrasi frontend lintas-domain**; CORS memang otomatis didaftarkan driver Gin (`gin/route.go:62`), bukan via middleware app. | `config/cors.go` |
| 9 | **Info** | Folder `frontend/` sengaja di-`.gitignore` (commit `00af982`) — sudah tepat dipisahkan. | `.gitignore` |
| 10 | **Info** | `hashing` bcrypt + `session` driver file ada, tapi API pure JSON (tidak pakai session/cookie). | `config/` |

---

## 8. Notes untuk Tim Frontend

- **Base URL dev:** `http://127.0.0.1:3030/api/v1` (bukan 3000).
- **CORS:** allow-all → tidak ada masalah cross-origin dari Next.js.
- **Login:** `POST /api/v1/auth/login` → return token + data user + data kampus. Semua endpoint lain butuh header `Authorization: Bearer <token>`.
- **Response format standar:**
  - Success: `{"status":"success","message":"...","data":{...}}`
  - Error: `{"status":"error","code":4xx,"message":"...","errors":{...}}`
- **File bukti** diakses via `file_url` = `/public/evidences/<nama_file>` (dilayani static route).
- **Kontrak detail:** lihat `KONTRAK_API.md` & `PANDUAN_API.md` di root repo.

---

## 9. Cara Menjalankan (Dev)

```bash
# 1. Pastikan MySQL berjalan di 127.0.0.1:8889 (db: ui_greenmetric, root/root)
#    atau sesuaikan DB_* di .env
./artisan migrate          # jalankan migration
./artisan db:seed          # seed kategori + 60 indikator + tier + user default

# 2. Jalankan server
go run .                   # atau: ./artisan serve, atau air (hot-reload)
# Server mendengarkan di http://127.0.0.1:3030

# 3. Cek build
go build ./...

# 4. Test (minimal saat ini)
go test ./tests/...
```

**Login uji coba:** `superadmin@gmail.com` / `secretpassword` (SUPER_ADMIN) · `adminkampus@gmail.com` (ADMIN_KAMPUS) · `operatorsi@gmail.com` (OPERATOR_SI)
