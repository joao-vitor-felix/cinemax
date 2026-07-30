package mock

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type RefreshTokenRepositoryMock struct {
	mock.Mock
}

func (r *RefreshTokenRepositoryMock) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := r.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (r *RefreshTokenRepositoryMock) GenerateToken(ctx context.Context, userId string) (*domain.RefreshToken, error) {
	args := r.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (r *RefreshTokenRepositoryMock) GenerateAndInvalidateUsedToken(ctx context.Context, token, userId string) (*domain.RefreshToken, error) {
	args := r.Called(ctx, token, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (r *RefreshTokenRepositoryMock) DeleteToken(ctx context.Context, token string) error {
	args := r.Called(ctx, token)
	return args.Error(0)
}

func (r *RefreshTokenRepositoryMock) DeleteTokensByUserID(ctx context.Context, userId string) error {
	args := r.Called(ctx, userId)
	return args.Error(0)
}
