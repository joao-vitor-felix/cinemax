package repository

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

const tokenExpiration = 7 * 24 * time.Hour

type PostgresRefreshTokenRepository struct {
	db *sql.DB
}

func NewPostgresRefreshTokenRepository(db *sql.DB) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db}
}

func (r *PostgresRefreshTokenRepository) GetByToken(token string) (*domain.RefreshToken, error) {
	query := `
	SELECT
		token,
		user_id,
		expires_at,
		used_at,
		created_at
	FROM refresh_tokens
	WHERE token = $1
	`

	var t domain.RefreshToken
	err := r.db.QueryRow(query, token).Scan(
		&t.Token,
		&t.UserId,
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

func (r *PostgresRefreshTokenRepository) GenerateToken(userId string) (*domain.RefreshToken, error) {
	tokenUUID := uuid.New().String()
	expiresAt := time.Now().Add(tokenExpiration)

	query := `
		INSERT INTO
			refresh_tokens (
			token,
			user_id,
			expires_at,
			created_at
		)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING
			token,
			user_id,
			expires_at,
			used_at,
			created_at
	`
	var t domain.RefreshToken
	err := r.db.QueryRow(query, tokenUUID, userId, expiresAt).Scan(
		&t.Token,
		&t.UserId,
		&t.ExpiresAt,
		&t.UsedAt,
		&t.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresRefreshTokenRepository) GenerateAndInvalidateUsedToken(token, userId string) (*domain.RefreshToken, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback transaction", "error", err)
		}
	}()

	updateQuery := `
		UPDATE
			refresh_tokens
		SET
			used_at = CURRENT_TIMESTAMP
		WHERE
			token = $1
		AND
			user_id = $2
	`
	_, err = tx.Exec(updateQuery, token, userId)
	if err != nil {
		return nil, err
	}

	newTokenUUID := uuid.New().String()
	expiresAt := time.Now().Add(tokenExpiration)

	insertQuery := `
		INSERT INTO
			refresh_tokens (
			token,
			user_id,
			expires_at,
			created_at
		)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING
			token,
			user_id,
			expires_at,
			used_at,
			created_at
	`
	var newToken domain.RefreshToken
	err = tx.QueryRow(insertQuery, newTokenUUID, userId, expiresAt).Scan(
		&newToken.Token,
		&newToken.UserId,
		&newToken.ExpiresAt,
		&newToken.UsedAt,
		&newToken.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &newToken, nil
}

func (r *PostgresRefreshTokenRepository) DeleteToken(token string) error {
	query := `
		DELETE FROM
			refresh_tokens
		WHERE
			token = $1
	`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *PostgresRefreshTokenRepository) DeleteTokensByUserID(userId string) error {
	query := `
		DELETE FROM
			refresh_tokens
		WHERE
			user_id = $1
	`
	_, err := r.db.Exec(query, userId)
	return err
}
