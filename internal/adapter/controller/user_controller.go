package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	_ "github.com/joao-vitor-felix/cinemax/internal/adapter/http/middleware"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type UserController struct {
	service port.UserService
}

func NewUserController(service port.UserService) *UserController {
	return &UserController{service}
}

// TODO: put it inside the struct
var validate = validator.New(validator.WithRequiredStructEnabled())

// Register godoc
// @Summary Register a user
// @Description Register a new user with the provided information.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body port.RegisterUserInput true "User registration data"
// @Success 201 {object} Resource "User registered successfully"
// @Failure 400 {object} middleware.ErrorResponse "Bad request"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /auth/sign-up [post]
func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) (map[string]any, error) {
	var body port.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("failed to decode request body", "error", err)
		return nil, domain.InvalidBodyError
	}

	if err := validate.Struct(body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			error := ve[0]
			errorMsg := BuildValidationErrorMessage(error.Field(), error.Tag())
			return nil, domain.ValidationError(errorMsg)
		}
		return nil, domain.InternalServerError
	}

	_, err := uc.service.Register(body)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"data": NewResource(
			nil,
			map[string]Link{
				"sign-in": {
					Href:   "/auth/sign-in",
					Method: "POST",
				},
			}),
		"status": http.StatusCreated,
	}, nil
}
