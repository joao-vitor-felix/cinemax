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

type RefreshTokenServiceMock struct {
	mock.Mock
}

func (s *RefreshTokenServiceMock) Execute(input port.RefreshTokenInput) (*port.RefreshTokenOutput, error) {
	args := s.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.RefreshTokenOutput), args.Error(1)
}

func TestRefreshTokenController(t *testing.T) {
	setupSut := func() (*controller.RefreshTokenController, *RefreshTokenServiceMock) {
		serviceMock := new(RefreshTokenServiceMock)
		sut := controller.NewRefreshTokenController(serviceMock)
		return sut, serviceMock
	}

	url := "/auth/refresh-token"

	t.Run("Execute", func(t *testing.T) {
		t.Run("should refresh tokens successfully and return 200", func(t *testing.T) {
			sut, service := setupSut()

			input := port.RefreshTokenInput{
				RefreshToken: "valid_refresh_token",
			}
			output := &port.RefreshTokenOutput{
				AccessToken:  "new_access_token",
				RefreshToken: "new_refresh_token",
			}

			service.On("Execute", input).Return(output, nil).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			res, err := sut.Execute(w, r)

			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.Status)
			require.Equal(t, output, res.Data.(*port.RefreshTokenOutput))

			service.AssertExpectations(t)
		})

		t.Run("should return an error when service fails", func(t *testing.T) {
			sut, service := setupSut()

			input := port.RefreshTokenInput{
				RefreshToken: "invalid_refresh_token",
			}

			service.On("Execute", input).Return(nil, domain.InvalidCredentialsError).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			_, err := sut.Execute(w, r)

			require.Error(t, err)
			require.Equal(t, domain.InvalidCredentialsError, err)

			service.AssertExpectations(t)
		})

		validationTests := []struct {
			name          string
			expectedError string
			body          port.RefreshTokenInput
		}{
			{
				name:          "refresh token is missing",
				expectedError: "(?i)required",
				body: port.RefreshTokenInput{
					RefreshToken: "",
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
