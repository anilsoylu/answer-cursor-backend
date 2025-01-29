package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// InitDB initializes the database connection
func InitDB(config *DBConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Run migrations
	migrateDB(config)

	return nil
}

// migrateDB runs database migrations
func migrateDB(config *DBConfig) {
	migrationURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)

	m, err := migrate.New(
		"file://internal/database/migrations",
		migrationURL)
	if err != nil {
		log.Printf("Migration failed to initialize: %v", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration failed: %v", err)
		return
	}
}

// DB returns the database instance
func DB() *gorm.DB {
	return db
} 