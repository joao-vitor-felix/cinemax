package unit

import (
	"testing"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/stretchr/testify/require"
)

func TestUserEntity(t *testing.T) {
	t.Parallel()

	t.Run("should return InvalidGenderError when gender is invalid", func(t *testing.T) {
		_, err := domain.NewUser(domain.User{
			Gender: "invalid",
		})
		require.Error(t, err)
		require.Equal(t, domain.InvalidGenderError, err)
	})

	t.Run("should return UserTooYoungError when user is under 13 years old", func(t *testing.T) {
		_, err := domain.NewUser(domain.User{
			Gender:      "male",
			DateOfBirth: "2015-01-01",
		})

		require.Error(t, err)
		require.Equal(t, domain.UserTooYoungError, err)
	})

	t.Run("should return a user successfully", func(t *testing.T) {
		input := domain.User{
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			Phone:       "+12125551234",
			Gender:      domain.Male,
			DateOfBirth: "1990-01-01",
		}

		user, err := domain.NewUser(input)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, input.FirstName, user.FirstName)
		require.Equal(t, input.LastName, user.LastName)
		require.Equal(t, input.Email, user.Email)
		require.Equal(t, input.Phone, user.Phone)
		require.Equal(t, input.Gender, user.Gender)
		require.Equal(t, input.DateOfBirth, user.DateOfBirth)
	})
}
