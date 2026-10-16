package controller_test

import (
	"context"
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

type ForgotPasswordUseCaseMock struct {
	mock.Mock
}

func (m *ForgotPasswordUseCaseMock) Execute(ctx context.Context, input port.ForgotPasswordInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func TestForgotPasswordController(t *testing.T) {
	setupSut := func() (*controller.ForgotPasswordController, *ForgotPasswordUseCaseMock) {
		useCaseMock := new(ForgotPasswordUseCaseMock)
		sut := controller.NewForgotPasswordController(useCaseMock)
		return sut, useCaseMock
	}

	url := "/auth/forgot-password"

	t.Run("Execute", func(t *testing.T) {
		t.Run("should return 200 successfully", func(t *testing.T) {
			sut, useCase := setupSut()
			input := port.ForgotPasswordInput{
				Email: "john.doe@example.com",
			}

			useCase.On("Execute", mock.Anything, input).Return(nil).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			sut.Execute(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			var body controller.MessageResponse
			require.NoError(t, controller.DecodeJSON(w.Body, &body))
			require.Equal(t, "if you have an e-mail registered, a code was sent to you.", body.Message)
			useCase.AssertExpectations(t)
		})

		t.Run("should return an error when user is not found", func(t *testing.T) {
			sut, useCase := setupSut()
			input := port.ForgotPasswordInput{
				Email: "john.doe@example.com",
			}

			notFoundError := domain.NotFoundError("user")
			useCase.On("Execute", mock.Anything, input).Return(notFoundError).Once()

			r := test.MakeRequest(http.MethodPost, url, input)
			w := httptest.NewRecorder()

			sut.Execute(w, r)

			require.Equal(t, http.StatusNotFound, w.Code)
			var body controller.ErrorResponse
			require.NoError(t, controller.DecodeJSON(w.Body, &body))
			require.Equal(t, notFoundError.Code, body.Code)
			require.Equal(t, notFoundError.Message, body.Message)
			useCase.AssertExpectations(t)
		})

		validationTests := []struct {
			name          string
			expectedError string
			body          port.ForgotPasswordInput
		}{
			{
				name:          "missing email",
				expectedError: "(?i)required",
				body: port.ForgotPasswordInput{
					Email: "",
				},
			},
			{
				name:          "invalid email format",
				expectedError: "(?i)valid email",
				body: port.ForgotPasswordInput{
					Email: "invalid-email",
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
