package service

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type RefreshTokenService struct {
	refreshTokenRepo port.RefreshTokenRepository
	userRepo         port.UserRepository
	tokenIssuer      port.TokenIssuer
}

func NewRefreshTokenService(
	refreshTokenRepo port.RefreshTokenRepository,
	userRepo port.UserRepository,
	tokenIssuer port.TokenIssuer,
) port.RefreshTokenService {
	return &RefreshTokenService{
		refreshTokenRepo: refreshTokenRepo,
		userRepo:         userRepo,
		tokenIssuer:      tokenIssuer,
	}
}

func (s *RefreshTokenService) Execute(input port.RefreshTokenInput) (*port.RefreshTokenOutput, error) {
	refreshToken, err := s.refreshTokenRepo.GetByToken(input.RefreshToken)

	if err != nil {
		return nil, err
	}

	if refreshToken == nil {
		return nil, domain.NotFoundError("token")
	}

	if refreshToken.IsUsed() {
		err = s.refreshTokenRepo.DeleteTokensByUserID(refreshToken.UserID)
		if err != nil {
			return nil, err
		}
		return nil, domain.InvalidCredentialsError
	}

	if refreshToken.IsExpired() {
		return nil, domain.InvalidCredentialsError
	}

	user, err := s.userRepo.FindByID(refreshToken.UserID)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.NotFoundError("user")
	}

	accessToken, err := s.tokenIssuer.Generate(port.AccessTokenPayload{
		ID:    user.ID.String(),
		Email: user.Email,
	})

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.refreshTokenRepo.GenerateAndInvalidateUsedToken(refreshToken.Token, user.ID.String())
	if err != nil {
		return nil, err
	}

	return &port.RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
	}, nil
}
