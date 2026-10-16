package mock

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/stretchr/testify/mock"
)

type TokenIssuerMock struct {
	mock.Mock
	port.TokenIssuer // Embedded interface allows partial mocking (no need to implement all methods)
}

func (m *TokenIssuerMock) Generate(claims port.AccessTokenPayload) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}
