package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/joao-vitor-felix/cinemax/internal/core/service"
	m "github.com/joao-vitor-felix/cinemax/test/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type HasherMock struct {
	mock.Mock
}

func (h *HasherMock) Hash(password string) (string, error) {
	args := h.Called(password)
	return args.String(0), args.Error(1)
}

func (h *HasherMock) Compare(hash, password string) error {
	args := h.Called(hash, password)
	return args.Error(0)
}

func setupSut() (*service.UserService, *m.UserRepositoryMock, *HasherMock) {
	repoMock := new(m.UserRepositoryMock)
	hasherMock := new(HasherMock)
	userService := service.NewUserService(repoMock, hasherMock)
	return userService, repoMock, hasherMock
}

func TestUserService(t *testing.T) {
	input := port.RegisterUserInput{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
		Phone:       "+12125551234",
		Password:    "password123",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
	}

	t.Run("Register", func(t *testing.T) {
		t.Parallel()
		t.Run("should register a new user successfully", func(t *testing.T) {
			userService, repoMock, hasherMock := setupSut()

			expectedHash := "hashed_password_123"
			repoMock.On("IsContactInfoAvailable", input.Email, input.Phone).Return(true, nil)
			hasherMock.On("Hash", input.Password).Return(expectedHash, nil)
			repoMock.On("Create", mock.AnythingOfType("*domain.User")).Return(&domain.User{
				ID:              "",
				FirstName:       input.FirstName,
				LastName:        input.LastName,
				Email:           input.Email,
				Phone:           input.Phone,
				PasswordHash:    expectedHash,
				DateOfBirth:     input.DateOfBirth,
				Gender:          input.Gender,
				ProfilePhotoURL: nil,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}, nil).Once()

			user, err := userService.Register(input)

			require.NoError(t, err)
			require.Equal(t, input.FirstName, user.FirstName)
			require.Equal(t, input.LastName, user.LastName)
			require.Equal(t, input.Email, user.Email)
			require.Equal(t, input.Phone, user.Phone)
			require.Equal(t, input.DateOfBirth, user.DateOfBirth)
			require.Equal(t, input.Gender, user.Gender)
			require.Equal(t, expectedHash, user.PasswordHash)

			repoMock.AssertExpectations(t)
			hasherMock.AssertExpectations(t)
		})

		t.Run("should return AppError when user validation fails", func(t *testing.T) {
			userService, _, _ := setupSut()
			invalidInput := input
			invalidInput.Gender = "invalid_gender"

			_, err := userService.Register(invalidInput)

			require.Error(t, err)
			var appErr *domain.AppError
			require.ErrorAs(t, err, &appErr, "error should be of type *domain.AppError")
		})

		t.Run("should return ContactInfoUnavailableError when contact data is already taken", func(t *testing.T) {
			userService, repoMock, _ := setupSut()
			repoMock.On("IsContactInfoAvailable", input.Email, input.Phone).Return(false, nil).Once()
			_, err := userService.Register(input)

			require.Error(t, err)
			appErr, ok := err.(*domain.AppError)
			require.True(t, ok, "error should be *domain.AppError")
			require.Equal(t, domain.ContactInfoUnavailableError.Code, appErr.Code)
			require.Equal(t, domain.ContactInfoUnavailableError.Message, appErr.Message)
			require.Equal(t, domain.ContactInfoUnavailableError.StatusCode, appErr.StatusCode)

			repoMock.AssertExpectations(t)
		})

		t.Run("should return error when IsContactInfoAvailable fails", func(t *testing.T) {
			userService, repoMock, _ := setupSut()
			expectedErr := errors.New("database connection failed")
			repoMock.On("IsContactInfoAvailable", input.Email, input.Phone).Return(true, expectedErr).Once()

			_, err := userService.Register(input)

			require.Error(t, err)
			require.Equal(t, expectedErr, err)

			repoMock.AssertExpectations(t)
		})

		t.Run("should return error when password hashing fails", func(t *testing.T) {
			userService, mockRepo, hasherMock := setupSut()
			expectedErr := errors.New("hashing failed")
			mockRepo.On("IsContactInfoAvailable", input.Email, input.Phone).Return(true, nil).Once()
			hasherMock.On("Hash", input.Password).Return("", expectedErr).Once()

			_, err := userService.Register(input)

			require.Error(t, err)
			require.Equal(t, expectedErr, err)
			mockRepo.AssertExpectations(t)
			hasherMock.AssertExpectations(t)
		})

		t.Run("should return error when Create fails", func(t *testing.T) {
			userService, mockRepo, hasherMock := setupSut()
			expectedErr := errors.New("database insert failed")
			expectedHash := "hashed_password_123"

			mockRepo.On("IsContactInfoAvailable", input.Email, input.Phone).Return(true, nil).Once()
			hasherMock.On("Hash", input.Password).Return(expectedHash, nil).Once()
			mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil, expectedErr).Once()

			_, err := userService.Register(input)

			require.Error(t, err)
			require.Equal(t, expectedErr, err)
			mockRepo.AssertExpectations(t)
			hasherMock.AssertExpectations(t)
		})
	})
}
