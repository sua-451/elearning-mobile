# Panduan File Gambar Dekorasi & Maskot Karakter Anime
===================================================================
Direktori ini (`frontend/public/decorations/`) disiapkan khusus untuk meletakkan gambar dekorasi dan ilustrasi anime yang mempercantik tampilan web MobileIndex.

Format yang didukung: PNG, WebP, JPG, SVG (disarankan PNG transparan atau WebP untuk karakter).

Daftar file yang dapat Anda letakkan di folder ini:
-------------------------------------------------------------------
1. `hero-mascot.png`
   - Fungsi: Karakter maskot anime utama di header banner halaman Home.
   - Posisi: Tampil di sebelah kanan banner hero pada desktop dan di atas/tengah pada mobile.
   - Rekomendasi: Karakter pose menyapa/belajar dengan background transparan (PNG/WebP), ukuran disarankan ~400x500px atau ~600x800px.

2. `bg-wallpaper.png`
   - Fungsi: Wallpaper background di seluruh halaman web (fixed background layer).
   - Posisi: Melapisi layar di balik konten dengan transparansi lembut (~10-15%) dan gradasi pelindung kontras teks.
   - Rekomendasi: Lanskap pemandangan anime, pola abstrak, atau art beresolusi tinggi (misal: 1920x1080px).

3. `menu-corner.png`
   - Fungsi: Karakter mini / chibi / stiker aksen di pojok kanan bawah layar.
   - Posisi: Mengambang di sudut kanan bawah (`fixed bottom-0 right-0`), pointer-events-none sehingga tidak mengganggu klik.
   - Rekomendasi: Karakter chibi duduk/mengintip dengan background transparan, ukuran ~150x150px hingga ~250x250px.

Catatan Keamanan Tampilan:
- Semua slot gambar telah dilengkapi penanganan fallback (error handling & CSS background placeholder) sehingga jika file belum dimasukkan atau gagal dimuat, web akan tetap tampil bersih tanpa ikon gambar rusak (broken image).
