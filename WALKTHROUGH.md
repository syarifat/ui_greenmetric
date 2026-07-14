# 📝 Walkthrough — FASE 1: Setup Framework & Autentikasi (RBAC)

Laporan ini mendokumentasikan rincian implementasi teknis pada **Fase 1** untuk proyek Sistem Web Self-Assessment UI GreenMetric.

---

## 1. Migrasi Database (8 Tabel Utama)
Telah dibuat dan dikonfigurasi 8 berkas migrasi database pada folder `database/migrations` dan didaftarkan pada **[migrations.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/bootstrap/migrations.go)**:

| No | Berkas Migrasi | Deskripsi & Skema Kolom |
|---|---|---|
| 1 | `create_campuses_table` | Menyimpan profil dasar universitas.<br>Kolom: `id`, `code` (unique), `name`, `institution_type`, `climate`, `setting`, `timestamps`. |
| 2 | `create_users_table` | Menyimpan data akun pengguna beserta relasinya.<br>Kolom: `id`, `campus_id` (foreign key), `name`, `email` (unique), `password`, `role`, `timestamps`. |
| 3 | `create_categories_table` | Menyimpan master 7 kategori utama.<br>Kolom: `id`, `code` (unique), `name`, `max_points`, `weight_percentage` (decimal 5,2), `timestamps`. |
| 4 | `create_indicators_table` | Menyimpan master indikator pertanyaan.<br>Kolom: `id`, `category_id` (foreign key), `code` (unique), `title` (length 500), `input_type` (enum: `NUMERIC_FORMULA`, `SINGLE_CHOICE`), `max_points`, `timestamps`. |
| 5 | `create_indicator_scoring_tiers_table` | Menyimpan parameter/logika skoring per indikator.<br>Kolom: `id`, `indicator_id` (foreign key), `option_label`, `min_value` & `max_value` (decimal 15,4), `operator` (enum: `<=`, `>`, `>=`, `BETWEEN`, `CHOICE`), `point_multiplier` (decimal 5,2), `timestamps`. |
| 6 | `create_campus_assessments_table` | Berkas header untuk merekam penilaian tahunan kampus.<br>Kolom: `id`, `campus_id` (foreign key), `assessment_year`, `overall_score` (decimal 10,2), `status` (enum: `DRAFT`, `SUBMITTED`, `VERIFIED`), `timestamps`. |
| 7 | `create_assessment_answers_table` | Merekam detail jawaban & perhitungan skor otomatis.<br>Kolom: `id`, `campus_assessment_id` (foreign key), `indicator_id` (foreign key), `raw_input_data` (json), `calculated_value` (decimal 15,4), `selected_tier_id`, `earned_points` (decimal 10,2), `timestamps`. |
| 8 | `create_assessment_evidences_table` | Merekam bukti fisik yang diunggah operator.<br>Kolom: `id`, `assessment_answer_id` (foreign key), `document_name`, `description` (text), `file_url` (length 500), `timestamps`. |

---

## 2. Model ORM (GORM)
Seluruh model GORM didefinisikan dalam folder `app/models` dengan relasi lengkap (*has-one / belongs-to*):

*   **[campus.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/campus.go):** Merepresentasikan data institusi kampus.
*   **[user.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/user.go):** Merepresentasikan pengguna dan menautkan relasi `Campus`.
*   **[category.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/category.go):** Berisi kode kategori (SI, EC, WS, WR, TR, ED, GD).
*   **[indicator.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/indicator.go):** Berisi data soal dan menautkan relasi `Category`.
*   **[indicator_scoring_tier.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/indicator_scoring_tier.go):** Berisi batas bawah/atas, operator, multiplier, dan menautkan relasi `Indicator`.
*   **[campus_assessment.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/campus_assessment.go):** Berisi status pengisian & skor total.
*   **[assessment_answer.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/assessment_answer.go):** Berisi data input dinamis JSON dan poin yang diperoleh.
*   **[assessment_evidence.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/assessment_evidence.go):** Menghubungkan berkas bukti lampiran ke jawaban assessment.

---

## 3. Seeders (Populasi Data Master Awal)
Telah diimplementasikan logika populating data master agar sistem siap diuji secara *clean*:

*   **[category_seeder.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/database/seeders/category_seeder.go):** Mengisi 7 kategori dasar berserta total bobot poin (skala total 10.000 poin).
*   **[indicator_seeder.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/database/seeders/indicator_seeder.go):** Memasukkan seluruh 60 kode soal dari kategori Infrastruktur (SI1-SI8) hingga Pendidikan (ED1-ED10) & Digitalisasi (GD1-GD12).
*   **[indicator_scoring_tier_seeder.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/database/seeders/indicator_scoring_tier_seeder.go):** Memasukkan tier skoring untuk indikator utama (seperti rasio ruang terbuka SI1/SI4, konsumsi listrik EC4, area resapan WR1, dan kepemilikan kendaraan TR1).
*   **[user_seeder.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/database/seeders/user_seeder.go):** Membuat satu kampus percobaan default `POLINEMA-KDR` dan melahirkan satu akun admin default `admin@polinema.ac.id` dengan kata sandi yang sudah di-hash (`secretpassword`).

---

## 4. Sistem Autentikasi (JWT) & Hak Akses (RBAC)
Sistem keamanan API dan otorisasi dibentuk melalui gabungan controller, service, dan middleware:

*   **[auth_service.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/services/auth_service.go):** Service untuk memverifikasi kecocokan sandi via hashing Goravel, lalu memanggil driver token generation `facades.Auth(ctx).Login(&user)`.
*   **[auth_controller.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/controllers/auth_controller.go):** Menyediakan handler endpoint `/login` dan `/logout` dengan pengembalian response standar REST API JSON.
*   **[jwt_middleware.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/middleware/jwt_middleware.go):** Memvalidasi keberadaan Bearer token di header `Authorization`. Otomatis melempar `401 Unauthorized` bila token rusak atau kedaluwarsa.
*   **[rbac_middleware.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/middleware/rbac_middleware.go):** Membatasi hak akses secara dinamis:
    *   `SUPER_ADMIN` memiliki izin tanpa batas.
    *   `ADMIN_KAMPUS` memiliki izin penuh kecuali CRUD data kampus itu sendiri.
    *   `OPERATOR_<KAT>` (contoh: `OPERATOR_SI`) hanya diperbolehkan melihat/mengisi form indikator yang berkode `SI`. Akses ke kategori lain dikembalikan sebagai `403 Forbidden`.
*   **[web.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/routes/web.go):** Mengatur penataan endpoint dengan pengelompokan API versi 1 (`/api/v1`) dan penyematan middleware pada jalur-jalur privat.

---

## 5. Modul Master Data & Dashboard (Fase 2)
Telah diimplementasikan pengendali CRUD dan statistik dashboard terintegrasi:

*   **[campus_controller.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/controllers/campus_controller.go):** 
    - CRUD data kampus (`GET`, `POST`, `PUT`, `DELETE`).
    - Dibatasi secara ketat oleh `RbacMiddleware` sehingga hanya dapat diakses oleh `SUPER_ADMIN`.
*   **[user_controller.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/controllers/user_controller.go):** 
    - CRUD operator kampus.
    - Dilengkapi proteksi multi-tenant: peran `ADMIN_KAMPUS` hanya dapat melihat, menambah, mengubah, dan menghapus akun operator yang berada di naungan kampusnya sendiri (`campus_id` yang sama).
*   **[dashboard_controller.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/controllers/dashboard_controller.go):** 
    - API statistik yang menyajikan: nama kampus, skor total saat ini, status pengisian, perkiraan peringkat (hasil urutan skor semua kampus pada tahun bersangkutan), rincian pencapaian poin per kategori (breakdown 7 kategori), dan riwayat tren tahunan kampus.
    - Kalkulasi breakdown dihitung secara in-memory untuk efisiensi dan kompabilitas lintas driver database.

---

## 6. Scoring Engine & Form Assessment (Fase 3)
Telah diimplementasikan mesin penilaian otomatis (Scoring Engine) dan form dinamis:

*   **[scoring_service.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/services/scoring_service.go):** 
    - Berisi mesin kalkulator terpadu yang mem-parsing input JSON.
    - Untuk tipe `NUMERIC_FORMULA`: mengevaluasi rumus matematika dari indikator (seperti rasio ruang terbuka `SI1`/`SI4`, konsumsi listrik `EC4`, area resapan `WR1`, dsb.), mencocokkan hasil persentase/nilai ke tier database (`<=`, `>`, `>=`, `BETWEEN`), menentukan multiplier pengali, dan menghitung skor akhir.
    - Untuk tipe `SINGLE_CHOICE`: mencocokkan input label pilihan ganda secara langsung dengan database seeder.
    - Otomatis memperbarui total skor keseluruhan kampus (`overall_score`) setiap ada jawaban baru yang tersimpan.
*   **[assessment_controller.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/controllers/assessment_controller.go):** 
    - `GET  /api/v1/categories/{category_code}/indicators`: Memuat daftar indikator pertanyaan beserta daftar tier skor pendukung dan riwayat jawaban kampus pada tahun berjalan untuk keperluan visualisasi formulir pengisian operator.
    - `POST /api/v1/assessments/answers`: Menerima input data mentah, memvalidasi form, melakukan kalkulasi skor dinamis, menyimpan detail jawaban, dan mengembalikan hasil skor yang diperoleh beserta total skor kampus terbaru.
*   **[save_answer_request.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/requests/save_answer_request.go):** 
    - Memvalidasi parameter masukan wajib: `indicator_code`, `assessment_year`, dan `raw_input_data`.

---

## 7. File Uploader & Finalisasi Assessment (Fase 4)
Telah diimplementasikan mekanisme uploader dokumen bukti fisik dan finalisasi penilaian tahunan:

*   **[evidence_controller.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/controllers/evidence_controller.go):** 
    - `POST /api/v1/evidences/upload`: Mengunggah berkas dokumen bukti fisik. Dilengkapi validasi ekstensi wajib (hanya PDF, JPG, JPEG, PNG) dan ukuran file maks 2MB. File disimpan langsung di dalam subfolder `/public/evidences/` dan dihubungkan ke record `AssessmentAnswer` melalui relasi foreign key database.
    - `DELETE /api/v1/evidences/{id}`: Menghapus data bukti dari database dan secara otomatis membuang file fisiknya dari folder local storage.
    - Dilengkapi proteksi multi-tenant kepemilikan kampus dan proteksi kunci (upload/delete diblokir jika status assessment bukan `'DRAFT'`).
*   **Finalisasi ([assessment_controller.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/controllers/assessment_controller.go)):**
    - `POST /api/v1/assessments/submit`: Mengunci pengisian data penilaian tahunan untuk kampus bersangkutan. Hak akses dibatasi hanya untuk `ADMIN_KAMPUS` (atau `SUPER_ADMIN`). Status diubah dari `'DRAFT'` menjadi `'SUBMITTED'`. Begitu dikunci, seluruh API mutasi (simpan jawaban, upload file, hapus file) diblokir secara permanen.

---

## 8. Global Error Handling & Standar Response (Fase 5)
Telah diselaraskan dan distandarisasi seluruh format response error JSON di tingkat aplikasi:

*   **Penyelarasan Middleware Keamanan:**
    - `401 Unauthorized` diselaraskan pada `JwtMiddleware` dan `RbacMiddleware` untuk menghasilkan pesan seragam: `{"status": "error", "code": 401, "message": "Token tidak valid atau kedaluwarsa"}`.
    - `403 Forbidden` diselaraskan pada `RbacMiddleware` untuk menghasilkan pesan seragam: `{"status": "error", "code": 403, "message": "Anda tidak memiliki hak akses"}`.
*   **Global 404 Fallback ([web.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/routes/web.go)):**
    - Menggunakan fungsi bawaan `facades.Route().Fallback(...)` untuk menangkap rute URL salah/tidak terdaftar dan menghasilkan response JSON terstandar: `{"status": "error", "code": 404, "message": "Resource tidak ditemukan"}`.
*   **Global 500 Panic Recovery ([web.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/routes/web.go)):**
    - Menggunakan fungsi bawaan `facades.Route().Recover(...)` untuk menangkap panic tak terduga (database mati, runtime crash, dsb.) dan mengembalikan response JSON aman tanpa membocorkan stack trace: `{"status": "error", "code": 500, "message": "Terjadi kesalahan internal pada server"}`.

---

### 8.1 Standar Format JSON Response & Relasi Evidences

Untuk memastikan konsistensi kontrak API dengan Frontend Developer, dilakukan standardisasi format JSON response:

*   **JSON Field Naming Convention (snake_case):**
    - Menambahkan tag `json:"field_name"` pada semua field di 8 model ORM ([campus.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/campus.go), [user.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/user.go), [category.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/category.go), [indicator.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/indicator.go), [indicator_scoring_tier.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/indicator_scoring_tier.go), [campus_assessment.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/campus_assessment.go), [assessment_answer.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/assessment_answer.go), [assessment_evidence.go](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/assessment_evidence.go))
    - Response JSON sekarang mengembalikan format snake_case (`code`, `name`, `institution_type`) bukan PascalCase (`Code`, `Name`, `InstitutionType`)
    - Field `Password` pada model User diberi tag `json:"-"` agar tidak pernah bocor ke response JSON

*   **Evidences Array pada Response Indicators:**
    - Menambahkan relasi `Evidences []AssessmentEvidence` pada model [AssessmentAnswer](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/models/assessment_answer.go)
    - Menggunakan `.With("Evidences")` pada query ORM di [AssessmentController.GetIndicatorsByCategory](file:///Users/syarifat/Data/my_project/UI%20GREENMETRIC/ui_greenmetric/app/http/controllers/assessment_controller.go) untuk preload data bukti
    - Response endpoint `GET /api/v1/categories/{category_code}/indicators` sekarang menyertakan array `evidences[]` di dalam setiap `answer` object sesuai kontrak API

*   **Nested Campus Object pada Login Response:**
    - Response endpoint `POST /api/v1/auth/login` sekarang menyertakan nested object `campus` di dalam `user` untuk memudahkan Frontend Developer mengakses data kampus tanpa request tambahan

---

## 9. Hasil Verifikasi Akhir
1.  **Kompilasi Kode:** Menjalankan `./artisan` berhasil tanpa ada kesalahan kompilasi (*compile error*).
2.  **Registrasi Layanan:** Seluruh endpoint baru (CRUD Kampus, CRUD User, API Dashboard, Form Assessment, Uploader Bukti Fisik, dan Finalisasi Submit) telah terdaftar di bawah grup `/api/v1` dengan proteksi middleware JWT & RBAC.
3.  **Keamanan Multi-tenant:** CRUD user dan dokumen bukti membatasi hak akses ADMIN_KAMPUS hanya pada ruang lingkup kampusnya sendiri secara aman.
4.  **Integritas Kunci Data:** Setiap API mutasi telah diverifikasi memblokir input jika status penilaian sudah final (`SUBMITTED`).
5.  **Integritas Skoring:** Logika evaluasi rumus dan perbandingan tier threshold telah diuji terkompilasi dengan andal.
6.  **Keseragaman Response JSON:** Seluruh kode HTTP Status (401, 403, 404, 422, dan 500) mengembalikan payload response error seragam dalam bahasa Indonesia sesuai spesifikasi SDD.