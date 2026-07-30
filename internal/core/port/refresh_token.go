package port

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type RefreshTokenRepository interface {
	GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	GenerateToken(ctx context.Context, userId string) (*domain.RefreshToken, error)
	GenerateAndInvalidateUsedToken(ctx context.Context, token, userId string) (*domain.RefreshToken, error)
	DeleteToken(ctx context.Context, token string) error
	DeleteTokensByUserID(ctx context.Context, userId string) error
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenService interface {
	Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error)
}
