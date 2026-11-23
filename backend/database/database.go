package database

import (
	"fmt"
	"log"
	"os"

	"e-learning-mobile-backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "elearning_mobile"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("⚠️ Gagal koneksi ke database MySQL (%s). Pastikan MySQL berjalan dan database '%s' sudah dibuat.\nError: %v", dsn, dbName, err)
		return
	}

	log.Println("✅ Terhubung ke database MySQL!")

	// Auto-Migrate
	err = DB.AutoMigrate(&models.User{}, &models.Material{}, &models.Lesson{})
	if err != nil {
		log.Printf("Gagal auto-migrate: %v", err)
		return
	}

	Seed()
}

func Seed() {
	if DB == nil {
		return
	}

	// 1. Seed Admin User
	var adminCount int64
	DB.Model(&models.User{}).Count(&adminCount)
	if adminCount == 0 {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Username: "admin",
			Password: string(hashedPass),
		}
		DB.Create(&admin)
		log.Println("🌱 Seeder: Akun admin default berhasil dibuat (Username: admin | Password: admin123)")
	}

	// 2. Seed Default Material & Lessons
	var materialCount int64
	DB.Model(&models.Material{}).Count(&materialCount)
	if materialCount == 0 {
		material := models.Material{
			Title:       "Pengenalan Dasar Dart",
			Slug:        "pengenalan-dasar-dart",
			Description: "Pelajari sintaks dasar bahasa pemrograman Dart yang menjadi fondasi utama framework Flutter.",
			OrderIndex:  1,
			Lessons: []models.Lesson{
				{
					Title:         "1. Pengantar Bahasa Dart & Struktur main()",
					Slug:          "pengantar-dart-dan-main",
					OrderIndex:    1,
					ImageURL:      "https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&w=1200&q=80",
					Content:       "### Selamat Datang di Dunia Dart!\n\nDart adalah bahasa pemrograman modern berorientasi objek yang dioptimalkan untuk pembuatan aplikasi cepat di berbagai platform (Mobile, Web, Desktop) oleh Google.\n\nSetiap program Dart dieksekusi pertama kali dari fungsi **`main()`**.\n\nContoh struktur dasar:\n```dart\nvoid main() {\n  print('Halo Pengembang Flutter!');\n}\n```",
					Notes:         "💡 Tips: Fungsi main() adalah titik awal (entry point) mutlak dalam program Dart dan Flutter.",
					StarterCode:   "void main() {\n  var greeting = 'Halo dari Dart!';\n  print(greeting);\n}",
					EnableDartpad: true,
				},
				{
					Title:         "2. Variabel & Tipe Data Dasar",
					Slug:          "variabel-dan-tipe-data",
					OrderIndex:    2,
					ImageURL:      "",
					Content:       "### Tipe Data Bawaan di Dart\n\nDart mendukung *type inference* menggunakan kata kunci `var`, atau deklarasi tipe data eksplisit:\n- `int` untuk bilangan bulat\n- `double` untuk bilangan desimal\n- `String` untuk teks\n- `bool` untuk boolean (`true` / `false`)\n\nSelain itu, Dart memiliki fitur **Null Safety** bawaan yang mencegah kesalahan error null reference saat runtime.",
					Notes:         "📌 Catatan: Gunakan tipe eksplisit seperti `int` atau `String` jika ingin kode lebih ekspresif dan mudah dibaca oleh tim.",
					StarterCode:   "void main() {\n  int umur = 21;\n  double tinggi = 175.5;\n  String nama = 'Bintang';\n  bool isStudent = true;\n\n  print('Nama: $nama, Umur: $umur tahun, Status: $isStudent');\n}",
					EnableDartpad: true,
				},
				{
					Title:         "3. Fungsi & Parameter",
					Slug:          "fungsi-dan-parameter",
					OrderIndex:    3,
					ImageURL:      "",
					Content:       "### Membuat Fungsi di Dart\n\nFungsi di Dart mendukung parameter wajib (*positional*) dan parameter bernama (*named parameters* yang dibungkus dengan kurung kurawal `{}`). Parameter bernama sangat umum digunakan dalam pembuatan Widget Flutter.\n\n```dart\nvoid sapa({required String nama, String salam = 'Halo'}) {\n  print('$salam, $nama!');\n}\n```",
					Notes:         "✨ Parameter bernama (named parameters) memudahkan pembacaan fungsi ketika memiliki banyak argumen.",
					StarterCode:   "void sapa({required String nama, String role = 'Mobile Developer'}) {\n  print('Halo $nama! Selamat belajar sebagai $role.');\n}\n\nvoid main() {\n  sapa(nama: 'Alex');\n  sapa(nama: 'Citra', role: 'Flutter Architect');\n}",
					EnableDartpad: true,
				},
			},
		}
		DB.Create(&material)
		log.Println("🌱 Seeder: Contoh materi 'Pengenalan Dasar Dart' berhasil ditambahkan!")
	}
}
