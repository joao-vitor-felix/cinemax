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

type SignInService struct {
	mock.Mock
}

func (s *SignInService) Execute(input port.SignInInput) (*port.SignInOutput, error) {
	args := s.Called(input)
	return args.Get(0).(*port.SignInOutput), args.Error(1)
}

func TestSignInController(t *testing.T) {
	setupSut := func() (*controller.SignInController, *SignInService) {
		serviceMock := new(SignInService)
		sut := controller.NewSignInController(serviceMock)
		return sut, serviceMock
	}

	url := "/auth/sign-in"

	t.Run("Execute", func(t *testing.T) {
		t.Run("should sign in and return 200 alongside the tokens", func(t *testing.T) {
			sut, service := setupSut()

			input := port.SignInInput{
				Email:    "john@john.com",
				Password: "12345678",
			}
			output := &port.SignInOutput{
				AccessToken:  "access_token_123",
				RefreshToken: "refresh_token_123",
			}

			service.On("Execute", input).Return(output, nil).Once()

			//TODO: make a helper function to make requests in controller tests
			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			res, err := sut.Execute(w, r)

			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.Status)
			require.Equal(t, output, res.Data.(*controller.Resource[*port.SignInOutput]).Data)

			service.AssertExpectations(t)
		})

		t.Run("should return an error when service fails", func(t *testing.T) {
			sut, service := setupSut()

			input := port.SignInInput{
				Email:    "john@john.com",
				Password: "12345678",
			}

			service.On("Execute", input).Return(&port.SignInOutput{}, domain.InvalidCredentialsError).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			_, err := sut.Execute(w, r)

			require.Error(t, err)
			require.Equal(t, domain.InvalidCredentialsError, err)

			service.AssertExpectations(t)
		})

		validationTests :=
			[]struct {
				name          string
				expectedError string
				body          port.SignInInput
			}{
				{
					name:          "email is missing",
					expectedError: "(?i)required",
					body: port.SignInInput{
						Email:    "",
						Password: "12345678",
					},
				},
				{
					name:          "password is missing",
					expectedError: "(?i)required",
					body: port.SignInInput{
						Email:    "john@john.com",
						Password: "",
					},
				},
				{
					name:          "invalid is email",
					expectedError: "(?i)valid email",
					body: port.SignInInput{
						Email:    "john",
						Password: "12345678",
					},
				},
				{
					name:          "password is too short",
					expectedError: "(?i)minimum length",
					body: port.SignInInput{
						Email:    "john@john.com",
						Password: "123",
					},
				},
				{
					name:          "password is too long",
					expectedError: "(?i)maximum length",
					body: port.SignInInput{
						Email:    "john@john.com",
						Password: "1234567890123",
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
