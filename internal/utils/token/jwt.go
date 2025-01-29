package token

import (
	"fmt"
	"time"

	"github.com/anilsoylu/answer-backend/internal/utils/errors"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

// Initialize sets the JWT secret key
func Initialize(secret string) {
	jwtKey = []byte(secret)
}

// Claims represents the JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token
func GenerateToken(userID uint, role string) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, errors.ErrInternalServer
	}

	return tokenString, expirationTime.Unix(), nil
}

// ValidateToken validates the JWT token
func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.ErrUnauthorized
		}
		return nil, errors.ErrUnauthorized
	}

	if !token.Valid {
		return nil, errors.ErrUnauthorized
	}

	return claims, nil
} 