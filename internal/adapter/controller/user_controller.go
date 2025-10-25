package controller

import (
	"encoding/json"
	"fmt"
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

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) (any, int, error) {
	var body port.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, 0, domain.InvalidBodyError
	}

	if err := validate.Struct(body); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		//FIXME: find a better way to display validation errors
		return nil, 0, domain.ValidationError(validationErrors.Error())
	}

	user, err := uc.service.Register(body)
	if err != nil {
		return nil, 0, err
	}

	return NewResource(user, map[string]Link{
		"self": {
			Href:   fmt.Sprintf("/users/%s", user.ID),
			Method: "GET",
		},
	}), http.StatusCreated, nil
}
