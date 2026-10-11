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

type PasswordHasherMock struct {
	mock.Mock
}

func (h *PasswordHasherMock) Hash(password string) (string, error) {
	args := h.Called(password)
	return args.String(0), args.Error(1)
}

func (h *PasswordHasherMock) Compare(hash *string, password string) error {
	args := h.Called(hash, password)
	return args.Error(0)
}

func setupResetPasswordSut() (
	port.ResetPasswordUseCase,
	*m.PasswordResetRepositoryMock,
	*m.UserRepositoryMock,
	*PasswordHasherMock,
	*m.UnitOfWorkMock,
) {
	passwordResetRepo := new(m.PasswordResetRepositoryMock)
	userRepo := new(m.UserRepositoryMock)
	passwordHasher := new(PasswordHasherMock)
	uow := new(m.UnitOfWorkMock)

	sut := usecase.NewResetPasswordUseCase(passwordResetRepo, userRepo, passwordHasher, uow)
	return sut, passwordResetRepo, userRepo, passwordHasher, uow
}

func TestResetPasswordUseCase(t *testing.T) {
	input := port.ResetPasswordInput{
		Token:       "plain-token-from-user",
		NewPassword: "newSecurePassword123",
	}

	t.Run("Execute", func(t *testing.T) {
		t.Run("should reset password successfully", func(t *testing.T) {
			sut, passwordResetRepo, userRepo, passwordHasher, uow := setupResetPasswordSut()

			userID := uuid.New().String()
			expiresAt := time.Now().Add(10 * time.Minute)
			validToken := &domain.PasswordResetToken{
				Token:     "hashed_token_value",
				UserID:    userID,
				ExpiresAt: expiresAt,
			}

			hashedPassword := "hashed_new_password"

			passwordResetRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(validToken, nil).Once()
			passwordHasher.On("Hash", input.NewPassword).Return(hashedPassword, nil).Once()
			userRepo.On("UpdatePassword", mock.Anything, userID, hashedPassword).Return(nil).Once()
			passwordResetRepo.On("MarkTokenAsUsed", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			uow.On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil).Once()

			err := sut.Execute(context.Background(), input)

			require.NoError(t, err)
			passwordResetRepo.AssertExpectations(t)
			passwordHasher.AssertExpectations(t)
			uow.AssertExpectations(t)
		})

		t.Run("should return NotFoundError when token does not exist", func(t *testing.T) {
			sut, passwordResetRepo, _, _, _ := setupResetPasswordSut()

			passwordResetRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Once()

			err := sut.Execute(context.Background(), input)

			require.Error(t, err)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "TOKEN_NOT_FOUND", appErr.Code)
		})

		t.Run("should return TOKEN_USED error if token already used", func(t *testing.T) {
			sut, passwordResetRepo, _, _, _ := setupResetPasswordSut()

			usedAt := time.Now().Add(-1 * time.Minute)
			usedToken := &domain.PasswordResetToken{
				Token:     "hashed_token_value",
				UserID:    uuid.New().String(),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				UsedAt:    &usedAt,
			}

			passwordResetRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(usedToken, nil).Once()

			err := sut.Execute(context.Background(), input)

			require.Error(t, err)
			require.ErrorIs(t, err, domain.TokenAlreadyUsedError)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "TOKEN_USED", appErr.Code)
		})

		t.Run("should return TOKEN_EXPIRED error if token is expired", func(t *testing.T) {
			sut, passwordResetRepo, _, _, _ := setupResetPasswordSut()

			expiredToken := &domain.PasswordResetToken{
				Token:     "hashed_token_value",
				UserID:    uuid.New().String(),
				ExpiresAt: time.Now().Add(-10 * time.Minute),
			}

			passwordResetRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(expiredToken, nil).Once()

			err := sut.Execute(context.Background(), input)

			require.Error(t, err)
			require.ErrorIs(t, err, domain.TokenExpiredError)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, "TOKEN_EXPIRED", appErr.Code)
		})

		t.Run("should bubble up error when uow fails", func(t *testing.T) {
			sut, passwordResetRepo, _, passwordHasher, uow := setupResetPasswordSut()

			validToken := &domain.PasswordResetToken{
				Token:     "hashed_token_value",
				UserID:    uuid.New().String(),
				ExpiresAt: time.Now().Add(10 * time.Minute),
			}
			hashedPassword := "hashed_new_password"
			expectedErr := errors.New("db error in uow")

			passwordResetRepo.On("GetByToken", mock.Anything, mock.AnythingOfType("string")).Return(validToken, nil).Once()
			passwordHasher.On("Hash", input.NewPassword).Return(hashedPassword, nil).Once()
			uow.On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(expectedErr).Once()

			err := sut.Execute(context.Background(), input)

			require.ErrorIs(t, err, expectedErr)
		})
	})
}
