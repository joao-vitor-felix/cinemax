package port

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type PasswordResetRepository interface {
	GetByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error)
	SaveToken(ctx context.Context, tokenHash, userID string) (*domain.PasswordResetToken, error)
	MarkTokenAsUsed(ctx context.Context, token string) error
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

// type PasswordResetUseCase interface {
// 	Execute(ctx context.Context, userID string, input PasswordResetTokenInput) error
// }

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, input ForgotPasswordInput) error
}
