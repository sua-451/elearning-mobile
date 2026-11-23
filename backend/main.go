package main

import (
	"log"
	"os"

	"e-learning-mobile-backend/database"
	"e-learning-mobile-backend/handlers"
	"e-learning-mobile-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present
	_ = godotenv.Load()

	// Connect to MySQL and Run Migrations/Seeders
	database.Connect()

	// Ensure uploads directory exists
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		_ = os.MkdirAll("./uploads", 0755)
	}

	r := gin.Default()

	// CORS Config
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// Static Files for Uploads
	r.Static("/uploads", "./uploads")

	// Public Routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "app": "Dart & Flutter Mobile E-Learning API"})
		})
		api.POST("/auth/login", handlers.Login)
		api.GET("/materials", handlers.GetMaterials)
		api.GET("/materials/:slug", handlers.GetMaterialBySlug)
		api.GET("/materials/:slug/lessons/:lesson_slug", handlers.GetLessonBySlug)
	}

	// Protected Admin Routes
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.GET("/me", handlers.Me)
		admin.POST("/upload", handlers.UploadImage)

		// Materials CRUD
		admin.POST("/materials", handlers.AdminCreateMaterial)
		admin.PUT("/materials/:id", handlers.AdminUpdateMaterial)
		admin.DELETE("/materials/:id", handlers.AdminDeleteMaterial)

		// Lessons CRUD
		admin.POST("/lessons", handlers.AdminCreateLesson)
		admin.PUT("/lessons/:id", handlers.AdminUpdateLesson)
		admin.DELETE("/lessons/:id", handlers.AdminDeleteLesson)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server REST API siap di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
