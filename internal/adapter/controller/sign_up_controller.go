package controller

import (
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type SignUpController struct {
	service port.SignUpService
}

func NewSignUpController(service port.SignUpService) *SignUpController {
	return &SignUpController{service}
}

// Execute godoc
//
//	@Summary		Sign up a user
//	@Description	Register a new user with the provided information.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		port.SignUpInput	true	"User registration data"
//	@Success		201		"User registered successfully"
//	@Failure		400		{object}	ErrorResponse		"Bad request"
//	@Failure		500		{object}	ErrorResponse		"Internal server error"
//	@Router			/auth/sign-up [post]
func (c *SignUpController) Execute(w http.ResponseWriter, r *http.Request) {
	var body port.SignUpInput
	if err := DecodeJSON(r.Body, &body); err != nil {
		slog.Error("failed to decode request body", "error", err)
		WriteError(
			w,
			http.StatusBadRequest,
			domain.InvalidBodyError.Code,
			domain.InvalidBodyError.Message,
		)
		return
	}

	if err := ValidateStruct(body); err != nil {
		slog.Error("validation error", "error", err)
		HandleError(w, err)
		return
	}

	_, err := c.service.Execute(r.Context(), body)
	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, nil)
}
