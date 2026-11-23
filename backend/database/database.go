package database

import (
	"fmt"
	"log"
	"os"
	"time"

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
	// Retry connection to handle database startup delay in Docker
	for i := 1; i <= 15; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		log.Printf("⏳ Menunggu database MySQL siap (percobaan %d/15)...", i)
		time.Sleep(2 * time.Second)
	}

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

	// 1. Migrasi pembersihan akun lama 'amiya' jika ada
	var oldAmiya models.User
	if err := DB.Where("username = ?", "amiya").First(&oldAmiya).Error; err == nil {
		DB.Delete(&oldAmiya)
		log.Println("🔄 Seeder: Akun lama 'amiya' telah dibersihkan dari database.")
	}

	// 2. Inisialisasi akun admin default 'admin'
	var adminUser models.User
	if err := DB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		adminUser = models.User{
			Username: "admin",
			Password: string(hashedPass),
		}
		DB.Create(&adminUser)
		log.Println("🌱 Seeder: Akun admin default berhasil dibuat (Username: admin | Password: admin123)")
	}
}
