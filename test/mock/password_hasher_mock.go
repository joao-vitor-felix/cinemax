package mock

import (
	"github.com/stretchr/testify/mock"
)

type PasswordHasherMock struct {
	mock.Mock
}

func (h *PasswordHasherMock) Hash(password string) (string, error) {
	args := h.Called(password)
	return args.String(0), args.Error(1)
}

func (h *PasswordHasherMock) Compare(hash *string, password string) (bool, error) {
	args := h.Called(hash, password)
	return args.Bool(0), args.Error(1)
}
