package controller

import (
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type ResetPasswordController struct {
	useCase port.ResetPasswordUseCase
}

func NewResetPasswordController(useCase port.ResetPasswordUseCase) *ResetPasswordController {
	return &ResetPasswordController{useCase}
}

// Execute godoc
//
//	@Summary		Reset password
//	@Description	Resets the user's password using a valid token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		port.ResetPasswordInput	true	"Reset Password Input"
//	@Success		200		{object}	MessageResponse				"Message response"
//	@Failure		400		{object}	ErrorResponse				"Bad request"
//	@Failure		404		{object}	ErrorResponse				"Not Found"
//	@Failure		500		{object}	ErrorResponse				"Internal server error"
//	@Router			/auth/reset-password [post]
func (c *ResetPasswordController) Execute(w http.ResponseWriter, r *http.Request) {
	var body port.ResetPasswordInput

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

	if err := c.useCase.Execute(r.Context(), body); err != nil {
		HandleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, MessageResponse{
		Message: "Password has been successfully reset.",
	})
}
