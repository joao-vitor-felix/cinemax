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
		Data:   output,
		Status: http.StatusOK,
	}, nil
}
