package errors

import "net/http"

// AppError represents the application-wide error structure
type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

// New creates a new AppError
func New(statusCode int, code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// Common application errors
var (
	ErrInvalidInput = &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "INVALID_INPUT",
		Message:    "Invalid input parameters",
	}

	ErrUnauthorized = &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    "Authentication required",
	}

	ErrForbidden = &AppError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    "Access forbidden",
	}

	ErrNotFound = &AppError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
	}

	ErrInternalServer = &AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "Internal server error",
	}

	ErrUserExists = &AppError{
		StatusCode: http.StatusConflict,
		Code:       "USER_EXISTS",
		Message:    "User already exists",
	}

	ErrInvalidCredentials = &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
	}
) 