package repository

import (
	"database/sql"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db}
}

func (r *PostgresUserRepository) Create(user *domain.User) (*domain.User, error) {
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
		RETURNING id
	`
	err := r.db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Phone, user.PasswordHash, user.DateOfBirth, user.Gender, user.ProfilePhotoURL).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) IsContactInfoAvailable(email, phone string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1 OR phone = $2
		)
	`
	var exists bool
	err := r.db.QueryRow(query, email, phone).Scan(&exists)
	if err != nil {
		return false, err
	}
	return !exists, nil
}
