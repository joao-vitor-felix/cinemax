package controller_test

import (
	"context"
	"errors"
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

type ResetPasswordUseCaseMock struct {
	mock.Mock
}

func (m *ResetPasswordUseCaseMock) Execute(ctx context.Context, input port.ResetPasswordInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func TestResetPasswordController(t *testing.T) {
	setupSut := func() (*controller.ResetPasswordController, *ResetPasswordUseCaseMock) {
		useCaseMock := new(ResetPasswordUseCaseMock)
		sut := controller.NewResetPasswordController(useCaseMock)
		return sut, useCaseMock
	}

	url := "/auth/reset-password"

	t.Run("Execute", func(t *testing.T) {
		t.Run("should return 200 successfully", func(t *testing.T) {
			sut, useCase := setupSut()
			input := port.ResetPasswordInput{
				Token:       "some-valid-token",
				NewPassword: "newSecurePassword123",
			}

			useCase.On("Execute", mock.Anything, input).Return(nil).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			sut.Execute(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			var body controller.MessageResponse
			require.NoError(t, controller.DecodeJSON(w.Body, &body))
			require.Equal(t, "Password has been successfully reset.", body.Message)
			useCase.AssertExpectations(t)
		})

		t.Run("should return an error when token is invalid or expired", func(t *testing.T) {
			sut, useCase := setupSut()
			input := port.ResetPasswordInput{
				Token:       "invalid-token",
				NewPassword: "newSecurePassword123",
			}

			expectedErr := domain.TokenExpiredError
			useCase.On("Execute", mock.Anything, input).Return(expectedErr).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			sut.Execute(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			var body controller.ErrorResponse
			require.NoError(t, controller.DecodeJSON(w.Body, &body))
			require.Equal(t, expectedErr.Code, body.Code)
			require.Equal(t, expectedErr.Message, body.Message)
			useCase.AssertExpectations(t)
		})

		t.Run("should handle internal server error", func(t *testing.T) {
			sut, useCase := setupSut()
			input := port.ResetPasswordInput{
				Token:       "some-valid-token",
				NewPassword: "newSecurePassword123",
			}

			useCase.On("Execute", mock.Anything, input).Return(errors.New("db error")).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			sut.Execute(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
			useCase.AssertExpectations(t)
		})

		validationTests := []struct {
			name          string
			expectedError string
			body          port.ResetPasswordInput
		}{
			{
				name:          "missing token",
				expectedError: "(?i)required",
				body: port.ResetPasswordInput{
					Token:       "",
					NewPassword: "newSecurePassword123",
				},
			},
			{
				name:          "missing new password",
				expectedError: "(?i)required",
				body: port.ResetPasswordInput{
					Token:       "some-valid-token",
					NewPassword: "",
				},
			},
			{
				name:          "password too short",
				expectedError: "(?i)min",
				body: port.ResetPasswordInput{
					Token:       "some-valid-token",
					NewPassword: "short",
				},
			},
		}

		for _, tt := range validationTests {
			t.Run("should throw when "+tt.name, func(t *testing.T) {
				sut, _ := setupSut()

				r := test.MakeRequest(http.MethodPost, url, tt.body)
				w := httptest.NewRecorder()

				sut.Execute(w, r)

				require.Equal(t, http.StatusBadRequest, w.Code)
				var body controller.ErrorResponse
				require.NoError(t, controller.DecodeJSON(w.Body, &body))
				require.Equal(t, "VALIDATION_ERROR", body.Code)
				require.Regexp(t, tt.expectedError, body.Message)
			})
		}
	})
}
