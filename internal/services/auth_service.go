package services

import (
	"errors"
	"time"

	"github.com/anilsoylu/answer-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotActive       = errors.New("user is not active")
	ErrUserNotFound        = errors.New("user not found")
	ErrDuplicateEntry     = errors.New("duplicate entry")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(user *models.User) error {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err == nil {
		return ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := s.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(identifier, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if user.Status != models.StatusActive {
		return nil, ErrUserNotActive
	}

	// Update last login date
	user.LastLoginDate = time.Now()
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) UpdateUserRole(userID uint, newRole models.UserRole, requester *models.User) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check permissions and role hierarchy
	if user.Role == models.RoleSuperAdmin && !requester.IsRootAdmin {
		return errors.New("only root admin can change SUPER_ADMIN's role")
	}

	if newRole == models.RoleSuperAdmin && !requester.IsRootAdmin {
		return errors.New("only root admin can assign SUPER_ADMIN role")
	}

	if requester.Role == models.RoleAdmin {
		if user.Role != models.RoleUser && user.Role != models.RoleEditor {
			return errors.New("admin can only modify USER and EDITOR roles")
		}
	}

	user.Role = newRole
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateUserStatus(userID uint, newStatus models.UserStatus, requesterRole models.UserRole) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check permissions
	if user.Role == models.RoleSuperAdmin {
		return errors.New("SUPER_ADMIN's status cannot be changed")
	}

	if requesterRole == models.RoleAdmin && user.Role == models.RoleAdmin {
		return errors.New("admins cannot change other admin's status")
	}

	user.Status = newStatus
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateProfile(userID uint, username, email, password, avatar string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Check for duplicate username/email
	if username != "" && username != user.Username {
		var existingUser models.User
		if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
			return nil, ErrDuplicateEntry
		}
		user.Username = username
	}

	if email != "" && email != user.Email {
		var existingUser models.User
		if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
			return nil, ErrDuplicateEntry
		}
		user.Email = email
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash and save new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetUserByID(userID uint, user *models.User) error {
	if err := s.db.First(user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
} 