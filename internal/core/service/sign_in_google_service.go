package service

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

const googleProvider = "google"

type SignInGoogleService struct {
	oauthService     port.OAuthService
	userRepo         port.UserRepository
	tokenIssuer      port.TokenIssuer
	refreshTokenRepo port.RefreshTokenRepository
}

func NewSignInGoogleService(
	oauthService port.OAuthService,
	userRepo port.UserRepository,
	tokenIssuer port.TokenIssuer,
	refreshTokenRepo port.RefreshTokenRepository,
) port.SignInGoogleService {
	return &SignInGoogleService{
		oauthService:     oauthService,
		userRepo:         userRepo,
		tokenIssuer:      tokenIssuer,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *SignInGoogleService) Execute(ctx context.Context, input port.SignInGoogleInput) (*port.SignInOutput, error) {
	accessToken, err := s.oauthService.GetAccessToken(input.Code)
	if err != nil {
		return nil, domain.InvalidCredentialsError
	}

	userInfo, err := s.oauthService.GetUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	if !userInfo.VerifiedEmail {
		return nil, domain.EmailNotVerifiedError
	}

	user, err := s.resolveUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	appAccessToken, err := s.tokenIssuer.Generate(port.AccessTokenPayload{
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
		AccessToken:  appAccessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *SignInGoogleService) resolveUser(ctx context.Context, userInfo *port.UserInfo) (*domain.User, error) {
	userFromProviderID, err := s.userRepo.FindByProviderUserID(ctx, googleProvider, userInfo.ID)
	if err != nil {
		return nil, err
	}

	if userFromProviderID != nil {
		return userFromProviderID, nil
	}

	userFromDb, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, err
	}

	if userFromDb == nil {
		newUser := domain.NewOAuthUser(userInfo.GivenName, userInfo.FamilyName, userInfo.Email)
		newIdentity := &domain.FederatedIdentity{
			Provider:       googleProvider,
			ProviderUserID: userInfo.ID,
		}

		if err := s.userRepo.CreateWithIdentity(ctx, newUser, newIdentity); err != nil {
			return nil, err
		}

		return newUser, nil
	}

	if _, err := s.userRepo.LinkIdentity(ctx, userFromDb.ID.String(), googleProvider, userInfo.ID); err != nil {
		return nil, err
	}

	return userFromDb, nil
}
