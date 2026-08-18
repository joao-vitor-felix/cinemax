package repository

import (
	"context"
	"database/sql"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			date_of_birth,
			gender,
			profile_photo_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.DateOfBirth,
		user.Gender,
		user.ProfilePhotoURL,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) IsContactInfoAvailable(ctx context.Context, email, phone string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1 OR phone = $2
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email, phone).Scan(&exists)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			date_of_birth,
			gender,
			profile_photo_url,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.DateOfBirth,
		&user.Gender,
		&user.ProfilePhotoURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			date_of_birth,
			gender,
			profile_photo_url,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.DateOfBirth,
		&user.Gender,
		&user.ProfilePhotoURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
