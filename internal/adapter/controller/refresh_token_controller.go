package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type RefreshTokenController struct {
	service port.RefreshTokenService
}

func NewRefreshTokenController(service port.RefreshTokenService) *RefreshTokenController {
	return &RefreshTokenController{service}
}

// Execute godoc
//
// @Summary Refresh access token
// @Description Refresh the access token using a valid refresh token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshToken body port.RefreshTokenInput true "Refresh token"
// @Success 200 {object} port.RefreshTokenOutput "Access token refreshed successfully"
// @Failure 400 {object} ErrorResponse "Bad request (invalid body or validation error)"
// @Failure 404 {object} ErrorResponse "Not found (token or user not found)"
// @Failure 401 {object} ErrorResponse "Unauthorized (invalid refresh token)"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/refresh-token [post]
func (c *RefreshTokenController) Execute(w http.ResponseWriter, r *http.Request) (Response, error) {
	var input port.RefreshTokenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("failed to decode request body", "error", err)
		return Response{}, domain.InvalidBodyError
	}

	if err := ValidateStruct(input); err != nil {
		slog.Error("validation error", "error", err)
		return Response{}, err
	}

	output, err := c.service.Execute(input)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Data:   output,
		Status: http.StatusOK,
	}, nil
}
