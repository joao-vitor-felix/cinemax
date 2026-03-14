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

func (c *RefreshTokenController) Handle(w http.ResponseWriter, r *http.Request) (Response, error) {
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
