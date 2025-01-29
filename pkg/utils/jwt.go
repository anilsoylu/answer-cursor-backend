package utils

import (
	"os"
	"time"

	"github.com/anilsoylu/answer-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint, username, email string, role models.UserRole, status models.UserStatus) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"role":     role,
		"status":   status,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
} 