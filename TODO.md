# 🛠️ Backend TODO List — UI GreenMetric (Goravel + MySQL)

> Referensi: [SDD.md](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/SDD.md) · [SRS.md](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/SRS.md) · [KRITERIA_INDIKATOR_PENILAIAN.md](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/KRITERIA_INDIKATOR_PENILAIAN.md)

---

## ⚡ FASE 1: Setup Framework & Autentikasi (RBAC)

### 1.1 Setup Awal Proyek
- `[ ]` Pastikan MySQL lokal sudah berjalan dan buat database kosong `ui_greenmetric`
- `[ ]` Lengkapi kredensial database di [.env](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/.env) (`DB_USERNAME`, `DB_PASSWORD`)
- `[ ]` Set JWT secret: `./artisan jwt:secret`

### 1.2 Buat Migration (8 Tabel)
- `[x]` `./artisan make:migration create_campuses_table`
- `[x]` `./artisan make:migration create_users_table`
- `[x]` `./artisan make:migration create_categories_table`
- `[x]` `./artisan make:migration create_indicators_table`
- `[x]` `./artisan make:migration create_indicator_scoring_tiers_table`
- `[x]` `./artisan make:migration create_campus_assessments_table`
- `[x]` `./artisan make:migration create_assessment_answers_table`
- `[x]` `./artisan make:migration create_assessment_evidences_table`
- `[ ]` Jalankan semua migrasi: `./artisan migrate`

### 1.3 Buat Model (Goravel ORM / GORM)
- `[x]` `./artisan make:model Campus` (code, name, institution_type, climate, setting)
- `[x]` `./artisan make:model User` (campus_id, name, email, password, role ENUM RBAC)
- `[x]` `./artisan make:model Category` (code, name, max_points, weight_percentage)
- `[x]` `./artisan make:model Indicator` (category_id, code, title, input_type, max_points)
- `[x]` `./artisan make:model IndicatorScoringTier` (indicator_id, option_label, min_value, max_value, operator, point_multiplier)
- `[x]` `./artisan make:model CampusAssessment` (campus_id, assessment_year, overall_score, status)
- `[x]` `./artisan make:model AssessmentAnswer` (campus_assessment_id, indicator_id, raw_input_data JSON, calculated_value, selected_tier_id, earned_points)
- `[x]` `./artisan make:model AssessmentEvidence` (assessment_answer_id, document_name, description, file_url)

### 1.4 Seeder — Master Data
- `[x]` `./artisan make:seeder CategorySeeder` — isi 7 kategori (SI=1100, EC=2000, WS=1700, WR=1100, TR=1700, ED=1300, GD=1100)
- `[x]` `./artisan make:seeder IndicatorSeeder` — isi seluruh indikator per kategori (SI1–SI8, EC1–EC10, WS1–WS6, WR1–WR6, TR1–TR8, ED1–ED10, GD1–GD12)
- `[x]` `./artisan make:seeder IndicatorScoringTierSeeder` — isi logika tier skor (operator + multiplier) untuk setiap indikator sesuai [KRITERIA_INDIKATOR_PENILAIAN.md](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/KRITERIA_INDIKATOR_PENILAIAN.md)
- `[x]` `./artisan make:seeder UserSeeder` — buat akun SUPER_ADMIN default
- `[ ]` Jalankan semua seeder: `./artisan db:seed`

### 1.5 Autentikasi JWT
- `[x]` `./artisan make:controller AuthController` — implementasi login/logout
- `[x]` Buat `app/services/auth_service.go` — validasi kredensial, generate JWT token, invalidate token
- `[x]` `./artisan make:middleware JwtMiddleware` — parse & validasi JWT header `Authorization: Bearer`
- `[x]` `./artisan make:middleware RbacMiddleware` — cek `role` user sebelum akses endpoint kategori tertentu
- `[x]` Daftarkan middleware di `routes/web.go`

### 1.6 Endpoint Auth
- `[x]` `POST /api/v1/auth/login` → response token + data user + data kampus
- `[x]` `POST /api/v1/auth/logout` → invalidate token

---

## 📊 FASE 2: Modul Master Data & Dashboard

### 2.1 Manajemen Kampus (SUPER_ADMIN)
- `[x]` `./artisan make:controller CampusController`
- `[x]` `GET    /api/v1/campuses` — list semua kampus
- `[x]` `POST   /api/v1/campuses` — tambah kampus baru
- `[x]` `PUT    /api/v1/campuses/{id}` — update data kampus
- `[x]` `DELETE /api/v1/campuses/{id}` — hapus kampus

### 2.2 Manajemen User Internal Kampus (ADMIN_KAMPUS)
- `[x]` `./artisan make:controller UserController`
- `[x]` `GET    /api/v1/users` — list user milik kampus sendiri
- `[x]` `POST   /api/v1/users` — buat akun operator (7 role kategori)
- `[x]` `PUT    /api/v1/users/{id}` — update akun operator
- `[x]` `DELETE /api/v1/users/{id}` — hapus akun operator

### 2.3 Dashboard & Analisis Tren
- `[x]` `./artisan make:controller DashboardController`
- `[x]` `GET /api/v1/assessments/dashboard` — response meliputi:
  - `[x]` `campus_name`, `current_year`, `assessment_status`
  - `[x]` `overall_score` & `max_overall_score` (10.000 poin)
  - `[x]` `estimated_rank` (bandingkan dengan benchmark tahun sebelumnya)
  - `[x]` `category_breakdown[]` — earned & max per kategori (SI, EC, WS, WR, TR, ED, GD)
  - `[x]` `trend_history[]` — array riwayat skor tahunan untuk grafik garis FE

---

## 🧮 FASE 3: Core — Scoring Engine (Mesin Kalkulasi Otomatis)

### 3.1 Service Layer Skoring
- `[x]` Buat `app/services/scoring_service.go` dengan logika:
  1. `[x]` Ambil data master `Indicator` + `IndicatorScoringTier` dari MySQL
  2. `[x]` Eksekusi rumus matematika dari `raw_input_data` sesuai jenis indikator:
     - `[x]` Tipe `NUMERIC_FORMULA` → hitung persentase / rasio (SI1, SI4, EC4, EC5, EC8, TR1, dll.)
     - `[x]` Tipe `SINGLE_CHOICE` → gunakan nilai pilihan secara langsung (SI5, SI6, EC3, EC6, EC7, dll.)
  3. `[x]` **Threshold Matching** — cocokkan hasil kalkulasi dengan `indicator_scoring_tiers` berdasarkan operator (`<=`, `>`, `>=`, `BETWEEN`, `CHOICE`)
  4. `[x]` Hitung `earned_points = point_multiplier × max_points`
  5. `[x]` Simpan hasil ke `assessment_answers`
  6. `[x]` Auto-update `overall_score` di `campus_assessments`

### 3.2 Endpoint Form Assessment
- `[x]` `./artisan make:controller AssessmentController`
- `[x]` `GET  /api/v1/categories/{category_code}/indicators` — ambil soal + jawaban terakhir per kategori (RBAC: OPERATOR_SI hanya bisa akses `SI`, dst.)
- `[x]` `POST /api/v1/assessments/answers` — simpan jawaban + trigger `ScoringService`, return skor langsung ke FE

### 3.3 Request Validation
- `[x]` `./artisan make:request SaveAnswerRequest` — validasi `indicator_code`, `assessment_year`, `raw_input_data`
- `[x]` Validasi domain: nilai input tidak boleh negatif, `luas_dasar` tidak boleh melebihi `luas_total`, dll. (ditangani di ScoringService sebelum kalkulasi)

---

## 📁 FASE 4: File Uploader Dinamis & Finalisasi

### 4.1 Upload Bukti Dokumen (Repeater)
- `[x]` `./artisan make:controller EvidenceController`
- `[x]` `POST /api/v1/evidences/upload` — endpoint `multipart/form-data`:
  - `[x]` Validasi format file: hanya `.pdf`, `.jpg`, `.png`
  - `[x]` Validasi ukuran: maks **2 MB** per file
  - `[x]` Simpan file ke Goravel Local Storage (`storage/evidences/` via `public_dir` disk)
  - `[x]` Simpan record ke tabel `assessment_evidences`
  - `[x]` Response: `evidence_id`, `document_name`, `file_url`
- `[x]` `DELETE /api/v1/evidences/{evidence_id}` — hapus dokumen bukti

### 4.2 Finalisasi & Submit Assessment
- `[x]` `POST /api/v1/assessments/submit` — khusus `ADMIN_KAMPUS`:
  - `[x]` Validasi status dan kepemilikan kampus
  - `[x]` Ubah status `campus_assessments` → `'SUBMITTED'` (data terkunci, tidak bisa diubah)
  - `[x]` Response: `status`, `final_overall_score`

---

## ⚠️ FASE 5: Error Handling & Standar Response JSON

- `[x]` Implementasi **5 skenario error standar** sesuai SDD:
  - `[x]` `422 Unprocessable Entity` — error validasi form (dengan field errors di Request binder)
  - `[x]` `403 Forbidden` — akses RBAC ditolak (role tidak sesuai kategori di middleware)
  - `[x]` `401 Unauthorized` — token JWT tidak valid / expired di middleware
  - `[x]` `404 Not Found` — data/record tidak ditemukan (ditangani global Fallback)
  - `[x]` `500 Internal Server Error` — error server/database internal (ditangani global Recover)
- `[x]` Buat response helper JSON standar:
  ```go
  // Success: {"status": "success", "message": "...", "data": {...}}
  // Error:   {"status": "error", "code": 4XX, "message": "...", "errors": {...}}
  ```
- `[x]` Lakukan standardisasi format JSON response (snake_case) dengan `json:` tags di semua model
- `[x]` Tambahkan `evidences` array pada response GetIndicatorsByCategory

---

## 📋 Ringkasan Progress

| Fase | Deskripsi | Status |
| :--- | :--- | :---: |
| **Fase 1** | Setup Framework, DB Migration, Model, Seeder, JWT + RBAC | `[x]` |
| **Fase 2** | Master Data (Kampus, User) + Dashboard API | `[x]` |
| **Fase 3** | Scoring Engine + Form Assessment API | `[x]` |
| **Fase 4** | File Uploader Dinamis + Submit & Finalisasi | `[x]` |
| **Fase 5** | Error Handling & Standar Response JSON | `[x]` |

> **Total Indikator yang Perlu Didukung Scoring Engine:**
> SI(8) + EC(10) + WS(6) + WR(6) + TR(8) + ED(10) + GD(12) = **60 Indikator**
