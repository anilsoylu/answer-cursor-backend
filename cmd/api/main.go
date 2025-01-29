package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anilsoylu/answer-backend/internal/database"
	"github.com/anilsoylu/answer-backend/internal/handlers"
	"github.com/anilsoylu/answer-backend/internal/services"
	"github.com/anilsoylu/answer-backend/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	dbConfig := &database.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if err := database.InitDB(dbConfig); err != nil {
		log.Fatal("Could not initialize database: ", err)
	}

	// Initialize services
	authService := services.NewAuthService(database.DB())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORS())

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.PUT("/password", middleware.AuthMiddleware(), authHandler.UpdatePassword)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.PUT("/role", authHandler.UpdateUserRole)
				users.PATCH("/status", authHandler.UpdateUserStatus)
				users.PUT("/profile", authHandler.UpdateProfile)
				users.POST("/ban", authHandler.BanUser)
				users.POST("/freeze", authHandler.FreezeAccount)
				users.DELETE("/:id", authHandler.DeleteAccount)
			}
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
} 