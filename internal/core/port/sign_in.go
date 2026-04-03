package port

import "github.com/joao-vitor-felix/cinemax/internal/core/domain"

type SignInRepository interface {
	FindByEmail(email string) (*domain.User, error)
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=12"`
}

type SignInOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignInService interface {
	Execute(input SignInInput) (*SignInOutput, error)
}
