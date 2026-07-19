# 📘 Panduan Penggunaan REST API — UI GreenMetric

Dokumen ini berisi panduan lengkap penggunaan seluruh REST API, daftar parameter yang **Wajib (Required)** dan **Opsional (Optional)**, serta contoh payload Request & Response JSON untuk mempermudah integrasi Frontend (FE).

---

## 📌 Ketentuan Umum & Header
Setiap endpoint privat (membutuhkan login) wajib melampirkan token JWT di header HTTP:
*   **Key:** `Authorization`
*   **Value:** `Bearer <token_jwt_anda>` (tanpa tanda kurung siku)

---

## 🔑 1. Modul Autentikasi (Auth)

### A. Login User
*   **Endpoint:** `POST /api/v1/auth/login`
*   **Akses:** Publik (Tanpa Token)
*   **Payload Request (JSON):**
    *   `email` `[Wajib]` (string): Email terdaftar (contoh: `admin@polinema.ac.id`)
    *   `password` `[Wajib]` (string): Sandi minimal 6 karakter (contoh: `secretpassword`)

```json
{
  "email": "admin@polinema.ac.id",
  "password": "secretpassword"
}
```

*   **Response Sukses (200 OK):**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "campus_id": 1,
      "name": "Syarif Admin Green Campus",
      "email": "admin@polinema.ac.id",
      "role": "ADMIN_KAMPUS",
      "campus": {
        "id": 1,
        "code": "POLINEMA-KDR",
        "name": "Politeknik Negeri Malang - PSDKU Kediri",
        "institution_type": "Vocational",
        "climate": "Tropical",
        "setting": "Suburban"
      }
    }
  }
}
```

*   **Response Gagal (401 Unauthorized):**
```json
{
  "status": "error",
  "code": 401,
  "message": "invalid email or password"
}
```

---

## 🏫 2. Modul Manajemen Kampus (SUPER_ADMIN Only)

### A. Tambah Kampus Baru
*   **Endpoint:** `POST /api/v1/campuses`
*   **Akses:** Privat (Hanya `SUPER_ADMIN`)
*   **Payload Request (JSON):**
    *   `code` `[Wajib]` (string): Kode unik kampus (contoh: `POLINEMA-KDR`)
    *   `name` `[Wajib]` (string): Nama lengkap kampus
    *   `institution_type` `[Opsional]` (string): `Vocational` / `University`
    *   `climate` `[Opsional]` (string): `Tropical` / `Temperate`
    *   `setting` `[Opsional]` (string): `Urban` / `Suburban` / `Rural`

```json
{
  "code": "POLINEMA-KDR",
  "name": "Politeknik Negeri Malang - PSDKU Kediri",
  "institution_type": "Vocational",
  "climate": "Tropical",
  "setting": "Suburban"
}
```

---

## 👤 3. Modul Manajemen User / Operator (ADMIN_KAMPUS / SUPER_ADMIN)

### A. Tambah User Baru
*   **Endpoint:** `POST /api/v1/users`
*   **Akses:** Privat (`ADMIN_KAMPUS` atau `SUPER_ADMIN`)
*   **Payload Request (JSON):**
    *   `name` `[Wajib]` (string): Nama lengkap operator
    *   `email` `[Wajib]` (string): Email unik untuk login
    *   `password` `[Wajib]` (string): Sandi minimal 6 karakter
    *   `role` `[Wajib]` (string): Pilihan: `ADMIN_KAMPUS`, `OPERATOR_SI`, `OPERATOR_EC`, `OPERATOR_WS`, `OPERATOR_WR`, `OPERATOR_TR`, `OPERATOR_ED`, `OPERATOR_GD`
    *   `campus_id` `[Kondisional/Wajib jika SUPER_ADMIN]` (integer): Wajib dikirim hanya jika yang membuat adalah `SUPER_ADMIN`. Jika pembuatnya `ADMIN_KAMPUS`, field ini opsional dan otomatis memetakan ke kampus admin yang login.

```json
{
  "name": "Operator Green Infrastructure",
  "email": "op.si@polinema.ac.id",
  "password": "operatorpassword",
  "role": "OPERATOR_SI"
}
```

---

## 🛠️ 3.1 Modul Manajemen Soal & Field Dinamis (Google Form-Style) — SUPER_ADMIN ONLY

Modul ini memungkinkan Super Admin untuk membuat, melihat, memperbarui, dan menghapus soal (indikator) beserta jenis input jawabannya secara dinamis.

### A. Tambah Soal & Dynamic Fields Baru
*   **Endpoint:** `POST /api/v1/admin/indicators`
*   **Akses:** Privat (Hanya `SUPER_ADMIN`)
*   **Payload Request (JSON):**
    *   `category_id` `[Wajib]` (integer): ID kategori (1 s/d 7)
    *   `code` `[Wajib]` (string): Kode unik baru (contoh: `SI9`)
    *   `title` `[Wajib]` (string): Pertanyaan/deskripsi soal
    *   `input_type` `[Wajib]` (string): `NUMERIC_FORMULA` / `SINGLE_CHOICE`
    *   `max_points` `[Wajib]` (integer): Bobot poin maksimal (contoh: `100`)
    *   `fields` `[Wajib]` (array of objects): Struktur kolom input jawaban yang harus diisi operator. Setiap objek berisi:
        *   `key` (string): Nama kunci/parameter di payload submit (contoh: `luas_total`)
        *   `label` (string): Label yang ditampilkan ke operator (contoh: `Total Luas Kampus`)
        *   `type` (string): Tipe data (`int`, `float`, `date`, `varchar`, `choice`)
        *   `required` (boolean): Apakah wajib diisi?

```json
{
  "category_id": 1,
  "code": "SI9",
  "title": "Jumlah kebun hidroponik di dalam kampus",
  "input_type": "NUMERIC_FORMULA",
  "max_points": 100,
  "fields": [
    {
      "key": "jumlah_kebun",
      "label": "Jumlah Kebun Hidroponik Aktif",
      "type": "int",
      "required": true
    }
  ]
}
```

### B. List Semua Indikator & Fields
*   **Endpoint:** `GET /api/v1/admin/indicators`
*   **Akses:** Privat (Hanya `SUPER_ADMIN`)
*   **Response Sukses (200 OK):**
```json
{
  "status": "success",
  "message": "Indicators retrieved successfully",
  "data": [
    {
      "id": 61,
      "category_id": 1,
      "code": "SI9",
      "title": "Jumlah kebun hidroponik di dalam kampus",
      "input_type": "NUMERIC_FORMULA",
      "max_points": 100,
      "fields": [
        {
          "id": 150,
          "indicator_id": 61,
          "key": "jumlah_kebun",
          "label": "Jumlah Kebun Hidroponik Aktif",
          "type": "int",
          "required": true
        }
      ]
    }
  ]
}
```

### C. Update Soal & Fields
*   **Endpoint:** `PUT /api/v1/admin/indicators/{id}`
*   **Akses:** Privat (Hanya `SUPER_ADMIN`)
*   **Payload:** Sama dengan `POST` (semua field opsional, hanya kirim yang ingin diubah). Jika `fields` dikirim, backend akan menimpa seluruh field lama dengan yang baru.

### D. Hapus Soal
*   **Endpoint:** `DELETE /api/v1/admin/indicators/{id}`
*   **Akses:** Privat (Hanya `SUPER_ADMIN`)

---

## 🎨 3.2 Panduan Render Form Dinamis (Untuk Frontend)

Dengan sistem Form Dinamis ini, Frontend **tidak boleh** melakukan *hardcode* komponen input form per indikator. Frontend harus me-loop array `fields` yang dikembalikan dari API `/categories/{category_code}/indicators` secara dinamis.

### Alur Rendering Dinamis:
1.  **Dapatkan Indikator & Fields** dari API `/api/v1/categories/{category_code}/indicators`.
2.  Iterasi list `indicators`. Untuk setiap indikator, tampilkan judul (`title`) dan lakukan iterasi atas array `fields`:
    *   Jika `field.type` == `"int"` atau `"float"`:
        Render input bertipe angka: `<input type="number" step="any" name="field.key" />`.
    *   Jika `field.type` == `"date"`:
        Render input kalender: `<input type="date" name="field.key" />`.
    *   Jika `field.type` == `"varchar"`:
        Render input teks biasa: `<input type="text" name="field.key" />`.
    *   Jika `field.type` == `"choice"`:
        Render komponen pilihan (`<select>` atau Radio Group). Pilihan opsi didapat secara dinamis dari array `tiers` pada indikator bersangkutan (loop `tier.option_label` sebagai label dan value-nya).

### Contoh Implementasi Singkat (Vue 3 / React Concept):
```html
<!-- Vue 3 Template Dynamic Form -->
<div v-for="indicator in indicators" :key="indicator.id" class="question-card">
  <h3>{{ indicator.code }} - {{ indicator.title }}</h3>
  
  <div v-for="field in indicator.fields" :key="field.id" class="input-group">
    <label>{{ field.label }}</label>
    
    <!-- Render jika tipe angka -->
    <input v-if="field.type === 'int' || field.type === 'float'" 
           type="number" step="any" :required="field.required"
           v-model="answers[indicator.code][field.key]" />
           
    <!-- Render jika tipe tanggal -->
    <input v-else-if="field.type === 'date'" 
           type="date" :required="field.required"
           v-model="answers[indicator.code][field.key]" />
           
    <!-- Render jika tipe teks -->
    <input v-else-if="field.type === 'varchar'" 
           type="text" :required="field.required"
           v-model="answers[indicator.code][field.key]" />
           
    <!-- Render jika tipe pilihan ganda -->
    <select v-else-if="field.type === 'choice'" :required="field.required"
            v-model="answers[indicator.code][field.key]">
      <option v-for="tier in indicator.tiers" :key="tier.id" :value="tier.option_label">
        {{ tier.option_label }}
      </option>
    </select>
  </div>
</div>
```

---

## 📊 4. Modul Dashboard & Formulir Kategori

### A. Ambil Summary Dashboard
*   **Endpoint:** `GET /api/v1/assessments/dashboard`
*   **Akses:** Privat (Seluruh Role Kampus)
*   **Payload Request:** Kosong (`GET`)
*   **Response Sukses (200 OK):**
```json
{
  "status": "success",
  "message": "Dashboard statistics loaded successfully",
  "data": {
    "campus_name": "Politeknik Negeri Malang - PSDKU Kediri",
    "current_year": 2026,
    "assessment_status": "DRAFT",
    "overall_score": 6450.50,
    "max_overall_score": 10000,
    "estimated_rank": 1,
    "category_breakdown": [
      {
        "category_code": "SI",
        "category_name": "Setting and Infrastructure",
        "earned_points": 850.00,
        "max_points": 1100
      }
    ],
    "trend_history": [
      {
        "year": 2026,
        "score": 6450.50
      }
    ]
  }
}
```

---

## 🧮 5. Modul Core Scoring Engine

### A. Simpan & Hitung Jawaban Form
*   **Endpoint:** `POST /api/v1/assessments/answers`
*   **Akses:** Privat (`ADMIN_KAMPUS` atau `OPERATOR_<KATEGORI>` terkait)
*   **Payload Request (JSON):**
    *   `indicator_code` `[Wajib]` (string): Kode indikator (contoh: `SI1`, `EC3`, dll.)
    *   `assessment_year` `[Wajib]` (integer): Tahun berjalan (contoh: `2026`)
    *   `raw_input_data` `[Wajib]` (object/JSON): Data masukan form yang bervariasi nilainya.

#### 💡 Spesifikasi Pengisian `raw_input_data` Berdasarkan Tipe Indikator:

#### 1. Jika Tipe Indikator `SINGLE_CHOICE` (Pilihan Opsi/Ganda)
Hanya membutuhkan satu parameter kunci: `option_label`.
*   **Contoh Payload Request:**
```json
{
  "indicator_code": "EC3",
  "assessment_year": 2026,
  "raw_input_data": {
    "option_label": ">= 3 sumber energi terbarukan"
  }
}
```

#### 2. Jika Tipe Indikator `NUMERIC_FORMULA` (Kalkulasi Rumus)
Wajib mengirimkan parameter angka spesifik sesuai rumus kode indikator berikut:

| Kode Indikator | Parameter Wajib di dalam `raw_input_data` | Deskripsi |
| :--- | :--- | :--- |
| **`SI1`** | `"luas_total"` `(float)`, `"luas_dasar"` `(float)` | Luas ruang terbuka |
| **`SI2`** | `"luas_hutan"` `(float)`, `"luas_total"` `(float)` | Luas area riset hutan |
| **`SI3`** | `"luas_vegetasi"` `(float)`, `"luas_total"` `(float)` | Luas vegetasi tanam |
| **`SI4`** | `"luas_total"` `(float)`, `"luas_dasar"` `(float)`, `"populasi"` `(float)` | Ruang terbuka per orang |
| **`EC1`** | `"persentase_alat_hemat_energi"` `(float)` | Persentase alat hemat energi |
| **`EC2`** | `"luas_smart_building"` `(float)`, `"luas_total_bangunan"` `(float)` | Luas smart building |
| **`EC4`** | `"total_listrik"` `(float)`, `"populasi"` `(float)` | KWh listrik per orang |
| **`EC5`** | `"produksi_energi_terbarukan"` `(float)`, `"total_penggunaan_energi"` `(float)` | Rasio energi terbarukan |
| **`EC8`** | `"jejak_karbon"` `(float)`, `"populasi"` `(float)` | Jejak karbon per orang |
| **`WR1`** | `"area_resapan"` `(float)`, `"luas_total"` `(float)` | Luas area resapan air |
| **`WR5`** | `"air_olahan_dikonsumsi"` `(float)`, `"total_air_dikonsumsi"` `(float)` | Rasio konsumsi air olahan |
| **`TR1`** | `"total_kendaraan"` `(float)`, `"populasi"` `(float)` | Rasio kendaraan per orang |
| **`TR4`** | `"total_zev"` `(float)`, `"populasi"` `(float)` | Rasio kendaraan ramah lingkungan |
| **`TR5`** | `"luas_parkir"` `(float)`, `"luas_total"` `(float)` | Rasio area parkir permukaan |
| **`ED1`** | `"mk_keberlanjutan"` `(float)`, `"total_mk"` `(float)` | Rasio mata kuliah hijau |
| **`ED2`** | `"dana_riset_keberlanjutan"` `(float)`, `"total_dana_riset"` `(float)` | Rasio dana penelitian hijau |
| **`ED3`** | `"publikasi_keberlanjutan"` `(float)`, `"total_publikasi"` `(float)` | Rasio publikasi ilmiah hijau |
| **`ED10`**| `"lulusan_green_jobs"` `(float)`, `"total_lulusan"` `(float)` | Rasio lulusan yang bekerja hijau |
| **`GD1`** | `"anggaran_keberlanjutan"` `(float)`, `"total_anggaran"` `(float)` | Rasio anggaran hijau universitas |
| **`GD8`** | `"pimpinan_perempuan"` `(float)`, `"total_pimpinan"` `(float)` | Rasio pimpinan perempuan |

*   **Contoh Payload Request `SI1`:**
```json
{
  "indicator_code": "SI1",
  "assessment_year": 2026,
  "raw_input_data": {
    "luas_total": 100000,
    "luas_dasar": 15000
  }
}
```

*   **Response Sukses Hasil Perhitungan (200 OK):**
```json
{
  "status": "success",
  "message": "Answer saved and calculated successfully",
  "data": {
    "earned_points": 100,
    "calculated_value": 85,
    "overall_score": 6450.50
  }
}
```

---

## 📁 6. Modul Upload Bukti Dinamis (Evidence)

### A. Unggah Bukti
*   **Endpoint:** `POST /api/v1/evidences/upload`
*   **Content-Type:** `multipart/form-data`
*   **Akses:** Privat (`ADMIN_KAMPUS` / `OPERATOR_<KATEGORI>`)
*   **Payload Request (Form-Data):**
    *   `indicator_code` `[Wajib]` (string): Contoh `SI1`
    *   `assessment_year` `[Wajib]` (string/integer): Contoh `2026`
    *   `document_name` `[Wajib]` (string): Contoh `SK Lahan Terbuka`
    *   `file` `[Wajib]` (file binary): File PDF/PNG/JPG (Maks 2MB)
    *   `description` `[Opsional]` (string): Catatan tambahan mengenai dokumen

*   **Response Sukses (200 OK):**
```json
{
  "status": "success",
  "message": "Evidence uploaded successfully",
  "data": {
    "id": 1,
    "assessment_answer_id": 12,
    "document_name": "SK Lahan Terbuka",
    "description": "SK resmi luas lahan",
    "file_url": "http://localhost:3030/public/evidences/file-17192837.pdf"
  }
}
```

---

## 🔒 7. Modul Finalisasi (Submit)

### A. Kunci Data Evaluasi
*   **Endpoint:** `POST /api/v1/assessments/submit`
*   **Akses:** Privat (Khusus `ADMIN_KAMPUS` atau `SUPER_ADMIN`)
*   **Payload Request (JSON):**
    *   `assessment_year` `[Wajib]` (integer): Tahun evaluasi yang ingin dikunci (contoh: `2026`)
    *   `campus_id` `[Kondisional/Wajib jika SUPER_ADMIN]` (integer): ID Kampus yang ingin dikunci (hanya wajib dikirim jika diakses oleh `SUPER_ADMIN`).

```json
{
  "assessment_year": 2026
}
```

---

## 📋 8. Katalog Lengkap Indikator & Kategori (Referensi Backend & Frontend)

Berikut adalah daftar lengkap 60 indikator UI GreenMetric yang didukung oleh sistem, terbagi ke dalam 7 kategori penilaian:

### 1. Setting and Infrastructure (SI) — Bobot 11% (Maks 1.100 Poin)
| Kode | Nama Indikator | Tipe Input | Maks Poin |
| :--- | :--- | :--- | :--- |
| **`SI1`** | Rasio luas ruang terbuka terhadap total luas | `NUMERIC_FORMULA` | 200 |
| **`SI2`** | Total luas hutan riset/masyarakat terhadap luas total | `NUMERIC_FORMULA` | 100 |
| **`SI3`** | Total luas vegetasi tanam terhadap luas total | `NUMERIC_FORMULA` | 200 |
| **`SI4`** | Luas ruang terbuka per orang | `NUMERIC_FORMULA` | 200 |
| **`SI5`** | Fasilitas penunjang disabilitas/kebutuhan khusus | `SINGLE_CHOICE` | 100 |
| **`SI6`** | Fasilitas keamanan dan keselamatan kampus | `SINGLE_CHOICE` | 100 |
| **`SI7`** | Infrastruktur kesehatan mahasiswa & staf | `SINGLE_CHOICE` | 100 |
| **`SI8`** | Konservasi flora, fauna, dan satwa liar | `SINGLE_CHOICE` | 100 |

### 2. Energy and Climate Change (EC) — Bobot 20% (Maks 2.000 Poin)
| Kode | Nama Indikator | Tipe Input | Maks Poin |
| :--- | :--- | :--- | :--- |
| **`EC1`** | Penggunaan peralatan hemat energi | `NUMERIC_FORMULA` | 200 |
| **`EC2`** | Implementasi smart building | `NUMERIC_FORMULA` | 300 |
| **`EC3`** | Jumlah sumber energi terbarukan di kampus | `SINGLE_CHOICE` | 300 |
| **`EC4`** | Total penggunaan listrik per orang (kWh/tahun) | `NUMERIC_FORMULA` | 200 |
| **`EC5`** | Rasio produksi energi terbarukan terhadap total energi | `NUMERIC_FORMULA` | 200 |
| **`EC6`** | Elemen green building yang diterapkan | `SINGLE_CHOICE` | 200 |
| **`EC7`** | Program pengurangan emisi gas rumah kaca (GHG) | `SINGLE_CHOICE` | 200 |
| **`EC8`** | Total jejak karbon per orang (metrik ton/tahun) | `NUMERIC_FORMULA` | 200 |
| **`EC9`** | Jumlah program inovatif energi & iklim | `SINGLE_CHOICE` | 100 |
| **`EC10`**| Program universitas berdampak perubahan iklim | `SINGLE_CHOICE` | 100 |

### 3. Waste (WS) — Bobot 17% (Maks 1.700 Poin)
| Kode | Nama Indikator | Tipe Input | Maks Poin |
| :--- | :--- | :--- | :--- |
| **`WS1`** | Program 3R (Reduce, Reuse, Recycle) | `SINGLE_CHOICE` | 200 |
| **`WS2`** | Program pengurangan kertas dan plastik | `SINGLE_CHOICE` | 300 |
| **`WS3`** | Pengolahan sampah organik | `SINGLE_CHOICE` | 300 |
| **`WS4`** | Pengolahan sampah anorganik | `SINGLE_CHOICE` | 300 |
| **`WS5`** | Pengolahan sampah beracun & berbahaya (B3) | `SINGLE_CHOICE` | 300 |
| **`WS6`** | Pengolahan air limbah (sewage) | `SINGLE_CHOICE` | 300 |

### 4. Water (WR) — Bobot 11% (Maks 1.100 Poin)
| Kode | Nama Indikator | Tipe Input | Maks Poin |
| :--- | :--- | :--- | :--- |
| **`WR1`** | Persentase area resapan air | `NUMERIC_FORMULA` | 100 |
| **`WR2`** | Program konservasi air | `SINGLE_CHOICE` | 200 |
| **`WR3`** | Implementasi program daur ulang air | `SINGLE_CHOICE` | 200 |
| **`WR4`** | Penggunaan peralatan hemat air | `SINGLE_CHOICE` | 200 |
| **`WR5`** | Rasio konsumsi air olahan | `NUMERIC_FORMULA` | 200 |
| **`WR6`** | Pengendalian pencemaran air di kampus | `SINGLE_CHOICE` | 200 |

### 5. Transportation (TR) — Bobot 17% (Maks 1.700 Poin)
| Kode | Nama Indikator | Tipe Input | Maks Poin |
| :--- | :--- | :--- | :--- |
| **`TR1`** | Rasio kendaraan bermesin pembakaran per orang | `NUMERIC_FORMULA` | 200 |
| **`TR2`** | Layanan shuttle bus kampus | `SINGLE_CHOICE` | 250 |
| **`TR3`** | Ketersediaan kendaraan bebas emisi (ZEV) | `SINGLE_CHOICE` | 200 |
| **`TR4`** | Rasio ZEV terhadap populasi kampus | `NUMERIC_FORMULA` | 200 |
| **`TR5`** | Rasio luas area parkir permukaan terhadap luas kampus | `NUMERIC_FORMULA` | 200 |
| **`TR6`** | Inisiatif pengurangan area parkir (3 tahun terakhir) | `SINGLE_CHOICE` | 200 |
| **`TR7`** | Inisiatif pembatasan kendaraan pribadi | `SINGLE_CHOICE` | 200 |
| **`TR8`** | Kualitas jalur pejalan kaki & pesepeda | `SINGLE_CHOICE` | 250 |

### 6. Education and Research (ED) — Bobot 13% (Maks 1.300 Poin)
| Kode | Nama Indikator | Tipe Input | Maks Poin |
| :--- | :--- | :--- | :--- |
| **`ED1`** | Rasio mata kuliah keberlanjutan terhadap total MK | `NUMERIC_FORMULA` | 200 |
| **`ED2`** | Rasio dana penelitian keberlanjutan | `NUMERIC_FORMULA` | 200 |
| **`ED3`** | Rasio publikasi ilmiah bertema keberlanjutan | `NUMERIC_FORMULA` | 200 |
| **`ED4`** | Jumlah kegiatan/event bertema keberlanjutan | `SINGLE_CHOICE` | 100 |
| **`ED5`** | Kegiatan organisasi mahasiswa terkait lingkungan | `SINGLE_CHOICE` | 150 |
| **`ED6`** | Kegiatan pelestarian budaya di kampus | `SINGLE_CHOICE` | 100 |
| **`ED7`** | Kolaborasi keberlanjutan tingkat internasional | `SINGLE_CHOICE` | 100 |
| **`ED8`** | Program pengabdian masyarakat bertema lingkungan | `SINGLE_CHOICE` | 100 |
| **`ED9`** | Jumlah start-up/inkubasi bisnis hijau | `SINGLE_CHOICE` | 100 |
| **`ED10`**| Persentase alumni yang bekerja di bidang hijau (Green Jobs) | `NUMERIC_FORMULA` | 50 |

### 7. Governance and Digitalization (GD) — Bobot 11% (Maks 1.100 Poin)
| Kode | Nama Indikator | Tipe Input | Maks Poin |
| :--- | :--- | :--- | :--- |
| **`GD1`** | Persentase anggaran universitas untuk keberlanjutan | `NUMERIC_FORMULA` | 200 |
| **`GD2`** | Situs web khusus informasi keberlanjutan kampus | `SINGLE_CHOICE` | 200 |
| **`GD3`** | Penerbitan Laporan Berkelanjutan (Sustainability Report) | `SINGLE_CHOICE` | 100 |
| **`GD4`** | Laporan Keuangan audited yang dipublikasikan | `SINGLE_CHOICE` | 100 |
| **`GD5`** | Unit kerja/kantor pengelola keberlanjutan | `SINGLE_CHOICE` | 100 |
| **`GD6`** | Penggunaan TIK untuk perencanaan & evaluasi hijau | `SINGLE_CHOICE` | 50 |
| **`GD7`** | Kebijakan adopsi teknologi cerdas (AI/IoT) | `SINGLE_CHOICE` | 50 |
| **`GD8`** | Rasio pemimpin perempuan terhadap total pimpinan | `NUMERIC_FORMULA` | 100 |
| **`GD9`** | Penerapan sistem integritas & anti-korupsi | `SINGLE_CHOICE` | 50 |
| **`GD10`**| Ketersediaan Whistleblowing System (pengaduan) | `SINGLE_CHOICE` | 50 |
| **`GD11`**| Program literasi digital berbasis LMS | `SINGLE_CHOICE` | 50 |
| **`GD12`**| Kode etik tertulis sivitas akademika | `SINGLE_CHOICE` | 50 |

```
