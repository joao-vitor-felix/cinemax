package usecase

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

const googleProvider = "google"

type SignInGoogleUseCase struct {
	oauthService     port.OAuthService
	userRepo         port.UserRepository
	tokenIssuer      port.TokenIssuer
	refreshTokenRepo port.RefreshTokenRepository
}

func NewSignInGoogleUseCase(
	oauthService port.OAuthService,
	userRepo port.UserRepository,
	tokenIssuer port.TokenIssuer,
	refreshTokenRepo port.RefreshTokenRepository,
) *SignInGoogleUseCase {
	return &SignInGoogleUseCase{
		oauthService:     oauthService,
		userRepo:         userRepo,
		tokenIssuer:      tokenIssuer,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (u *SignInGoogleUseCase) Execute(ctx context.Context, input port.SignInGoogleInput) (*port.SignInOutput, error) {
	accessToken, err := u.oauthService.GetAccessToken(input.Code)
	if err != nil {
		return nil, domain.InvalidCredentialsError
	}

	userInfo, err := u.oauthService.GetUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	if !userInfo.VerifiedEmail {
		return nil, domain.EmailNotVerifiedError
	}

	user, err := u.resolveUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	appAccessToken, err := u.tokenIssuer.Generate(port.AccessTokenPayload{
		ID:    user.ID.String(),
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.refreshTokenRepo.GenerateToken(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	return &port.SignInOutput{
		AccessToken:  appAccessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (u *SignInGoogleUseCase) resolveUser(ctx context.Context, userInfo *port.UserInfo) (*domain.User, error) {
	userFromProviderID, err := u.userRepo.GetByProviderUserID(ctx, googleProvider, userInfo.ID)
	if err != nil {
		return nil, err
	}

	if userFromProviderID != nil {
		return userFromProviderID, nil
	}

	userFromDb, err := u.userRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, err
	}

	if userFromDb == nil {
		newUser := domain.NewOAuthUser(userInfo.GivenName, userInfo.FamilyName, userInfo.Email)
		newIdentity := &domain.FederatedIdentity{
			Provider:       googleProvider,
			ProviderUserID: userInfo.ID,
		}

		userID, err := u.userRepo.CreateWithIdentity(ctx, newUser.FirstName, newUser.LastName, newUser.Email, newIdentity)
		if err != nil {
			return nil, err
		}

		return &domain.User{
			ID:        userID,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newUser.Email,
		}, nil
	}

	if _, err := u.userRepo.LinkIdentity(ctx, userFromDb.ID.String(), googleProvider, userInfo.ID); err != nil {
		return nil, err
	}

	return userFromDb, nil
}
