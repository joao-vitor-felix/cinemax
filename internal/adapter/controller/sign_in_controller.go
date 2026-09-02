package controller

import (
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type SignInController struct {
	service port.SignInService
}

func NewSignInController(service port.SignInService) *SignInController {
	return &SignInController{service}
}

// Execute godoc
//
//	@Summary		Sign in a user
//	@Description	Authenticate a user and return access and refresh tokens.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		port.SignInInput	true	"User credentials"
//	@Success		200			{object}	Resource[port.SignInOutput]	"User authenticated successfully"
//	@Failure		400			{object}	ErrorResponse		"Bad request (invalid body or validation error)"
//	@Failure		401			{object}	ErrorResponse		"Unauthorized (invalid credentials)"
//	@Failure		500			{object}	ErrorResponse		"Internal server error"
//	@Router			/auth/sign-in [post]
func (c *SignInController) Execute(w http.ResponseWriter, r *http.Request) {
	var body port.SignInInput
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
		e := domain.ValidationError(err.Error())
		WriteError(w, http.StatusBadRequest, e.Code, e.Message)
		return
	}
	output, err := c.service.Execute(r.Context(), body)

	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, NewResource(
		output,
		map[string]Link{
			"refresh-token": {
				Href:   "/auth/refresh-token",
				Method: "POST",
			},
		},
	))
}
