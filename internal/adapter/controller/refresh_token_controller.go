package controller

import (
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type RefreshTokenController struct {
	useCase port.RefreshTokenUseCase
}

func NewRefreshTokenController(useCase port.RefreshTokenUseCase) *RefreshTokenController {
	return &RefreshTokenController{useCase}
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
// @Router /auth/refresh [post]
func (c *RefreshTokenController) Execute(w http.ResponseWriter, r *http.Request) {
	var input port.RefreshTokenInput
	if err := DecodeJSON(r.Body, &input); err != nil {
		slog.Error("failed to decode request body", "error", err)
		WriteError(
			w,
			http.StatusBadRequest,
			domain.InvalidBodyError.Code,
			domain.InvalidBodyError.Message,
		)
		return
	}

	if err := ValidateStruct(input); err != nil {
		slog.Error("validation error", "error", err)
		HandleError(w, err)
		return
	}

	output, err := c.useCase.Execute(r.Context(), input)
	if err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, output)
}
