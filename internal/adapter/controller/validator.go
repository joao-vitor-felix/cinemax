package controller

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

func ValidateStruct(input any) error {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			validationError := ve[0]
			errorMsg := BuildValidationErrorMessage(validationError.Field(), validationError.Tag())
			return domain.ValidationError(errorMsg)
		}
		return domain.InternalServerError
	}
	return nil
}

func BuildValidationErrorMessage(field, tag string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s does not meet the minimum length requirement", field)
	case "max":
		return fmt.Sprintf("%s exceeds the maximum length requirement", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "e164":
		return fmt.Sprintf("%s must be a valid E.164 formatted phone number", field)
	case "datetime":
		return fmt.Sprintf("%s must be a valid date in the format YYYY-MM-DD", field)
	case "oneof":
		return fmt.Sprintf("%s contains an invalid value", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
