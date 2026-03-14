package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type UserController struct {
	service port.UserService
}

func NewUserController(service port.UserService) *UserController {
	return &UserController{service}
}

// Register godoc
// @Summary Register a user
// @Description Register a new user with the provided information.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body port.RegisterUserInput true "User registration data"
// @Success 201 {object} Resource "User registered successfully"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/sign-up [post]
func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) (Response, error) {
	var body port.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("failed to decode request body", "error", err)
		return Response{}, domain.InvalidBodyError
	}

	if err := ValidateStruct(body); err != nil {
		slog.Error("validation error", "error", err)
		return Response{}, err
	}

	_, err := uc.service.Register(body)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Data: NewResource(
			nil,
			map[string]Link{
				"sign-in": {
					Href:   "/auth/sign-in",
					Method: "POST",
				},
			}),
		Status: http.StatusCreated,
	}, nil
}
