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

	if refreshToken == nil {
		return nil, domain.NotFoundError("token")
	}

	if err != nil {
		return nil, err
	}

	if refreshToken.IsUsed() {
		_ = s.refreshTokenRepo.DeleteTokensByUserID(refreshToken.UserId)
		return nil, domain.InvalidCredentialsError
	}

	if refreshToken.IsExpired() {
		return nil, domain.InvalidCredentialsError
	}

	user, err := s.userRepo.FindByID(refreshToken.UserId)

	if user == nil {
		return nil, domain.NotFoundError("user")
	}

	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenIssuer.Generate(port.AccessTokenPayload{
		Id:    user.ID.String(),
		Email: user.Email,
	})

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.refreshTokenRepo.GenerateAndDeleteUsedToken(refreshToken.Token, user.ID.String())
	if err != nil {
		return nil, err
	}

	return &port.RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
	}, nil
}
