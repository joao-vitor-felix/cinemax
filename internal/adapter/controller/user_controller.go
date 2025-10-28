package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type UserController struct {
	service port.UserService
}

func NewUserController(service port.UserService) *UserController {
	return &UserController{service}
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) (map[string]any, error) {
	var body port.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, domain.InvalidBodyError
	}

	if err := validate.Struct(body); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		//FIXME: find a better way to display validation errors
		return nil, domain.ValidationError(validationErrors.Error())
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
