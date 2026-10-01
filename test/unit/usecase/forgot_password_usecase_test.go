package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/joao-vitor-felix/cinemax/internal/core/usecase"
	cinemaxmock "github.com/joao-vitor-felix/cinemax/test/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestForgotPasswordUseCase(t *testing.T) {
	setupSut := func() (*usecase.ForgotPasswordUseCase, *cinemaxmock.PasswordResetRepositoryMock, *cinemaxmock.UserRepositoryMock) {
		passwordResetRepo := new(cinemaxmock.PasswordResetRepositoryMock)
		userRepo := new(cinemaxmock.UserRepositoryMock)
		sut := usecase.NewPasswordResetUseCase(passwordResetRepo, userRepo)
		return sut, passwordResetRepo, userRepo
	}

	input := port.ForgotPasswordInput{
		Email: "john.doe@example.com",
	}

	t.Run("Execute", func(t *testing.T) {
		t.Run("should successfully execute and generate a token", func(t *testing.T) {
			sut, passwordResetRepo, userRepo := setupSut()

			userID := uuid.New()
			user := &domain.User{
				ID:    userID,
				Email: input.Email,
			}

			userRepo.On("GetByEmail", mock.Anything, input.Email).Return(user, nil).Once()
			passwordResetRepo.On("SaveToken", mock.Anything, mock.AnythingOfType("string"), userID.String()).Return(&domain.PasswordResetToken{}, nil).Once()

			err := sut.Execute(context.Background(), input)

			require.NoError(t, err)
			userRepo.AssertExpectations(t)
			passwordResetRepo.AssertExpectations(t)
		})

		t.Run("should return error when GetByEmail fails", func(t *testing.T) {
			sut, passwordResetRepo, userRepo := setupSut()

			expectedErr := errors.New("db error")
			userRepo.On("GetByEmail", mock.Anything, input.Email).Return(nil, expectedErr).Once()

			err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, expectedErr)
			userRepo.AssertExpectations(t)
			passwordResetRepo.AssertNotCalled(t, "SaveToken")
		})

		t.Run("should return NotFoundError when user does not exist", func(t *testing.T) {
			sut, passwordResetRepo, userRepo := setupSut()

			userRepo.On("GetByEmail", mock.Anything, input.Email).Return(nil, nil).Once()

			err := sut.Execute(context.Background(), input)

			require.Error(t, err)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "USER_NOT_FOUND", appErr.Code)

			userRepo.AssertExpectations(t)
			passwordResetRepo.AssertNotCalled(t, "SaveToken")
		})

		t.Run("should return error when SaveToken fails", func(t *testing.T) {
			sut, passwordResetRepo, userRepo := setupSut()

			userID := uuid.New()
			user := &domain.User{
				ID:    userID,
				Email: input.Email,
			}

			expectedErr := errors.New("db error on save")
			userRepo.On("GetByEmail", mock.Anything, input.Email).Return(user, nil).Once()
			passwordResetRepo.On("SaveToken", mock.Anything, mock.AnythingOfType("string"), userID.String()).Return(nil, expectedErr).Once()

			err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, expectedErr)
			userRepo.AssertExpectations(t)
			passwordResetRepo.AssertExpectations(t)
		})
	})
}
