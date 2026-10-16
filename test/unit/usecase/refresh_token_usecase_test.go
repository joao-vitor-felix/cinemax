package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/joao-vitor-felix/cinemax/internal/core/usecase"
	m "github.com/joao-vitor-felix/cinemax/test/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRefreshTokenSut() (
	port.RefreshTokenUseCase,
	*m.RefreshTokenRepositoryMock,
	*m.UserRepositoryMock,
	*m.TokenIssuerMock,
) {
	refreshTokenRepo := new(m.RefreshTokenRepositoryMock)
	userRepo := new(m.UserRepositoryMock)
	tokenIssuer := new(m.TokenIssuerMock)
	sut := usecase.NewRefreshTokenUseCase(refreshTokenRepo, userRepo, tokenIssuer)
	return sut, refreshTokenRepo, userRepo, tokenIssuer
}

func TestRefreshTokenUseCase(t *testing.T) {
	input := port.RefreshTokenInput{
		Token: "old_refresh_token",
	}

	userID := uuid.New()
	mockUser := &domain.User{
		ID:    userID,
		Email: "john@example.com",
	}

	validOldToken := &domain.RefreshToken{
		Token:     input.Token,
		UserID:    userID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}

	t.Run("Execute", func(t *testing.T) {
		t.Run("should refresh tokens successfully", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			expectedAccessToken := "new_access_token"
			expectedNewRefreshToken := "new_refresh_token"

			refreshTokenRepo.On("GetByToken", mock.Anything, input.Token).Return(validOldToken, nil).Once()
			userRepo.On("GetByID", mock.Anything, validOldToken.UserID).Return(mockUser, nil).Once()

			tokenIssuer.On("Generate", port.AccessTokenPayload{
				ID:    userID.String(),
				Email: mockUser.Email,
			}).Return(expectedAccessToken, nil).Once()

			newRt := &domain.RefreshToken{
				Token: expectedNewRefreshToken,
			}
			refreshTokenRepo.On("GenerateAndInvalidateUsedToken", mock.Anything, input.Token, userID.String()).Return(newRt, nil).Once()

			output, err := sut.Execute(context.Background(), input)

			require.NoError(t, err)
			require.NotNil(t, output)
			require.Equal(t, expectedAccessToken, output.AccessToken)
			require.Equal(t, expectedNewRefreshToken, output.RefreshToken)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
		})

		t.Run("should return error when token is not found", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			refreshTokenRepo.On("GetByToken", mock.Anything, input.Token).Return(nil, nil).Once()

			output, err := sut.Execute(context.Background(), input)

			require.Error(t, err)
			require.Nil(t, output)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "TOKEN_NOT_FOUND", appErr.Code)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "GetByID")
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error and clear tokens when old token was already used", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			usedTime := time.Now().Add(-1 * time.Hour)
			usedToken := &domain.RefreshToken{
				Token:     input.Token,
				UserID:    userID.String(),
				ExpiresAt: time.Now().Add(time.Hour * 1),
				UsedAt:    &usedTime,
			}

			refreshTokenRepo.On("GetByToken", mock.Anything, input.Token).Return(usedToken, nil).Once()
			refreshTokenRepo.On("DeleteTokensByUserID", mock.Anything, userID.String()).Return(nil).Once()

			output, err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, domain.InvalidCredentialsError)
			require.Nil(t, output)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "GetByID")
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error when old token is expired", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			expiredToken := &domain.RefreshToken{
				Token:     input.Token,
				UserID:    userID.String(),
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			}

			refreshTokenRepo.On("GetByToken", mock.Anything, input.Token).Return(expiredToken, nil).Once()

			output, err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, domain.InvalidCredentialsError)
			require.Nil(t, output)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertNotCalled(t, "GetByID")
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error when user is not found", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			refreshTokenRepo.On("GetByToken", mock.Anything, input.Token).Return(validOldToken, nil).Once()
			userRepo.On("GetByID", mock.Anything, validOldToken.UserID).Return(nil, nil).Once()

			output, err := sut.Execute(context.Background(), input)

			require.Error(t, err)
			require.Nil(t, output)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "USER_NOT_FOUND", appErr.Code)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertNotCalled(t, "Generate")
		})

		t.Run("should return error when access token generation fails", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			expectedErr := errors.New("token gen fail")

			refreshTokenRepo.On("GetByToken", mock.Anything, input.Token).Return(validOldToken, nil).Once()
			userRepo.On("GetByID", mock.Anything, validOldToken.UserID).Return(mockUser, nil).Once()

			tokenIssuer.On("Generate", port.AccessTokenPayload{
				ID:    userID.String(),
				Email: mockUser.Email,
			}).Return("", expectedErr).Once()

			output, err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, output)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
		})

		t.Run("should return error when replacing old token fails", func(t *testing.T) {
			sut, refreshTokenRepo, userRepo, tokenIssuer := setupRefreshTokenSut()

			expectedErr := errors.New("db error replacing token")

			refreshTokenRepo.On("GetByToken", mock.Anything, input.Token).Return(validOldToken, nil).Once()
			userRepo.On("GetByID", mock.Anything, validOldToken.UserID).Return(mockUser, nil).Once()

			tokenIssuer.On("Generate", port.AccessTokenPayload{
				ID:    userID.String(),
				Email: mockUser.Email,
			}).Return("at_123", nil).Once()

			refreshTokenRepo.On("GenerateAndInvalidateUsedToken", mock.Anything, input.Token, userID.String()).Return(nil, expectedErr).Once()

			output, err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, output)

			refreshTokenRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
		})
	})
}
