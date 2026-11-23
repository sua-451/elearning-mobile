# Dart & Flutter E-Learning Platform (Mobile Development)

Platform katalog dan baca materi pembelajaran mobile development (**Dart & Flutter**) yang sangat ringan, minimalis, dan cepat. Dibangun menggunakan arsitektur modern:
- **Backend**: Golang (Gin + GORM + MySQL + JWT)
- **Frontend**: Svelte (Vite + Tailwind CSS v4)
- **Database**: MySQL

---

## 📁 Struktur Folder Proyek

```text
e-learning-mobile/
├── backend/
│   ├── database/       # Koneksi MySQL, auto-migrate, & seeder otomatis
│   ├── handlers/       # REST API handler (Public & Admin CRUD)
│   ├── middleware/     # JWT authentication middleware
│   ├── models/         # Schema model (User, Material, Lesson)
│   ├── uploads/        # Direktori penyimpanan upload gambar materi
│   ├── .env.example    # Template konfigurasi environment backend
│   ├── go.mod / go.sum
│   └── main.go
│
├── frontend/
│   ├── src/
│   │   ├── lib/api.js  # Utilitas request API & session client
│   │   ├── pages/      # Home, MaterialDetail (Codelab Split-View), AdminLogin, AdminDashboard
│   │   ├── App.svelte  # Komponen utama & router hash ringkas
│   │   └── app.css     # Styling palet minimalis + Dart blue (#0175C2)
│   ├── package.json
│   └── vite.config.js
│
└── README.md
```

---

## 🚀 Fitur Unggulan

1. **Halaman Beranda (`/`):**
   - Header ringkas dan langsung menampilkan grid modul materi yang siap dipelajari.
2. **Halaman Baca Materi (`/materials/[slug]` - Gaya Codelab Split-View):**
   - **Sidebar Kiri (Collapsible):** Daftar bab (*lessons*) ringkas yang dapat dibuka/tutup untuk memaksimalkan ruang baca.
   - **Mode Standar:** Layout terpusat (`max-w-3xl`) dengan tipografi Markdown bersih jika simulasi tidak diaktifkan.
   - **Mode Simulasi Dart:** Layout **Split-Screen 2 Kolom** di desktop (kiri: materi & starter code dengan tombol *Copy Code*; kanan: *sticky* **Embedded DartPad** + tombol buka ke DartPad.dev di tab baru). Di perangkat mobile, otomatis berubah menjadi tab toggle fleksibel antara "Materi" dan "DartPad".
3. **Admin Dashboard (`/admin/dashboard`):**
   - Manajemen CRUD penuh untuk Modul Materi dan Bab Pelajaran.
   - Fitur upload file gambar lokal ke server backend atau input URL gambar langsung.
   - Toggle aktivasi DartPad dan input *Starter Code* per bab.
4. **Sistem Ringan & Bersih:**
   - Tanpa registrasi publik, tanpa kuis berlebih, tanpa chart analitik berat.

---

## 🔑 Kredensial Admin Default

Saat backend dijalankan pertama kali, sistem seeder otomatis membuat akun admin default berikut:

- **Username:** `admin`
- **Password:** `admin123`

---

## 🛠️ Panduan Menjalankan Sistem

### 1. Prasyarat
- **Go** (versi 1.20+)
- **Node.js** (versi 18+) & **npm**
- **MySQL Server** (XAMPP / Laragon / Docker / MySQL Standalone)

---

### 2. Konfigurasi & Menjalankan Backend (Golang)

1. Buat database baru di MySQL Anda, misalnya:
   ```sql
   CREATE DATABASE elearning_mobile;
   ```

2. Masuk ke direktori `backend/`:
   ```bash
   cd backend
   ```

3. Sesuaikan file `.env` jika kredensial MySQL Anda berbeda (default: root tanpa password di localhost):
   ```ini
   PORT=8080
   APP_URL=http://localhost:8080
   DB_USER=root
   DB_PASSWORD=
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=elearning_mobile
   JWT_SECRET=supersecret-flutter-dart-learning-key-2026
   ```

4. Jalankan backend:
   ```bash
   go run main.go
   ```
   *Backend akan otomatis melakukan migrasi tabel (`users`, `materials`, `lessons`) dan memasukkan seeder akun admin serta materi contoh "Pengenalan Dasar Dart".*

---

### 3. Menjalankan Frontend (Svelte)

1. Buka terminal baru dan masuk ke direktori `frontend/`:
   ```bash
   cd frontend
   ```

2. Jalankan development server:
   ```bash
   npm run dev
   ```

3. Buka browser pada URL yang ditampilkan di terminal (biasanya `http://localhost:5173/`).

---

### 4. Build untuk Production (Opsional)
- **Backend:**
  ```bash
  cd backend
  go build -o server.exe .
  ```
- **Frontend:**
  ```bash
  cd frontend
  npm run build
  ```
