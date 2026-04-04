package unit

import (
	"testing"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/stretchr/testify/require"
)

func TestUserEntity(t *testing.T) {
	t.Run("should return InvalidGenderError when gender is invalid", func(t *testing.T) {
		_, err := domain.NewUser("John", "Doe", "john.doe@example.com", "+12125551234", "1990-01-01", "invalid")
		require.Error(t, err)
		require.Equal(t, domain.InvalidGenderError, err)
	})

	t.Run("should return UserTooYoungError when user is under 13 years old", func(t *testing.T) {
		_, err := domain.NewUser("John", "Doe", "john.doe@example.com", "+12125551234", "2015-01-01", "male")

		require.Error(t, err)
		require.Equal(t, domain.UserTooYoungError, err)
	})

	t.Run("should return a user successfully", func(t *testing.T) {
		firstName := "John"
		lastName := "Doe"
		email := "john.doe@example.com"
		phone := "+12125551234"
		gender := domain.Male
		dateOfBirth := "1990-01-01"

		user, err := domain.NewUser(firstName, lastName, email, phone, dateOfBirth, gender)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, firstName, user.FirstName)
		require.Equal(t, lastName, user.LastName)
		require.Equal(t, email, user.Email)
		require.Equal(t, phone, user.Phone)
		require.Equal(t, gender, user.Gender)
		require.Equal(t, dateOfBirth, user.DateOfBirth)
	})
}
