package unit

import (
	"testing"
	"time"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/stretchr/testify/require"
)

func TestPasswordResetTokenEntity(t *testing.T) {
	t.Run("should return TokenAlreadyUsedError when token has been used", func(t *testing.T) {
		usedAt := time.Now()
		token := &domain.PasswordResetToken{
			Token:     "token123",
			UserID:    "user123",
			ExpiresAt: time.Now().Add(10 * time.Minute),
			UsedAt:    &usedAt,
		}

		require.True(t, token.IsUsed())
		err := token.Validate()
		require.ErrorIs(t, err, domain.TokenAlreadyUsedError)
	})

	t.Run("should return TokenExpiredError when token is expired", func(t *testing.T) {
		token := &domain.PasswordResetToken{
			Token:     "token123",
			UserID:    "user123",
			ExpiresAt: time.Now().Add(-10 * time.Minute),
		}

		require.True(t, token.IsExpired())
		err := token.Validate()
		require.ErrorIs(t, err, domain.TokenExpiredError)
	})

	t.Run("should return nil when token is valid", func(t *testing.T) {
		token := &domain.PasswordResetToken{
			Token:     "token123",
			UserID:    "user123",
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}

		require.False(t, token.IsUsed())
		require.False(t, token.IsExpired())
		err := token.Validate()
		require.NoError(t, err)
	})
}
