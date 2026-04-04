package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/joao-vitor-felix/cinemax/internal/core/service"
	m "github.com/joao-vitor-felix/cinemax/test/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TokenIssuerServiceMock struct {
	mock.Mock
}

func (t *TokenIssuerServiceMock) Generate(claims port.AccessTokenPayload) (string, error) {
	args := t.Called(claims)
	return args.String(0), args.Error(1)
}

func (t *TokenIssuerServiceMock) Validate(token string) (*port.AccessTokenPayload, error) {
	args := t.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.AccessTokenPayload), args.Error(1)
}

func setupRefreshTokenSut() (
	port.RefreshTokenService,
	*m.RefreshTokenRepositoryMock,
	*m.UserRepositoryMock,
	*TokenIssuerServiceMock,
) {
	refreshTokenRepo := new(m.RefreshTokenRepositoryMock)
	userRepo := new(m.UserRepositoryMock)
	tokenIssuer := new(TokenIssuerServiceMock)
	sut := service.NewRefreshTokenService(refreshTokenRepo, userRepo, tokenIssuer)
	return sut, refreshTokenRepo, userRepo, tokenIssuer
}

func TestRefreshTokenService(t *testing.T) {
	input := port.RefreshTokenInput{
		RefreshToken: "old_refresh_token",
	}

	userId := uuid.New()
	mockUser := &domain.User{
		ID:    userId,
		Email: "john@example.com",
	}

	validOldToken := &domain.RefreshToken{
		Token:     input.RefreshToken,
		UserID:    userId.String(),
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}

	t.Run("Execute", func(t *testing.T) {
		t.Run("should refresh tokens successfully", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			expectedAccessToken := "new_access_token"
			expectedNewRefreshToken := "new_refresh_token"

			refreshTokenRepo.On("GetByToken", input.RefreshToken).Return(validOldToken, nil).Once()
			userRepo.On("FindByID", validOldToken.UserID).Return(mockUser, nil).Once()

			tokenIssuer.On("Generate", port.AccessTokenPayload{
				ID:    userId.String(),
				Email: mockUser.Email,
			}).Return(expectedAccessToken, nil).Once()

			newRt := &domain.RefreshToken{
				Token: expectedNewRefreshToken,
			}
			refreshTokenRepo.On("GenerateAndInvalidateUsedToken", input.RefreshToken, userId.String()).Return(newRt, nil).Once()

			output, err := sut.Execute(input)

			require.NoError(t, err)
			require.NotNil(t, output)
			require.Equal(t, expectedAccessToken, output.AccessToken)
			require.Equal(t, expectedNewRefreshToken, output.RefreshToken)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
		})

		t.Run("should return error when old token is not found", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			refreshTokenRepo.On("GetByToken", input.RefreshToken).Return(nil, nil).Once()

			output, err := sut.Execute(input)

			require.Error(t, err)
			require.Nil(t, output)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "NOT_FOUND", appErr.Code)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "FindByID")
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error and clear tokens when old token was already used", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			usedTime := time.Now().Add(-1 * time.Hour)
			usedToken := &domain.RefreshToken{
				Token:     input.RefreshToken,
				UserID:    userId.String(),
				ExpiresAt: time.Now().Add(time.Hour * 1),
				UsedAt:    &usedTime,
			}

			refreshTokenRepo.On("GetByToken", input.RefreshToken).Return(usedToken, nil).Once()
			refreshTokenRepo.On("DeleteTokensByUserID", userId.String()).Return(nil).Once()

			output, err := sut.Execute(input)

			require.ErrorIs(t, err, domain.InvalidCredentialsError)
			require.Nil(t, output)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "FindByID")
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error when old token is expired", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			expiredToken := &domain.RefreshToken{
				Token:     input.RefreshToken,
				UserID:    userId.String(),
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			}

			refreshTokenRepo.On("GetByToken", input.RefreshToken).Return(expiredToken, nil).Once()

			output, err := sut.Execute(input)

			require.ErrorIs(t, err, domain.InvalidCredentialsError)
			require.Nil(t, output)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "FindByID")
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error when token is not found", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			refreshTokenRepo.On("GetByToken", input.RefreshToken).Return(nil, nil).Once()

			output, err := sut.Execute(input)

			require.Error(t, err)
			require.Nil(t, output)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "NOT_FOUND", appErr.Code)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "FindByID")
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error when user is not found", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			refreshTokenRepo.On("GetByToken", input.RefreshToken).Return(validOldToken, nil).Once()
			userRepo.On("FindByID", validOldToken.UserID).Return(nil, nil).Once()

			output, err := sut.Execute(input)

			require.Error(t, err)
			require.Nil(t, output)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "NOT_FOUND", appErr.Code)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error when access token generation fails", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			expectedErr := errors.New("token gen fail")

			refreshTokenRepo.On("GetByToken", input.RefreshToken).Return(validOldToken, nil).Once()
			userRepo.On("FindByID", validOldToken.UserID).Return(mockUser, nil).Once()

			tokenIssuer.On("Generate", port.AccessTokenPayload{
				ID:    userId.String(),
				Email: mockUser.Email,
			}).Return("", expectedErr).Once()

			output, err := sut.Execute(input)

			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, output)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
		})

		t.Run("should return error when replacing old token fails", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			expectedErr := errors.New("db error replacing token")

			refreshTokenRepo.On("GetByToken", input.RefreshToken).Return(validOldToken, nil).Once()
			userRepo.On("FindByID", validOldToken.UserID).Return(mockUser, nil).Once()

			tokenIssuer.On("Generate", port.AccessTokenPayload{
				ID:    userId.String(),
				Email: mockUser.Email,
			}).Return("at_123", nil).Once()

			refreshTokenRepo.On("GenerateAndInvalidateUsedToken", input.RefreshToken, userId.String()).Return(nil, expectedErr).Once()

			output, err := sut.Execute(input)

			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, output)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
		})
	})
}
