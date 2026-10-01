package mock

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type PasswordResetRepositoryMock struct {
	mock.Mock
}

func (m *PasswordResetRepositoryMock) GetByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	args := m.Called(ctx, token)
	var rt *domain.PasswordResetToken
	if args.Get(0) != nil {
		rt = args.Get(0).(*domain.PasswordResetToken)
	}
	return rt, args.Error(1)
}

func (m *PasswordResetRepositoryMock) SaveToken(ctx context.Context, tokenHash, userID string) (*domain.PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash, userID)
	var rt *domain.PasswordResetToken
	if args.Get(0) != nil {
		rt = args.Get(0).(*domain.PasswordResetToken)
	}
	return rt, args.Error(1)
}

func (m *PasswordResetRepositoryMock) MarkTokenAsUsed(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
