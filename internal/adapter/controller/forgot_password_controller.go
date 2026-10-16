package controller

import (
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type ForgotPasswordController struct {
	useCase port.ForgotPasswordUseCase
}

func NewForgotPasswordController(useCase port.ForgotPasswordUseCase) *ForgotPasswordController {
	return &ForgotPasswordController{useCase}
}

// Execute godoc
//
//	@Summary		Forgot password
//	@Description	Send a forgot password email
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		port.ForgotPasswordInput	true	"Forgot Password Input"
//	@Success		200		{object}	MessageResponse				"Message response"
//	@Failure		400		{object}	ErrorResponse				"Bad request"
//	@Failure		500		{object}	ErrorResponse				"Internal server error"
//	@Router			/auth/forgot-password [post]
func (c *ForgotPasswordController) Execute(w http.ResponseWriter, r *http.Request) {
	var body port.ForgotPasswordInput

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

	WriteJSON(w, 200, MessageResponse{
		Message: "if you have an e-mail registered, a code was sent to you.",
	})
}
