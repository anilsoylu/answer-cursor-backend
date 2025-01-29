package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, e.Field()+" is required")
			case "email":
				errorMessages = append(errorMessages, e.Field()+" must be a valid email")
			case "min":
				errorMessages = append(errorMessages, e.Field()+" must be at least "+e.Param()+" characters")
			case "max":
				errorMessages = append(errorMessages, e.Field()+" must be at most "+e.Param()+" characters")
			case "oneof":
				errorMessages = append(errorMessages, e.Field()+" must be one of: "+e.Param())
			default:
				errorMessages = append(errorMessages, e.Field()+" is invalid")
			}
		}
		return strings.Join(errorMessages, ", ")
	}
	return "Invalid input"
} 