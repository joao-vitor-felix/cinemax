package mock

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type RefreshTokenRepositoryMock struct {
	mock.Mock
}

func (r *RefreshTokenRepositoryMock) GetByToken(token string) (*domain.RefreshToken, error) {
	args := r.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (r *RefreshTokenRepositoryMock) GenerateToken(userId string) (*domain.RefreshToken, error) {
	args := r.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (r *RefreshTokenRepositoryMock) GenerateAndDeleteUsedToken(token, userId string) (*domain.RefreshToken, error) {
	args := r.Called(token, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (r *RefreshTokenRepositoryMock) DeleteToken(token string) error {
	args := r.Called(token)
	return args.Error(0)
}

func (r *RefreshTokenRepositoryMock) DeleteTokensByUserID(userId string) error {
	args := r.Called(userId)
	return args.Error(0)
}
