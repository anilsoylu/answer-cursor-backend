package response

import (
	"encoding/json"
	"net/http"

	"github.com/anilsoylu/answer-backend/internal/utils/errors"
)

// Response genel yanıt yapısı
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

// Error sends an error response
func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var errorResponse interface{}

	// Check if it's our custom error type
	if appErr, ok := err.(*errors.AppError); ok {
		statusCode = appErr.StatusCode
		errorResponse = appErr
	} else {
		// Default to internal server error for unknown error types
		statusCode = http.StatusInternalServerError
		errorResponse = map[string]string{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "An unexpected error occurred",
		}
	}

	w.WriteHeader(statusCode)

	resp := Response{
		Success: false,
		Error:   errorResponse,
	}

	json.NewEncoder(w).Encode(resp)
} 