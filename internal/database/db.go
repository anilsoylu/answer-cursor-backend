package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anilsoylu/answer-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConfig represents database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var (
	db *gorm.DB
)

// InitDB sets up the database connection
func InitDB(config *DBConfig) error {
	var err error

	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Istanbul",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	// Custom logger configuration
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open connection
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Connected Successfully to Database")

	// Create enum types
	if err := createEnumTypes(); err != nil {
		return fmt.Errorf("failed to create enum types: %v", err)
	}

	// Auto Migrate
	if err = db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database Migrated Successfully")
	return nil
}

// DB returns the database instance
func DB() *gorm.DB {
	if db == nil {
		log.Fatal("Database connection not initialized")
	}
	return db
}

// createEnumTypes creates the required enum types in PostgreSQL
func createEnumTypes() error {
	// Drop existing enum types if they exist
	db.Exec(`DROP TYPE IF EXISTS user_status CASCADE;`)
	db.Exec(`DROP TYPE IF EXISTS user_role CASCADE;`)

	// Create user_status enum
	if err := db.Exec(`CREATE TYPE user_status AS ENUM ('active', 'passive', 'banned', 'frozen');`).Error; err != nil {
		return err
	}

	// Create user_role enum
	if err := db.Exec(`CREATE TYPE user_role AS ENUM ('USER', 'EDITOR', 'ADMIN', 'SUPER_ADMIN');`).Error; err != nil {
		return err
	}

	return nil
} 