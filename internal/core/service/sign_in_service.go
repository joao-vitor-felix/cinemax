package service

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type SignInService struct {
	userRepo         port.UserRepository
	passwordHasher   port.PasswordHasher
	tokenIssuer      port.TokenIssuer
	refreshTokenRepo port.RefreshTokenRepository
}

func NewSignInService(
	userRepo port.UserRepository,
	passwordHasher port.PasswordHasher,
	tokenIssuer port.TokenIssuer,
	refreshTokenRepo port.RefreshTokenRepository,
) port.SignInService {
	return &SignInService{
		userRepo:         userRepo,
		passwordHasher:   passwordHasher,
		tokenIssuer:      tokenIssuer,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *SignInService) Execute(input port.SignInInput) (*port.SignInOutput, error) {
	user, err := s.userRepo.FindByEmail(input.Email)

	if err != nil {
		return nil, domain.InvalidCredentialsError
	}

	if user == nil {
		return nil, domain.InvalidCredentialsError
	}

	err = s.passwordHasher.Compare(user.PasswordHash, input.Password)
	if err != nil {
		return nil, domain.InvalidCredentialsError
	}

	accessToken, err := s.tokenIssuer.Generate(port.AccessTokenPayload{
		ID:    user.ID.String(),
		Email: user.Email,
	})

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshTokenRepo.GenerateToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &port.SignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}
