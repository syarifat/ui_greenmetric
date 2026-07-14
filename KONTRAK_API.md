# 🔌 Kontrak REST API — UI GreenMetric Self-Assessment System

Dokumen ini berisi spesifikasi teknis endpoint REST API backend untuk digunakan sebagai panduan integrasi oleh Frontend (FE) Developer.

---

## 📌 Ketentuan Umum

*   **Base URL:** `http://localhost:3000` (atau sesuaikan dengan konfigurasi `.env` host)
*   **Content-Type:** `application/json` (kecuali untuk upload file yang menggunakan `multipart/form-data`)
*   **Autentikasi:** Menggunakan **JWT Bearer Token**. Setiap request ke endpoint privat wajib menyertakan token di Header HTTP:
    ```http
    Authorization: Bearer <your_jwt_token>
    ```

---

## 🛑 Format Standar Response

### 1. Response Sukses (200 OK / 201 Created)
```json
{
  "status": "success",
  "message": "Deskripsi pesan sukses",
  "data": { ... } // Atau berupa array [ ... ] jika data list
}
```

### 2. Response Gagal / Error

#### A. 401 Unauthorized (Token kadaluwarsa atau tidak valid)
```json
{
  "status": "error",
  "code": 401,
  "message": "Token tidak valid atau kedaluwarsa"
}
```

#### B. 403 Forbidden (Akses peran/RBAC ditolak)
```json
{
  "status": "error",
  "code": 403,
  "message": "Anda tidak memiliki hak akses"
}
```

#### C. 404 Not Found (Endpoint salah atau data tidak ada)
```json
{
  "status": "error",
  "code": 404,
  "message": "Resource tidak ditemukan"
}
```

#### D. 422 Unprocessable Entity (Validasi Input Gagal)
```json
{
  "status": "error",
  "code": 422,
  "message": "Validasi gagal",
  "errors": {
    "email": "The email field is required.",
    "password": "The password field is required."
  }
}
```

#### E. 500 Internal Server Error (Kendala Internal Server/DB)
```json
{
  "status": "error",
  "code": 500,
  "message": "Terjadi kesalahan internal pada server"
}
```

---

## 🔑 1. Modul Autentikasi

### A. Login User
*   **Endpoint:** `POST /api/v1/auth/login`
*   **Akses:** Publik (Tanpa Token)
*   **Payload Request (JSON):**
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

### B. Logout User
*   **Endpoint:** `POST /api/v1/auth/logout`
*   **Akses:** Privat (Butuh Token)
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Logout successful"
    }
    ```

---

## 🏫 2. Modul Manajemen Kampus (SUPER_ADMIN)

### A. Ambil Semua Kampus
*   **Endpoint:** `GET /api/v1/campuses`
*   **Akses:** Privat (`SUPER_ADMIN` saja)
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Campuses retrieved successfully",
      "data": [
        {
          "id": 1,
          "code": "POLINEMA-KDR",
          "name": "Polinema Kampus Kediri",
          "institution_type": "Politeknik",
          "climate": "Tropis",
          "setting": "Sub-Urban"
        }
      ]
    }
    ```

### B. Tambah Kampus Baru
*   **Endpoint:** `POST /api/v1/campuses`
*   **Akses:** Privat (`SUPER_ADMIN` saja)
*   **Payload Request (JSON):**
    ```json
    {
      "code": "ITS",
      "name": "Institut Teknologi Sepuluh Nopember",
      "institution_type": "Universitas",
      "climate": "Tropis",
      "setting": "Urban"
    }
    ```
*   **Response Sukses (201 Created):**
    ```json
    {
      "status": "success",
      "message": "Campus created successfully",
      "data": {
        "id": 2,
        "code": "ITS",
        "name": "Institut Teknologi Sepuluh Nopember",
        "institution_type": "Universitas",
        "climate": "Tropis",
        "setting": "Urban"
      }
    }
    ```

### C. Update Data Kampus
*   **Endpoint:** `PUT /api/v1/campuses/{id}`
*   **Akses:** Privat (`SUPER_ADMIN` saja)
*   **Payload Request (JSON):** *(Kirim field yang ingin diubah saja)*
    ```json
    {
      "name": "ITS Surabaya",
      "setting": "Urban Pantai"
    }
    ```
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Campus updated successfully",
      "data": {
        "id": 2,
        "code": "ITS",
        "name": "ITS Surabaya",
        "institution_type": "Universitas",
        "climate": "Tropis",
        "setting": "Urban Pantai"
      }
    }
    ```

### D. Hapus Kampus
*   **Endpoint:** `DELETE /api/v1/campuses/{id}`
*   **Akses:** Privat (`SUPER_ADMIN` saja)
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Campus deleted successfully"
    }
    ```

---

## 👥 3. Modul Manajemen User/Operator Kampus

*   *Catatan:* `SUPER_ADMIN` bisa mengelola user seluruh kampus. `ADMIN_KAMPUS` hanya bisa mengelola operator di naungan kampusnya sendiri.

### A. Ambil Semua User/Operator
*   **Endpoint:** `GET /api/v1/users`
*   **Akses:** Privat (`ADMIN_KAMPUS` / `SUPER_ADMIN`)
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Users retrieved successfully",
      "data": [
        {
          "id": 2,
          "campus_id": 1,
          "name": "Operator Green Metric SI",
          "email": "op_si@polinema.ac.id",
          "role": "OPERATOR_SI"
        }
      ]
    }
    ```

### B. Tambah User/Operator Baru
*   **Endpoint:** `POST /api/v1/users`
*   **Akses:** Privat (`ADMIN_KAMPUS` / `SUPER_ADMIN`)
*   **Payload Request (JSON):**
    *   *Peran Opsional:* `ADMIN_KAMPUS`, `OPERATOR_SI`, `OPERATOR_EC`, `OPERATOR_WS`, `OPERATOR_WR`, `OPERATOR_TR`, `OPERATOR_ED`, `OPERATOR_GD`.
    *   *Khusus `SUPER_ADMIN` wajib menyertakan field `"campus_id"`.*
    ```json
    {
      "name": "Operator EC Polinema",
      "email": "op_ec@polinema.ac.id",
      "password": "operatorpassword",
      "role": "OPERATOR_EC"
    }
    ```
*   **Response Sukses (201 Created):**
    ```json
    {
      "status": "success",
      "message": "User created successfully",
      "data": {
        "id": 3,
        "campus_id": 1,
        "name": "Operator EC Polinema",
        "email": "op_ec@polinema.ac.id",
        "role": "OPERATOR_EC"
      }
    }
    ```

### C. Update Data User
*   **Endpoint:** `PUT /api/v1/users/{id}`
*   **Akses:** Privat (`ADMIN_KAMPUS` / `SUPER_ADMIN`)
*   **Payload Request (JSON):** *(Kirim field yang diubah saja)*
    ```json
    {
      "name": "Operator EC Edit",
      "password": "newpassword123"
    }
    ```
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "User updated successfully",
      "data": {
        "id": 3,
        "campus_id": 1,
        "name": "Operator EC Edit",
        "email": "op_ec@polinema.ac.id",
        "role": "OPERATOR_EC"
      }
    }
    ```

### D. Hapus User
*   **Endpoint:** `DELETE /api/v1/users/{id}`
*   **Akses:** Privat (`ADMIN_KAMPUS` / `SUPER_ADMIN`)
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "User deleted successfully"
    }
    ```

---

## 📊 4. Modul Dashboard & Analisis Tren

### A. Get Ringkasan Dashboard Kampus
*   **Endpoint:** `GET /api/v1/assessments/dashboard`
*   **Akses:** Privat (Semua User Terotentikasi)
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Dashboard statistics loaded successfully",
      "data": {
        "campus_name": "Polinema Kampus Kediri",
        "current_year": 2026,
        "assessment_status": "DRAFT",
        "overall_score": 1450.00,
        "max_overall_score": 10000,
        "estimated_rank": 1,
        "category_breakdown": [
          {
            "category_code": "SI",
            "category_name": "Setting and Infrastructure",
            "earned_points": 300.00,
            "max_points": 1100,
            "weight_percentage": 11.00
          },
          {
            "category_code": "EC",
            "category_name": "Energy and Climate Change",
            "earned_points": 1150.00,
            "max_points": 2100,
            "weight_percentage": 21.00
          }
          // ... (total 7 kategori dasar)
        ],
        "trend_history": [
          {
            "year": 2025,
            "score": 4500.00
          },
          {
            "year": 2026,
            "score": 1450.00
          }
        ]
      }
    }
    ```

---

## 📝 5. Modul Soal & Pengisian Jawaban (Core)

### A. Ambil Indikator Soal & Jawaban Terakhir per Kategori
*   **Endpoint:** `GET /api/v1/categories/{category_code}/indicators`
    *   *Contoh URL:* `/api/v1/categories/si/indicators`
*   **Akses:** Privat (Terbuka untuk `ADMIN_KAMPUS`, atau operator kategori terkait seperti `OPERATOR_SI` hanya bisa mengakses kode `si`)
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Indicators and answers loaded successfully",
      "data": {
        "category_code": "SI",
        "category_name": "Setting and Infrastructure",
        "indicators": [
          {
            "id": 1,
            "code": "SI1",
            "title": "Rasio luas ruang terbuka terhadap total luas",
            "input_type": "NUMERIC_FORMULA",
            "max_points": 200,
            "tiers": [
              {
                "id": 1,
                "indicator_id": 1,
                "option_label": "<= 1%",
                "min_value": null,
                "max_value": 1.0000,
                "operator": "<=",
                "point_multiplier": 0.05
              },
              {
                "id": 2,
                "indicator_id": 1,
                "option_label": "> 1% - 80%",
                "min_value": 1.0000,
                "max_value": 80.0000,
                "operator": "BETWEEN",
                "point_multiplier": 0.25
              }
              // ... list threshold scoring tiers
            ],
            "answer": { // NULL jika belum diisi sama sekali
              "id": 12,
              "raw_input_data": "{\"luas_total\":100000,\"luas_dasar\":15000}",
              "calculated_value": 85.0000, // persentase / nilai terhitung
              "selected_tier_id": 3,
              "earned_points": 100.00,
              "evidences": [
                {
                  "id": 101,
                  "assessment_answer_id": 12,
                  "document_name": "Peta Ruang Terbuka Kampus",
                  "description": "Area hijau mencakup 85% dari total lahan",
                  "file_url": "http://localhost:3000/storage/evidences/peta_2026.pdf"
                }
              ]
            }
          }
        ]
      }
    }
    ```

### B. Simpan Jawaban Indikator Soal
*   **Endpoint:** `POST /api/v1/assessments/answers`
*   **Akses:** Privat (`ADMIN_KAMPUS` / `OPERATOR_<KAT>`)
*   **Payload Request Tipe NUMERIC_FORMULA (JSON):**
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
*   **Payload Request Tipe SINGLE_CHOICE (JSON):**
    ```json
    {
      "indicator_code": "SI5",
      "assessment_year": 2026,
      "raw_input_data": {
        "option_label": "Fasilitas tersedia sebagian dan sudah beroperasi"
      }
    }
    ```
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Answer saved and calculated successfully",
      "data": {
        "earned_points": 100.00,
        "calculated_value": 85.0000, // NULL jika bertipe SINGLE_CHOICE
        "overall_score": 1450.00     // Total skor baru kampus setelah ditambahkan jawaban ini
      }
    }
    ```

---

## 📁 6. Modul Upload Bukti & Kunci Pengisian (Finalisasi)

### A. Upload Dokumen Bukti Fisik
*   **Endpoint:** `POST /api/v1/evidences/upload`
*   **Content-Type:** `multipart/form-data`
*   **Akses:** Privat (`ADMIN_KAMPUS` / `OPERATOR_<KAT>`)
*   **Payload Request (Form Data):**
    - `file`: (File PDF/JPG/PNG, Maks 2MB)
    - `assessment_answer_id`: (ID record dari `AssessmentAnswer` bersangkutan)
    - `description`: (Teks deskripsi opsional)
*   **Response Sukses (201 Created):**
    ```json
    {
      "status": "success",
      "message": "Evidence document uploaded successfully",
      "data": {
        "id": 102,
        "assessment_answer_id": 12,
        "document_name": "peta_ruang_terbuka.pdf",
        "description": "Dokumen lampiran tata ruang terbuka hijau Polinema",
        "file_url": "/public/evidences/171563212456_hash.pdf"
      }
    }
    ```

### B. Hapus Dokumen Bukti Fisik
*   **Endpoint:** `DELETE /api/v1/evidences/{id}`
    *   *Contoh URL:* `/api/v1/evidences/102`
*   **Akses:** Privat (`ADMIN_KAMPUS` / `OPERATOR_<KAT>`)
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Evidence document deleted successfully"
    }
    ```

### C. Finalisasi & Submit Assessment (Kunci Data)
*   **Endpoint:** `POST /api/v1/assessments/submit`
*   **Akses:** Privat (`ADMIN_KAMPUS` saja)
*   **Payload Request (JSON):**
    ```json
    {
      "assessment_year": 2026
    }
    ```
*   **Response Sukses (200 OK):**
    ```json
    {
      "status": "success",
      "message": "Assessment finalized and locked successfully",
      "data": {
        "campus_assessment_id": 4,
        "status": "SUBMITTED",
        "overall_score": 1450.00
      }
    }
    ```
