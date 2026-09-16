package repository

import (
	"context"
	"database/sql"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type PostgresPasswordResetTokenRepository struct {
	db *sql.DB
}

func NewPostgresPasswordResetTokenRepository(db *sql.DB) *PostgresPasswordResetTokenRepository {
	return &PostgresPasswordResetTokenRepository{db}
}

func (r *PostgresPasswordResetTokenRepository) GetByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	query := `
		SELECT
			token,
			user_id,
			expires_at,
			used_at,
			created_at
		FROM password_reset_tokens
		WHERE token = $1
	`

	var t domain.PasswordResetToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&t.Token,
		&t.UserID,
		&t.ExpiresAt,
		&t.UsedAt,
		&t.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresPasswordResetTokenRepository) GenerateToken(ctx context.Context, tokenHash, userID string) (*domain.PasswordResetToken, error) {
	query := `
		INSERT INTO password_reset_tokens (token, user_id, expires_at)
		VALUES ($1, $2, NOW() + INTERVAL '10 minutes')
		RETURNING token, user_id, expires_at, used_at
	`

	var t domain.PasswordResetToken
	err := r.db.QueryRowContext(ctx, query, tokenHash, userID).Scan(
		&t.Token,
		&t.UserID,
		&t.ExpiresAt,
		&t.UsedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresPasswordResetTokenRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE token = $1
	`
	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
