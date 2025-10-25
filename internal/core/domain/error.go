package domain

import (
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func ValidationError(error string) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    error,
		StatusCode: http.StatusBadRequest,
	}
}

var (
	InvalidBodyError = NewAppError(
		"INVALID_REQUEST_BODY",
		"the request body is invalid",
		http.StatusBadRequest,
	)

	InternalServerError = NewAppError(
		"INTERNAL_SERVER_ERROR",
		"an internal server error occurred",
		http.StatusInternalServerError,
	)

	// Thrown when email or phone is already in use
	ContactInfoUnavailableError = NewAppError(
		"CONTACT_DATA_UNAVAILABLE",
		"email or phone already in use",
		http.StatusConflict,
	)

	// Thrown when gender is invalid
	InvalidGenderError = NewAppError(
		"INVALID_GENDER",
		"invalid gender",
		http.StatusBadRequest,
	)

	// Thrown when user is under 13 years old
	UserTooYoungError = NewAppError(
		"TOO_YOUNG",
		"user must be at least 13 years old",
		http.StatusBadRequest,
	)
)
