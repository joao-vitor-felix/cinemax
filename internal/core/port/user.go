package port

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type UserRepository interface {
	Create(user *domain.User) (*domain.User, error)
	IsContactInfoAvailable(email, phone string) (bool, error)
}

type RegisterUserInput struct {
	FirstName       string        `json:"first_name" validate:"required,min=2"`
	LastName        string        `json:"last_name" validate:"required,min=2"`
	Email           string        `json:"email" validate:"required,email"`
	Phone           string        `json:"phone" validate:"required,e164"`
	Password        string        `json:"password" validate:"required,min=8,max=12"`
	DateOfBirth     string        `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Gender          domain.Gender `json:"gender" validate:"oneof=male female prefer_not_to_say other"`
	ProfilePhotoURL *string       `json:"profile_photo_url,omitempty" validate:"omitempty,url"`
}

type UserService interface {
	Register(input RegisterUserInput) (*domain.User, error)
}
