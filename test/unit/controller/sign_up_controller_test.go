package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joao-vitor-felix/cinemax/internal/adapter/controller"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/joao-vitor-felix/cinemax/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type SignUpService struct {
	mock.Mock
}

func (s *SignUpService) Execute(input port.SignUpInput) (*domain.User, error) {
	args := s.Called(input)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestSignUpController(t *testing.T) {
	setupSut := func() (*controller.SignUpController, *SignUpService) {
		serviceMock := new(SignUpService)
		sut := controller.NewSignUpController(serviceMock)
		return sut, serviceMock
	}

	url := "/auth/sign-up"

	t.Run("Execute", func(t *testing.T) {
		t.Run("should return 201 and register a new user successfully", func(t *testing.T) {
			sut, service := setupSut()
			input := port.SignUpInput{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Phone:       "+12125551234",
				Password:    "password123",
				DateOfBirth: "1990-01-01",
				Gender:      "male",
			}

			service.On("Execute", input).Return(&domain.User{}, nil).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			res, err := sut.Execute(w, r)

			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, res.Status)
			require.Equal(t, nil, res.Data.(*controller.Resource[any]).Data)
			service.AssertExpectations(t)
		})

		t.Run("should return an error when service fails", func(t *testing.T) {
			sut, service := setupSut()
			input := port.SignUpInput{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Phone:       "+12125551234",
				Password:    "password123",
				DateOfBirth: "1990-01-01",
				Gender:      "male",
			}

			service.On("Execute", input).Return(&domain.User{}, domain.ContactInfoUnavailableError).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			_, err := sut.Execute(w, r)

			require.Error(t, err)
			require.Equal(t, domain.ContactInfoUnavailableError, err)
			service.AssertExpectations(t)
		})

		validationTests := []struct {
			name          string
			expectedError string
			body          port.SignUpInput
		}{
			{
				name:          "missing first name",
				expectedError: "(?i)required",
				body: port.SignUpInput{
					FirstName:   "",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "first name too short",
				expectedError: "(?i)minimum length",
				body: port.SignUpInput{
					FirstName:   "J",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "missing last name",
				expectedError: "(?i)required",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "last name too short",
				expectedError: "(?i)minimum length",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "D",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "missing email",
				expectedError: "(?i)required",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "invalid email format",
				expectedError: "(?i)valid email",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "invalid-email",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "missing phone",
				expectedError: "(?i)required",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "invalid phone format (not e164)",
				expectedError: "(?i)must be a valid",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "12125551234",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "missing password",
				expectedError: "(?i)required",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "password too short",
				expectedError: "(?i)minimum length",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "pass",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "password too long",
				expectedError: "(?i)maximum length",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "passwordtoolong123",
					DateOfBirth: "1990-01-01",
					Gender:      "male",
				},
			},
			{
				name:          "missing date of birth",
				expectedError: "(?i)required",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "",
					Gender:      "male",
				},
			},
			{
				name:          "invalid date of birth format",
				expectedError: "(?i)must be a valid date",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "01-01-1990",
					Gender:      "male",
				},
			},
			{
				name:          "invalid gender",
				expectedError: "(?i)contains an invalid value",
				body: port.SignUpInput{
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john.doe@example.com",
					Phone:       "+12125551234",
					Password:    "password123",
					DateOfBirth: "1990-01-01",
					Gender:      "invalid",
				},
			},
		}

		for _, tt := range validationTests {
			t.Run("should throw when "+tt.name, func(t *testing.T) {
				sut, _ := setupSut()

				r := test.MakeRequest(http.MethodPost, url, tt.body)
				w := httptest.NewRecorder()

				_, err := sut.Execute(w, r)

				require.Error(t, err)
				var appErr *domain.AppError
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, http.StatusBadRequest, appErr.StatusCode)
				require.Equal(t, "VALIDATION_ERROR", appErr.Code)
				require.Regexp(t, tt.expectedError, appErr.Message)
			})
		}
	})
}
