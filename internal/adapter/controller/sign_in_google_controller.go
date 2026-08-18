package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type SignInGoogleController struct {
	service port.SignInGoogleService
}

func NewSignInGoogleController(service port.SignInGoogleService) *SignInGoogleController {
	return &SignInGoogleController{service}
}

// Execute godoc
//
//	@Summary		Sign in a user with Google OAuth
//	@Description	Authenticate a user using a Google OAuth code and return access and refresh tokens.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			code	body		port.SignInGoogleInput	true	"Google OAuth authorization code"
//	@Success		200		{object}	Resource[port.SignInOutput]	"User authenticated successfully"
//	@Failure		400		{object}	ErrorResponse		"Bad request (invalid body or validation error)"
//	@Failure		401		{object}	ErrorResponse		"Unauthorized (invalid or expired code, email not verified)"
//	@Failure		500		{object}	ErrorResponse		"Internal server error"
//	@Router			/auth/sign-in/google [post]
func (c *SignInGoogleController) Execute(w http.ResponseWriter, r *http.Request) (Response, error) {
	var body port.SignInGoogleInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("failed to decode request body", "error", err)
		return Response{}, domain.InvalidBodyError
	}

	if err := ValidateStruct(body); err != nil {
		slog.Error("validation error", "error", err)
		return Response{}, err
	}

	output, err := c.service.Execute(r.Context(), body)
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
