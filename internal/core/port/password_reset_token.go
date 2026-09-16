package port

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type PasswordResetTokenRepository interface {
	GetByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error)
	GenerateToken(ctx context.Context, tokenHash, userID string) (*domain.PasswordResetToken, error)
	MarkTokenAsUsed(ctx context.Context, token string) error
}
