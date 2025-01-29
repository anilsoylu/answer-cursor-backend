package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint, username, email, role, status string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"role":     role,
		"status":   status,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 saat geçerli
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
} 