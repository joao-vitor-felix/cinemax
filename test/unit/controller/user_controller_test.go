package unit

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

type UserServiceMock struct {
	mock.Mock
}

func (u *UserServiceMock) Register(input port.RegisterUserInput) (*domain.User, error) {
	args := u.Called(input)
	return args.Get(0).(*domain.User), args.Error(1)
}

func setupSut() (*controller.UserController, *UserServiceMock) {
	serviceMock := new(UserServiceMock)
	sut := controller.NewUserController(serviceMock)
	return sut, serviceMock
}

func TestUserController(t *testing.T) {
	t.Parallel()

	url := "/auth/sign-up"

	t.Run("Register", func(t *testing.T) {
		t.Run("should register a new user successfully", func(t *testing.T) {
			sut, service := setupSut()
			input := port.RegisterUserInput{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Phone:       "+12125551234",
				Password:    "password123",
				DateOfBirth: "1990-01-01",
				Gender:      "male",
			}

			service.On("Register", input).Return(&domain.User{}, nil).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			resp, err := sut.Register(w, r)

			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, resp.Status)
			require.IsType(t, resp.Data, &controller.Resource{})
			require.Equal(t, nil, resp.Data.(*controller.Resource).Data)
			service.AssertExpectations(t)
		})

		t.Run("should throw InvalidBodyError when a invalid body is provided", func(t *testing.T) {
			sut, _ := setupSut()
			r := test.MakeRequest(http.MethodPost, url, "invalid json")
			w := httptest.NewRecorder()

			_, err := sut.Register(w, r)

			require.Error(t, err)
			require.Equal(t, domain.InvalidBodyError, err)
		})

		t.Run("should return an error when service fails", func(t *testing.T) {
			sut, service := setupSut()
			input := port.RegisterUserInput{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Phone:       "+12125551234",
				Password:    "password123",
				DateOfBirth: "1990-01-01",
				Gender:      "male",
			}

			service.On("Register", input).Return(&domain.User{}, domain.ContactInfoUnavailableError).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			_, err := sut.Register(w, r)

			require.Error(t, err)
			require.Equal(t, domain.ContactInfoUnavailableError, err)
			service.AssertExpectations(t)
		})

		validationTests := []struct {
			name string
			body port.RegisterUserInput
		}{
			{
				name: "missing first name",
				body: port.RegisterUserInput{
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
				name: "first name too short",
				body: port.RegisterUserInput{
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
				name: "missing last name",
				body: port.RegisterUserInput{
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
				name: "last name too short",
				body: port.RegisterUserInput{
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
				name: "missing email",
				body: port.RegisterUserInput{
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
				name: "invalid email format",
				body: port.RegisterUserInput{
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
				name: "missing phone",
				body: port.RegisterUserInput{
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
				name: "invalid phone format (not e164)",
				body: port.RegisterUserInput{
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
				name: "missing password",
				body: port.RegisterUserInput{
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
				name: "password too short",
				body: port.RegisterUserInput{
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
				name: "password too long",
				body: port.RegisterUserInput{
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
				name: "missing date of birth",
				body: port.RegisterUserInput{
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
				name: "invalid date of birth format",
				body: port.RegisterUserInput{
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
				name: "invalid gender",
				body: port.RegisterUserInput{
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

		for _, tc := range validationTests {
			t.Run(tc.name, func(t *testing.T) {
				sut, _ := setupSut()

				r := test.MakeRequest(http.MethodPost, url, tc.body)
				w := httptest.NewRecorder()

				_, err := sut.Register(w, r)

				require.Error(t, err)
				var appErr *domain.AppError
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, http.StatusBadRequest, appErr.StatusCode)
				require.Equal(t, "VALIDATION_ERROR", appErr.Code)
			})
		}
	})
}
