package mock

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) IsContactInfoAvailable(email, phone string) (bool, error) {
	args := m.Called(email, phone)
	return args.Bool(0), args.Error(1)
}
