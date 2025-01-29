package repository

import (
	"strings"

	"github.com/anilsoylu/answer-backend/internal/models"
	"github.com/anilsoylu/answer-backend/internal/utils/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		// Check for unique constraint violation
		if strings.Contains(result.Error.Error(), "unique constraint") {
			if strings.Contains(result.Error.Error(), "username") {
				return errors.New(409, "USERNAME_EXISTS", "Username already exists")
			}
			if strings.Contains(result.Error.Error(), "email") {
				return errors.New(409, "EMAIL_EXISTS", "Email already exists")
			}
		}
		return errors.ErrInternalServer
	}
	return nil
}

// FindByIdentifier finds a user by email or username
func (r *UserRepository) FindByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, errors.ErrInternalServer
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return &user, nil
}

// UpdateLastLogin updates the last login date of a user
func (r *UserRepository) UpdateLastLogin(userID uint) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_date", gorm.Expr("NOW()")).Error; err != nil {
		return errors.ErrInternalServer
	}
	return nil
} 