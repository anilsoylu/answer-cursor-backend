package validator

import (
	"github.com/anilsoylu/answer-backend/internal/utils/errors"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct based on validate tags
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.ErrInternalServer
		}

		var errorMessage string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errorMessage = err.Field() + " is required"
			case "email":
				errorMessage = err.Field() + " must be a valid email"
			case "min":
				errorMessage = err.Field() + " must be at least " + err.Param() + " characters long"
			case "max":
				errorMessage = err.Field() + " must be at most " + err.Param() + " characters long"
			default:
				errorMessage = "Validation failed on " + err.Field()
			}
			break // Only return the first error
		}

		return errors.New(400, "VALIDATION_ERROR", errorMessage)
	}
	return nil
} 