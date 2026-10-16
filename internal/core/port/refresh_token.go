package port

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type RefreshTokenRepository interface {
	GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	GenerateToken(ctx context.Context, userID string) (*domain.RefreshToken, error)
	GenerateAndInvalidateUsedToken(ctx context.Context, token, userID string) (*domain.RefreshToken, error)
	DeleteToken(ctx context.Context, token string) error
	DeleteTokensByUserID(ctx context.Context, userID string) error
}

type RefreshTokenInput struct {
	Token string `json:"token" validate:"required"`
}

type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error)
}
