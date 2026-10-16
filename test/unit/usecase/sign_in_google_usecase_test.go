package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/joao-vitor-felix/cinemax/internal/core/usecase"
	m "github.com/joao-vitor-felix/cinemax/test/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupSignInGoogleSut() (
	port.SignInGoogleUseCase,
	*m.OAuthServiceMock,
	*m.UserRepositoryMock,
	*m.TokenIssuerMock,
	*m.RefreshTokenRepositoryMock,
) {
	oauthService := new(m.OAuthServiceMock)
	userRepo := new(m.UserRepositoryMock)
	tokenIssuer := new(m.TokenIssuerMock)
	refreshTokenRepo := new(m.RefreshTokenRepositoryMock)

	sut := usecase.NewSignInGoogleUseCase(oauthService, userRepo, tokenIssuer, refreshTokenRepo)

	return sut, oauthService, userRepo, tokenIssuer, refreshTokenRepo
}

func TestSignInGoogleUseCase(t *testing.T) {
	code := "valid_code"
	accessToken := "google_access_token"
	googleID := "12345"
	email := "test@example.com"

	t.Run("Execute", func(t *testing.T) {
		t.Run("should sign in a new user successfully", func(t *testing.T) {
			sut, oauthService, userRepo, tokenIssuer, refreshTokenRepo := setupSignInGoogleSut()

			appAccessToken := "app_access_token"
			expectedRefreshToken := "app_refresh_token"
			createdUserID := uuid.New()

			oauthService.On("GetAccessToken", code).Return(accessToken, nil).Once()
			oauthService.On("GetUserInfo", accessToken).Return(&port.UserInfo{
				ID:            googleID,
				Email:         email,
				VerifiedEmail: true,
				GivenName:     "John",
				FamilyName:    "Doe",
				Picture:       "http://example.com/pic.jpg",
			}, nil).Once()

			userRepo.On("GetByProviderUserID", mock.Anything, "google", googleID).Return(nil, nil).Once()
			userRepo.On("GetByEmail", mock.Anything, email).Return(nil, nil).Once()
			userRepo.On("CreateWithIdentity", mock.Anything, "John", "Doe", email, mock.AnythingOfType("*domain.FederatedIdentity")).Return(createdUserID, nil).Once()

			tokenIssuer.On("Generate", mock.AnythingOfType("port.AccessTokenPayload")).Return(appAccessToken, nil).Once()
			refreshTokenRepo.On("GenerateToken", mock.Anything, createdUserID.String()).Return(&domain.RefreshToken{Token: expectedRefreshToken}, nil).Once()

			output, err := sut.Execute(context.Background(), port.SignInGoogleInput{Code: code})

			require.NoError(t, err)
			require.NotNil(t, output)
			require.Equal(t, appAccessToken, output.AccessToken)
			require.Equal(t, expectedRefreshToken, output.RefreshToken)

			oauthService.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
		})

		t.Run("should sign in an existing linked user successfully", func(t *testing.T) {
			sut, oauthService, userRepo, tokenIssuer, refreshTokenRepo := setupSignInGoogleSut()

			userID := uuid.New()
			appAccessToken := "app_access_token"
			expectedRefreshToken := "app_refresh_token"

			oauthService.On("GetAccessToken", code).Return(accessToken, nil).Once()
			oauthService.On("GetUserInfo", accessToken).Return(&port.UserInfo{
				ID:            googleID,
				Email:         email,
				VerifiedEmail: true,
			}, nil).Once()

			user := &domain.User{ID: userID, Email: email}
			userRepo.On("GetByProviderUserID", mock.Anything, "google", googleID).Return(user, nil).Once()

			tokenIssuer.On("Generate", mock.AnythingOfType("port.AccessTokenPayload")).Return(appAccessToken, nil).Once()
			refreshTokenRepo.On("GenerateToken", mock.Anything, userID.String()).Return(&domain.RefreshToken{Token: expectedRefreshToken}, nil).Once()

			output, err := sut.Execute(context.Background(), port.SignInGoogleInput{Code: code})

			require.NoError(t, err)
			require.NotNil(t, output)
			require.Equal(t, appAccessToken, output.AccessToken)
			require.Equal(t, expectedRefreshToken, output.RefreshToken)

			oauthService.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
		})

		t.Run("should sign in an existing unlinked user successfully", func(t *testing.T) {
			sut, oauthService, userRepo, tokenIssuer, refreshTokenRepo := setupSignInGoogleSut()

			userID := uuid.New()
			appAccessToken := "app_access_token"
			expectedRefreshToken := "app_refresh_token"

			oauthService.On("GetAccessToken", code).Return(accessToken, nil).Once()
			oauthService.On("GetUserInfo", accessToken).Return(&port.UserInfo{
				ID:            googleID,
				Email:         email,
				VerifiedEmail: true,
			}, nil).Once()

			userRepo.On("GetByProviderUserID", mock.Anything, "google", googleID).Return(nil, nil).Once()

			user := &domain.User{ID: userID, Email: email}
			userRepo.On("GetByEmail", mock.Anything, email).Return(user, nil).Once()

			userRepo.On("LinkIdentity", mock.Anything, userID.String(), "google", googleID).Return(&domain.FederatedIdentity{}, nil).Once()

			tokenIssuer.On("Generate", mock.AnythingOfType("port.AccessTokenPayload")).Return(appAccessToken, nil).Once()
			refreshTokenRepo.On("GenerateToken", mock.Anything, userID.String()).Return(&domain.RefreshToken{Token: expectedRefreshToken}, nil).Once()

			output, err := sut.Execute(context.Background(), port.SignInGoogleInput{Code: code})

			require.NoError(t, err)
			require.NotNil(t, output)
			require.Equal(t, appAccessToken, output.AccessToken)
			require.Equal(t, expectedRefreshToken, output.RefreshToken)

			oauthService.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
		})

		t.Run("should return error when email is not verified", func(t *testing.T) {
			sut, oauthService, userRepo, tokenIssuer, refreshTokenRepo := setupSignInGoogleSut()

			oauthService.On("GetAccessToken", code).Return(accessToken, nil).Once()
			oauthService.On("GetUserInfo", accessToken).Return(&port.UserInfo{
				VerifiedEmail: false,
			}, nil).Once()

			output, err := sut.Execute(context.Background(), port.SignInGoogleInput{Code: code})

			require.Error(t, err)
			require.Equal(t, domain.EmailNotVerifiedError, err)
			require.Nil(t, output)

			oauthService.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "GetByProviderUserID")
			tokenIssuer.AssertNotCalled(t, "Generate")
			refreshTokenRepo.AssertNotCalled(t, "GenerateToken")
		})

		t.Run("should return error when access token retrieval fails", func(t *testing.T) {
			sut, oauthService, userRepo, tokenIssuer, refreshTokenRepo := setupSignInGoogleSut()

			invalidCode := "invalid_code"
			oauthService.On("GetAccessToken", invalidCode).Return("", errors.New("invalid_grant")).Once()

			output, err := sut.Execute(context.Background(), port.SignInGoogleInput{Code: invalidCode})

			require.Error(t, err)
			require.Equal(t, domain.InvalidCredentialsError, err)
			require.Nil(t, output)

			oauthService.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "GetByProviderUserID")
			tokenIssuer.AssertNotCalled(t, "Generate")
			refreshTokenRepo.AssertNotCalled(t, "GenerateToken")
		})
	})
}
