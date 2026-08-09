package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

func (r *PostgresUserRepository) FindByProviderUserID(ctx context.Context, provider, providerUserID string) (*domain.User, error) {
	query := `
	SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.email,
		u.phone,
		u.password_hash,
		u.date_of_birth,
		u.gender,
		u.profile_photo_url,
		u.created_at,
		u.updated_at
	FROM
		users u
	INNER JOIN
		federated_identities fi ON fi.user_id = u.id
	WHERE
		fi.provider = $1
		AND fi.provider_user_id = $2
	`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, provider, providerUserID).Scan(
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

func (r *PostgresUserRepository) LinkIdentity(ctx context.Context, userID, provider, providerUserID string) (*domain.FederatedIdentity, error) {
	query := `
	INSERT INTO federated_identities (user_id, provider, provider_user_id)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, provider, provider_user_id, created_at
	`

	var fi domain.FederatedIdentity

	err := r.db.QueryRowContext(ctx, query, userID, provider, providerUserID).Scan(
		&fi.ID,
		&fi.UserID,
		&fi.Provider,
		&fi.ProviderUserID,
		&fi.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &fi, nil
}

func (r *PostgresUserRepository) CreateWithIdentity(ctx context.Context, firstName, lastName, email string, identity *domain.FederatedIdentity) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	userQuery := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var userID uuid.UUID
	err = tx.QueryRowContext(
		ctx,
		userQuery,
		firstName,
		lastName,
		email,
	).Scan(&userID)
	if err != nil {
		return uuid.Nil, err
	}

	identity.UserID = userID

	identityQuery := `
		INSERT INTO federated_identities (user_id, provider, provider_user_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var identityID uuid.UUID
	err = tx.QueryRowContext(
		ctx,
		identityQuery,
		identity.UserID,
		identity.Provider,
		identity.ProviderUserID,
	).Scan(&identityID)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
