package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type PostgresPasswordResetRepository struct {
	db *sql.DB
}

func NewPostgresPasswordResetRepository(db *sql.DB) *PostgresPasswordResetRepository {
	return &PostgresPasswordResetRepository{db}
}

func (r *PostgresPasswordResetRepository) GetByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
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

func (r *PostgresPasswordResetRepository) SaveToken(ctx context.Context, tokenHash, userID string) (*domain.PasswordResetToken, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback transaction", "error", err)
		}
	}()

	updateQuery := `
		UPDATE password_reset_tokens
		SET expires_at = NOW()
		WHERE user_id = $1 AND used_at IS NULL AND expires_at > NOW()
	`
	_, err = tx.ExecContext(ctx, updateQuery, userID)
	if err != nil {
		return nil, err
	}

	insertQuery := `
		INSERT INTO password_reset_tokens (token, user_id, expires_at)
		VALUES ($1, $2, NOW() + INTERVAL '10 minutes')
		RETURNING token, user_id, expires_at, used_at
	`

	var t domain.PasswordResetToken
	err = tx.QueryRowContext(ctx, insertQuery, tokenHash, userID).Scan(
		&t.Token,
		&t.UserID,
		&t.ExpiresAt,
		&t.UsedAt,
	)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresPasswordResetRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
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
