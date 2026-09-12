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


func setupSignInSut() (
	port.SignInUseCase,
	*m.UserRepositoryMock,
	*HasherMock,
	*m.TokenIssuerMock,
	*m.RefreshTokenRepositoryMock,
) {
	userRepo := new(m.UserRepositoryMock)
	hasher := new(HasherMock)
	tokenIssuer := new(m.TokenIssuerMock)
	refreshTokenRepo := new(m.RefreshTokenRepositoryMock)
	sut := usecase.NewSignInUseCase(userRepo, hasher, tokenIssuer, refreshTokenRepo)
	return sut, userRepo, hasher, tokenIssuer, refreshTokenRepo
}

func TestSignInUseCase(t *testing.T) {
	input := port.SignInInput{
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	userID := uuid.New()
	passwordHash := "hashed_password"
	mockUser := &domain.User{
		ID:           userID,
		Email:        input.Email,
		PasswordHash: &passwordHash,
	}

	t.Run("Execute", func(t *testing.T) {
		t.Run("should sign in user successfully", func(t *testing.T) {
			sut, userRepo, hasher, tokenIssuer, refreshTokenRepo := setupSignInSut()

			expectedAccessToken := "access_token"
			expectedRefreshToken := "refresh_token"

			userRepo.On("FindByEmail", mock.Anything, input.Email).Return(mockUser, nil).Once()
			hasher.On("Compare", mockUser.PasswordHash, input.Password).Return(nil).Once()
			tokenIssuer.On("Generate", port.AccessTokenPayload{
				ID:    userID.String(),
				Email: mockUser.Email,
			}).Return(expectedAccessToken, nil).Once()

			rt := &domain.RefreshToken{
				Token:     expectedRefreshToken,
				UserID:    userID.String(),
				ExpiresAt: time.Now().Add(time.Hour * 24),
			}
			refreshTokenRepo.On("GenerateToken", mock.Anything, userID.String()).Return(rt, nil).Once()

			output, err := sut.Execute(context.Background(), input)

			require.NoError(t, err)
			require.NotNil(t, output)
			require.Equal(t, expectedAccessToken, output.AccessToken)
			require.Equal(t, expectedRefreshToken, output.RefreshToken)

			userRepo.AssertExpectations(t)
			hasher.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
		})

		t.Run("should return error when user is not found", func(t *testing.T) {
			sut, userRepo, hasher, tokenIssuer, refreshTokenRepo := setupSignInSut()

			userRepo.On("FindByEmail", mock.Anything, input.Email).Return(nil, nil).Once()

			output, err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, domain.InvalidCredentialsError)
			require.Nil(t, output)

			userRepo.AssertExpectations(t)
			hasher.AssertNotCalled(t, "Compare")
			tokenIssuer.AssertNotCalled(t, "Generate")
			refreshTokenRepo.AssertNotCalled(t, "GenerateToken")
		})

		t.Run("should return error when password does not match", func(t *testing.T) {
			sut, userRepo, hasher, tokenIssuer, refreshTokenRepo := setupSignInSut()

			userRepo.On("FindByEmail", mock.Anything, input.Email).Return(mockUser, nil).Once()
			hasher.On("Compare", mockUser.PasswordHash, input.Password).Return(errors.New("invalid password")).Once()

			output, err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, domain.InvalidCredentialsError)
			require.Nil(t, output)

			userRepo.AssertExpectations(t)
			hasher.AssertExpectations(t)
			tokenIssuer.AssertNotCalled(t, "Generate")
			refreshTokenRepo.AssertNotCalled(t, "GenerateToken")
		})

		t.Run("should return error when access token generation fails", func(t *testing.T) {
			sut, userRepo, hasher, tokenIssuer, refreshTokenRepo := setupSignInSut()

			expectedErr := errors.New("token generation failed")

			userRepo.On("FindByEmail", mock.Anything, input.Email).Return(mockUser, nil).Once()
			hasher.On("Compare", mockUser.PasswordHash, input.Password).Return(nil).Once()
			tokenIssuer.On("Generate", mock.Anything).Return("", expectedErr).Once()

			output, err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, output)

			userRepo.AssertExpectations(t)
			hasher.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
			refreshTokenRepo.AssertNotCalled(t, "GenerateToken")
		})

		t.Run("should return error when refresh token generation fails", func(t *testing.T) {
			sut, userRepo, hasher, tokenIssuer, refreshTokenRepo := setupSignInSut()

			expectedErr := errors.New("refresh token generation failed")

			userRepo.On("FindByEmail", mock.Anything, input.Email).Return(mockUser, nil).Once()
			hasher.On("Compare", mockUser.PasswordHash, input.Password).Return(nil).Once()
			tokenIssuer.On("Generate", mock.Anything).Return("access_token", nil).Once()
			refreshTokenRepo.On("GenerateToken", mock.Anything, userID.String()).Return(nil, expectedErr).Once()

			output, err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, output)

			userRepo.AssertExpectations(t)
			hasher.AssertExpectations(t)
			tokenIssuer.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
		})
	})
}
