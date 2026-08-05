package mock

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) IsContactInfoAvailable(ctx context.Context, email, phone string) (bool, error) {
	args := m.Called(ctx, email, phone)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) CreateWithIdentity(ctx context.Context, user *domain.User, identity *domain.FederatedIdentity) error {
	args := m.Called(ctx, user, identity)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindByProviderUserID(ctx context.Context, provider, providerUserID string) (*domain.User, error) {
	args := m.Called(ctx, provider, providerUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) LinkIdentity(ctx context.Context, userID, provider, providerUserID string) (*domain.FederatedIdentity, error) {
	args := m.Called(ctx, userID, provider, providerUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FederatedIdentity), args.Error(1)
}
