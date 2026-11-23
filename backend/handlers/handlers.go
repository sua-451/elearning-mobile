package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"e-learning-mobile-backend/database"
	"e-learning-mobile-backend/middleware"
	"e-learning-mobile-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- AUTH HANDLERS ---

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password harus diisi"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan token autentikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func Me(c *gin.Context) {
	username, _ := c.Get("username")
	userID, _ := c.Get("user_id")
	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
	})
}

// --- PUBLIC MATERIAL & LESSON HANDLERS ---

func GetMaterials(c *gin.Context) {
	var materials []models.Material
	if err := database.DB.Order("order_index asc, id asc").Preload("Lessons", func(db *gorm.DB) *gorm.DB {
		return db.Order("lessons.order_index asc, lessons.id asc")
	}).Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar materi"})
		return
	}

	c.JSON(http.StatusOK, materials)
}

func GetMaterialBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var material models.Material
	if err := database.DB.Preload("Lessons", func(db *gorm.DB) *gorm.DB {
		return db.Order("lessons.order_index asc, lessons.id asc")
	}).Where("slug = ?", slug).First(&material).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, material)
}

func GetLessonBySlug(c *gin.Context) {
	materialSlug := c.Param("slug")
	lessonSlug := c.Param("lesson_slug")

	var material models.Material
	if err := database.DB.Where("slug = ?", materialSlug).First(&material).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi induk tidak ditemukan"})
		return
	}

	var lesson models.Lesson
	if err := database.DB.Where("material_id = ? AND slug = ?", material.ID, lessonSlug).First(&lesson).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bab tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

// --- ADMIN CRUD: MATERIALS ---

type MaterialInput struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order_index"`
}

func AdminCreateMaterial(c *gin.Context) {
	var input MaterialInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := strings.TrimSpace(strings.ToLower(input.Slug))
	material := models.Material{
		Title:       input.Title,
		Slug:        slug,
		Description: input.Description,
		OrderIndex:  input.OrderIndex,
	}

	if err := database.DB.Create(&material).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug sudah digunakan atau data tidak valid"})
		return
	}

	c.JSON(http.StatusCreated, material)
}

func AdminUpdateMaterial(c *gin.Context) {
	id := c.Param("id")
	var material models.Material
	if err := database.DB.First(&material, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi tidak ditemukan"})
		return
	}

	var input MaterialInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	material.Title = input.Title
	material.Slug = strings.TrimSpace(strings.ToLower(input.Slug))
	material.Description = input.Description
	material.OrderIndex = input.OrderIndex

	if err := database.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug sudah digunakan atau data tidak valid"})
		return
	}

	c.JSON(http.StatusOK, material)
}

func AdminDeleteMaterial(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Material{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus materi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Materi berhasil dihapus"})
}

// --- ADMIN CRUD: LESSONS ---

type LessonInput struct {
	MaterialID    uint   `json:"material_id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Slug          string `json:"slug" binding:"required"`
	OrderIndex    int    `json:"order_index"`
	ImageURL      string `json:"image_url"`
	Content       string `json:"content"`
	Notes         string `json:"notes"`
	StarterCode   string `json:"starter_code"`
	Blocks        string `json:"blocks"`
	EnableDartpad bool   `json:"enable_dartpad"`
}

func AdminCreateLesson(c *gin.Context) {
	var input LessonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson := models.Lesson{
		MaterialID:    input.MaterialID,
		Title:         input.Title,
		Slug:          strings.TrimSpace(strings.ToLower(input.Slug)),
		OrderIndex:    input.OrderIndex,
		ImageURL:      input.ImageURL,
		Content:       input.Content,
		Notes:         input.Notes,
		StarterCode:   input.StarterCode,
		Blocks:        input.Blocks,
		EnableDartpad: input.EnableDartpad,
	}

	if err := database.DB.Create(&lesson).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat bab baru"})
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

func AdminUpdateLesson(c *gin.Context) {
	id := c.Param("id")
	var lesson models.Lesson
	if err := database.DB.First(&lesson, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bab tidak ditemukan"})
		return
	}

	var input LessonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson.MaterialID = input.MaterialID
	lesson.Title = input.Title
	lesson.Slug = strings.TrimSpace(strings.ToLower(input.Slug))
	lesson.OrderIndex = input.OrderIndex
	lesson.ImageURL = input.ImageURL
	lesson.Content = input.Content
	lesson.Notes = input.Notes
	lesson.StarterCode = input.StarterCode
	lesson.Blocks = input.Blocks
	lesson.EnableDartpad = input.EnableDartpad

	if err := database.DB.Save(&lesson).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui bab"})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func AdminDeleteLesson(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Lesson{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus bab"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bab berhasil dihapus"})
}

// --- FILE UPLOAD HANDLER ---

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File gambar tidak ditemukan dalam request"})
		return
	}

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strconv.FormatInt(time.Now().Unix(), 10), ext)
	targetPath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file upload"})
		return
	}

	baseURL := os.Getenv("APP_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	fileURL := fmt.Sprintf("%s/uploads/%s", strings.TrimRight(baseURL, "/"), filename)

	c.JSON(http.StatusOK, gin.H{
		"url":      fileURL,
		"filename": filename,
	})
}

// --- ADMIN USER MANAGEMENT HANDLERS ---

type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// AdminGetUsers lists all admin users (without exposing hashed passwords)
func AdminGetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Select("id, username, created_at, updated_at").Order("id asc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar admin"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// AdminCreateUser creates a new admin user
func AdminCreateUser(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username minimal 3 karakter"})
		return
	}
	if len(req.Password) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 5 karakter"})
		return
	}

	// Check if username already exists
	var existing models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}

	newUser := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan admin baru"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Admin berhasil ditambahkan",
		"user": gin.H{
			"id":       newUser.ID,
			"username": newUser.Username,
		},
	})
}

// AdminDeleteUser deletes an admin user by ID
func AdminDeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID admin tidak valid"})
		return
	}

	currentUserIDVal, exists := c.Get("user_id")
	if exists {
		var currentID uint64
		switch v := currentUserIDVal.(type) {
		case uint:
			currentID = uint64(v)
		case float64:
			currentID = uint64(v)
		case int:
			currentID = uint64(v)
		case uint64:
			currentID = v
		}
		if currentID == targetID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak dapat menghapus akun admin yang sedang aktif login"})
			return
		}
	}

	// Ensure there is at least one admin remaining
	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)
	if totalUsers <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak dapat menghapus satu-satunya akun admin"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin berhasil dihapus"})
}

// AdminChangePassword updates the password of the currently logged-in admin
func AdminChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama dan password baru wajib diisi"})
		return
	}

	if len(req.NewPassword) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password baru minimal 5 karakter"})
		return
	}

	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi autentikasi tidak ditemukan"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, currentUserIDVal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama tidak sesuai"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses kata sandi baru"})
		return
	}

	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diperbarui"})
}

