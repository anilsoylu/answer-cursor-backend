package password

import (
	"github.com/anilsoylu/answer-backend/internal/utils/errors"
	"golang.org/x/crypto/bcrypt"
)

// Hash creates a bcrypt hash of the password
func Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.ErrInternalServer
	}
	return string(hashedPassword), nil
}

// Verify checks if the provided password matches the hash
func Verify(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.ErrInvalidCredentials
		}
		return errors.ErrInternalServer
	}
	return nil
} 