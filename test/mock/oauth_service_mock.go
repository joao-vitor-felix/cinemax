package mock

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/stretchr/testify/mock"
)

type OAuthServiceMock struct {
	mock.Mock
}

func (m *OAuthServiceMock) GetAccessToken(code string) (string, error) {
	args := m.Called(code)
	return args.String(0), args.Error(1)
}

func (m *OAuthServiceMock) GetUserInfo(accessToken string) (*port.UserInfo, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.UserInfo), args.Error(1)
}

func (m *OAuthServiceMock) RevokeAccessToken(accessToken string) error {
	args := m.Called(accessToken)
	return args.Error(0)
}
