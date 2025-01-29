package seed

import (
	"log"
	"os"
	"time"

	"github.com/anilsoylu/answer-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateSuperAdmin creates a super admin user if it doesn't exist
func CreateSuperAdmin(db *gorm.DB) error {
	// Get admin credentials from environment variables
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	email := os.Getenv("ADMIN_EMAIL")
	role := os.Getenv("ADMIN_ROLE")
	status := os.Getenv("ADMIN_STATUS")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}

	// Create super admin user
	admin := models.User{
		Username:      username,
		Email:        email,
		Password:     string(hashedPassword),
		Avatar:       "/uploads/default/avatar.png",
		Status:       models.UserStatus(status),
		Role:         models.UserRole(role),
		IsRootAdmin:  true,
		CreatedAt:    time.Now(),
		LastLoginDate: time.Now(),
	}

	// Check if super admin already exists
	var existingAdmin models.User
	result := db.Where("email = ? OR username = ?", email, username).First(&existingAdmin)
	if result.Error == nil {
		log.Printf("Super admin already exists")
		return nil
	}

	// Create super admin if not exists
	result = db.Create(&admin)
	if result.Error != nil {
		log.Printf("Failed to create super admin user: %v", result.Error)
		return result.Error
	}

	log.Printf("Super admin user created successfully")
	return nil
} 