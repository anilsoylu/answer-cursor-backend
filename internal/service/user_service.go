package service

import (
	"github.com/anilsoylu/answer-backend/internal/models"
	"github.com/anilsoylu/answer-backend/internal/repository"
	"github.com/anilsoylu/answer-backend/internal/utils/errors"
	"github.com/anilsoylu/answer-backend/internal/utils/password"
	"github.com/anilsoylu/answer-backend/internal/utils/token"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register handles user registration
func (s *UserService) Register(req *models.RegisterRequest) error {
	// Create user instance
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Status:   models.StatusActive,
		Role:     models.RoleUser,
	}

	// Hash password
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Save user
	return s.repo.Create(user)
}

// Login handles user login
func (s *UserService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by identifier (email or username)
	user, err := s.repo.FindByIdentifier(req.Identifier)
	if err != nil {
		return nil, err
	}

	// Check password
	if err := password.Verify(user.Password, req.Password); err != nil {
		return nil, err
	}

	// Check if user is active
	if user.Status != models.StatusActive {
		return nil, errors.New(403, "ACCOUNT_INACTIVE", "Your account is not active")
	}

	// Generate token
	tokenString, expiresIn, err := token.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Update last login
	if err := s.repo.UpdateLastLogin(user.ID); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:     tokenString,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
	}, nil
} 