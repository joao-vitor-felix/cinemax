package port

import "github.com/joao-vitor-felix/cinemax/internal/core/domain"

type RefreshTokenRepository interface {
	GetByToken(token string) (*domain.RefreshToken, error)
	GenerateToken(userId string) (*domain.RefreshToken, error)
	GenerateAndInvalidateUsedToken(token, userId string) (*domain.RefreshToken, error)
	DeleteToken(token string) error
	DeleteTokensByUserID(userId string) error
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenService interface {
	Execute(input RefreshTokenInput) (*RefreshTokenOutput, error)
}
