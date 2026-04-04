package controller

import (
	"encoding/json"
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
func (c *SignInController) Execute(w http.ResponseWriter, r *http.Request) (Response, error) {
	var body port.SignInInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("failed to decode request body", "error", err)
		return Response{}, domain.InvalidBodyError
	}

	if err := ValidateStruct(body); err != nil {
		slog.Error("validation error", "error", err)
		return Response{}, err
	}
	output, err := c.service.Execute(body)

	if err != nil {
		return Response{}, err
	}

	return Response{
		Data: NewResource(
			output,
			map[string]Link{
				"refresh-token": {
					Href:   "/auth/refresh-token",
					Method: "POST",
				},
			},
		),
		Status: http.StatusOK,
	}, nil
}
