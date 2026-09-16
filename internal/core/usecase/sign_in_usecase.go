package usecase

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type SignInUseCase struct {
	userRepo         port.UserRepository
	passwordHasher   port.PasswordHasher
	tokenIssuer      port.TokenIssuer
	refreshTokenRepo port.RefreshTokenRepository
}

func NewSignInUseCase(
	userRepo port.UserRepository,
	passwordHasher port.PasswordHasher,
	tokenIssuer port.TokenIssuer,
	refreshTokenRepo port.RefreshTokenRepository,
) port.SignInUseCase {
	return &SignInUseCase{
		userRepo:         userRepo,
		passwordHasher:   passwordHasher,
		tokenIssuer:      tokenIssuer,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *SignInUseCase) Execute(ctx context.Context, input port.SignInInput) (*port.SignInOutput, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)

	if err != nil {
		return nil, err
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

	refreshToken, err := s.refreshTokenRepo.GenerateToken(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	return &port.SignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}
