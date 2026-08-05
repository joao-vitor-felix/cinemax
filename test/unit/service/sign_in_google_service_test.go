package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/joao-vitor-felix/cinemax/internal/core/service"
	appMock "github.com/joao-vitor-felix/cinemax/test/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSignInGoogleTest() (
	*appMock.OAuthServiceMock,
	*appMock.UserRepositoryMock,
	*TokenIssuerMock,
	*appMock.RefreshTokenRepositoryMock,
	port.SignInGoogleService,
) {
	oauthService := new(appMock.OAuthServiceMock)
	userRepo := new(appMock.UserRepositoryMock)
	tokenIssuer := new(TokenIssuerMock)
	refreshTokenRepo := new(appMock.RefreshTokenRepositoryMock)

	svc := service.NewSignInGoogleService(oauthService, userRepo, tokenIssuer, refreshTokenRepo)

	return oauthService, userRepo, tokenIssuer, refreshTokenRepo, svc
}

func TestSignInGoogleService_SuccessNewUser(t *testing.T) {
	oauthService, userRepo, tokenIssuer, refreshTokenRepo, svc := setupSignInGoogleTest()

	code := "valid_code"
	accessToken := "google_access_token"
	googleID := "12345"
	email := "test@example.com"
	appAccessToken := "app_access_token"

	oauthService.On("GetAccessToken", code).Return(accessToken, nil)
	oauthService.On("GetUserInfo", accessToken).Return(&port.UserInfo{
		ID:            googleID,
		Email:         email,
		VerifiedEmail: true,
		GivenName:     "John",
		FamilyName:    "Doe",
		Picture:       "http://example.com/pic.jpg",
	}, nil)

	userRepo.On("FindByProviderUserID", mock.Anything, "google", googleID).Return(nil, nil)
	userRepo.On("FindByEmail", mock.Anything, email).Return(nil, nil)
	userRepo.On("CreateWithIdentity", mock.Anything, mock.AnythingOfType("*domain.User"), mock.AnythingOfType("*domain.FederatedIdentity")).Return(nil)

	tokenIssuer.On("Generate", mock.AnythingOfType("port.AccessTokenPayload")).Return(appAccessToken, nil)
	refreshTokenRepo.On("GenerateToken", mock.Anything, mock.AnythingOfType("string")).Return(&domain.RefreshToken{Token: "app_refresh_token"}, nil)

	output, err := svc.Execute(context.Background(), port.SignInGoogleInput{Code: code})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, appAccessToken, output.AccessToken)
	assert.Equal(t, "app_refresh_token", output.RefreshToken)
}

func TestSignInGoogleService_SuccessExistingLinkedUser(t *testing.T) {
	oauthService, userRepo, tokenIssuer, refreshTokenRepo, svc := setupSignInGoogleTest()

	code := "valid_code"
	accessToken := "google_access_token"
	googleID := "12345"
	email := "test@example.com"
	userID := uuid.New()
	appAccessToken := "app_access_token"

	oauthService.On("GetAccessToken", code).Return(accessToken, nil)
	oauthService.On("GetUserInfo", accessToken).Return(&port.UserInfo{
		ID:            googleID,
		Email:         email,
		VerifiedEmail: true,
	}, nil)

	user := &domain.User{ID: userID, Email: email}
	userRepo.On("FindByProviderUserID", mock.Anything, "google", googleID).Return(user, nil)

	tokenIssuer.On("Generate", mock.AnythingOfType("port.AccessTokenPayload")).Return(appAccessToken, nil)
	refreshTokenRepo.On("GenerateToken", mock.Anything, userID.String()).Return(&domain.RefreshToken{Token: "app_refresh_token"}, nil)

	output, err := svc.Execute(context.Background(), port.SignInGoogleInput{Code: code})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, appAccessToken, output.AccessToken)
}

func TestSignInGoogleService_SuccessExistingUnlinkedUser(t *testing.T) {
	oauthService, userRepo, tokenIssuer, refreshTokenRepo, svc := setupSignInGoogleTest()

	code := "valid_code"
	accessToken := "google_access_token"
	googleID := "12345"
	email := "test@example.com"
	userID := uuid.New()
	appAccessToken := "app_access_token"

	oauthService.On("GetAccessToken", code).Return(accessToken, nil)
	oauthService.On("GetUserInfo", accessToken).Return(&port.UserInfo{
		ID:            googleID,
		Email:         email,
		VerifiedEmail: true,
	}, nil)

	userRepo.On("FindByProviderUserID", mock.Anything, "google", googleID).Return(nil, nil)

	user := &domain.User{ID: userID, Email: email}
	userRepo.On("FindByEmail", mock.Anything, email).Return(user, nil)

	userRepo.On("LinkIdentity", mock.Anything, userID.String(), "google", googleID).Return(&domain.FederatedIdentity{}, nil)

	tokenIssuer.On("Generate", mock.AnythingOfType("port.AccessTokenPayload")).Return(appAccessToken, nil)
	refreshTokenRepo.On("GenerateToken", mock.Anything, userID.String()).Return(&domain.RefreshToken{Token: "app_refresh_token"}, nil)

	output, err := svc.Execute(context.Background(), port.SignInGoogleInput{Code: code})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, appAccessToken, output.AccessToken)
}

func TestSignInGoogleService_EmailNotVerified(t *testing.T) {
	oauthService, _, _, _, svc := setupSignInGoogleTest()

	code := "valid_code"
	accessToken := "google_access_token"

	oauthService.On("GetAccessToken", code).Return(accessToken, nil)
	oauthService.On("GetUserInfo", accessToken).Return(&port.UserInfo{
		VerifiedEmail: false,
	}, nil)

	output, err := svc.Execute(context.Background(), port.SignInGoogleInput{Code: code})

	assert.Error(t, err)
	assert.Equal(t, domain.EmailNotVerifiedError, err)
	assert.Nil(t, output)
}

func TestSignInGoogleService_GetAccessTokenError(t *testing.T) {
	oauthService, _, _, _, svc := setupSignInGoogleTest()

	code := "invalid_code"
	oauthService.On("GetAccessToken", code).Return("", errors.New("invalid_grant"))

	output, err := svc.Execute(context.Background(), port.SignInGoogleInput{Code: code})

	assert.Error(t, err)
	assert.Equal(t, domain.InvalidCredentialsError, err)
	assert.Nil(t, output)
}
