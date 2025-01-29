package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anilsoylu/answer-backend/internal/database"
	"github.com/anilsoylu/answer-backend/internal/database/seed"
	"github.com/anilsoylu/answer-backend/internal/handlers"
	"github.com/anilsoylu/answer-backend/internal/routes"
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

	// Create super admin user
	if err := seed.CreateSuperAdmin(database.DB()); err != nil {
		log.Fatal("Failed to create super admin user:", err)
	}

	// Initialize services
	authService := services.NewAuthService(database.DB())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORS())

	// Setup routes
	routes.SetupAuthRoutes(router, authHandler)
	routes.SetupAdminRoutes(router, authHandler)

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