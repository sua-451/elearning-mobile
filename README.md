# MobileIndex - Platform Pembelajaran Dart & Flutter

Platform katalog dan baca materi pembelajaran mobile development (**Dart & Flutter**) yang sangat ringan, minimalis, dan cepat. Dibangun menggunakan arsitektur modern:
- **Backend**: Golang (Gin + GORM + MySQL/MariaDB + JWT)
- **Frontend**: Svelte (Vite + Tailwind CSS v4)
- **Deployment**: Single-Container Docker (Multi-stage build)

---

## 📁 Struktur Folder Proyek

```text
e-learning-mobile/
├── Dockerfile          # Multi-stage Docker build (Frontend + Backend -> Alpine Runtime)
├── docker-compose.yml  # Orkestrasi container App + MariaDB database
├── .dockerignore       # Pengecualian build context Docker
├── .gitignore          # Pengecualian Git terpusat
├── backend/
│   ├── database/       # Koneksi database, auto-migrate, & seeder otomatis
│   ├── handlers/       # REST API handler (Public & Admin CRUD)
│   ├── middleware/     # JWT authentication middleware
│   ├── models/         # Schema model (User, Material, Lesson)
│   ├── uploads/        # Direktori diagram teknis SVG & upload gambar materi
│   ├── .env.example    # Template konfigurasi environment backend
│   ├── go.mod / go.sum
│   └── main.go         # REST API server & static file server untuk frontend
│
├── frontend/
│   ├── src/
│   │   ├── lib/        # Utilitas API client & theme store (Dark/Light mode)
│   │   ├── components/ # Logo SVG & komponen reusable
│   │   ├── pages/      # Home, MaterialDetail, AdminLogin, AdminDashboard
│   │   ├── App.svelte  # Komponen utama & router hash ringkas
│   │   └── app.css     # Styling Tailwind CSS v4 & custom dark variant
│   ├── package.json
│   └── vite.config.js
│
└── README.md
```

---

## 🚀 Fitur Unggulan

1. **Halaman Beranda (`#/`):**
   - Header ringkas dan langsung menampilkan grid modul materi yang siap dipelajari.
2. **Halaman Baca Materi (`#/materials/[slug]` - Gaya Codelab Split-View):**
   - **Sidebar Kiri (Collapsible):** Daftar bab (*lessons*) ringkas yang dapat dibuka/tutup untuk memaksimalkan ruang baca.
   - **Dynamic Vertical Blocks:** Render konten modular (teks Markdown teknis, diagram teknis SVG, catatan callout, dan blok kode).
   - **Mode Simulasi Dart:** Layout **Split-Screen 2 Kolom** di desktop (kiri: materi bacaan & tombol salin *starter code*; kanan: *sticky* **Embedded DartPad** + tombol buka ke DartPad.dev di tab baru).
   - **Toggle Panel DartPad (Keep-Alive):** Menutup/membuka kembali DartPad tidak mereset atau me-reload iframe kode yang sedang diketik siswa.
3. **Admin Dashboard (`#/admin/dashboard`):**
   - Manajemen CRUD penuh untuk Modul Materi dan Bab Pelajaran.
   - **Dynamic Block Builder (seperti Google Forms):** Tambah, geser urutan (▲/▼), dan hapus blok materi secara visual.
   - Fitur upload file gambar lokal ke server backend atau input URL gambar langsung.
   - Toggle aktivasi DartPad dan input *Starter Code* per bab.
4. **Desain & Tema:**
   - Dukungan penuh **Dark Mode & Light Mode** dengan deteksi preferensi sistem otomatis.
   - Standarisasi ikon monokrom profesional menggunakan **Lucide Svelte** (tanpa emoji).

---

## 🔑 Kredensial Admin Default

Saat database pertama kali diinisialisasi oleh sistem seeder, akun admin default otomatis dibuat:

- **Username:** `admin`
- **Password:** `admin123`

---

## 🐳 Panduan Instalasi & Menjalankan Menggunakan Docker (Rekomendasi)

Dengan menggunakan Docker, Anda tidak perlu menginstall Golang, Node.js, maupun MySQL secara manual di komputer/server Anda. Semua dependensi dan proses build ditangani otomatis di dalam container.

### 1. Prasyarat
- **Docker** dan **Docker Compose** terpasang di komputer/server Anda.
  *(Untuk server Ubuntu/Debian, install via `curl -fsSL https://get.docker.com | sh`)*

---

### 2. Langkah Menjalankan Aplikasi

1. **Clone Repositori:**
   ```bash
   git clone https://github.com/USERNAME/REPO_NAME.git
   cd REPO_NAME
   ```

2. **Jalankan dengan Docker Compose:**
   ```bash
   docker compose up -d --build
   ```

   *Perintah ini akan secara otomatis:*
   - Menyiapkan database MariaDB/MySQL dengan volume data persisten (`db_data`).
   - Mem-build frontend Svelte menjadi bundle statis.
   - Mengompilasi binary Golang untuk Linux Alpine.
   - Menghubungkan aplikasi ke database dan menjalankan migrasi otomatis.
   - Menyajikan frontend UI dan REST API dalam satu container di port `8080`.

3. **Akses Aplikasi di Browser:**
   - **Katalog Pembelajaran:** `http://localhost:8080/` (atau `http://IP_SERVER:8080/`)
   - **Login Dashboard Admin:** `http://localhost:8080/#/admin/login`

---

### 3. Perintah Manajemen Docker yang Berguna

- **Melihat status container:**
  ```bash
  docker compose ps
  ```
- **Melihat log aplikasi secara realtime:**
  ```bash
  docker compose logs -f app
  ```
- **Menghentikan aplikasi:**
  ```bash
  docker compose down
  ```
- **Update aplikasi setelah pull commit terbaru dari Git:**
  ```bash
  git pull origin main
  docker compose up -d --build
  ```

---

## 💻 Panduan Pengembangan Lokal Tanpa Docker (Opsional)

Jika ingin melakukan pengembangan (*local development*) dengan hot-reload:

1. **Prasyarat:**
   - Go (v1.20+)
   - Node.js (v18+) & npm
   - MySQL Server lokal (Laragon / XAMPP)

2. **Jalankan Backend:**
   ```bash
   cd backend
   cp .env.example .env
   go run main.go
   ```

3. **Jalankan Frontend:**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
   *Frontend akan berjalan di `http://localhost:5173` dengan hot-reload dan otomatis terhubung ke backend `http://localhost:8080/api`.*
