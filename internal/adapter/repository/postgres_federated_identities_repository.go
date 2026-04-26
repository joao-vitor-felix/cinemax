package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type PostgresFederatedIdentitiesRepository struct {
	db *sql.DB
}

func NewPostgresFederatedIdentitiesRepository(db *sql.DB) *PostgresFederatedIdentitiesRepository {
	return &PostgresFederatedIdentitiesRepository{db}
}

func (r *PostgresFederatedIdentitiesRepository) GetByProviderUserID(provider, providerUserID string) (*domain.FederatedIdentity, error) {
	query := `
	SELECT
		id,
		user_id,
		provider,
		provider_user_id,
		created_at
	FROM
		federated_identities
	WHERE
		provider = $1
		AND provider_user_id = $2
	`

	var fi domain.FederatedIdentity

	err := r.db.QueryRow(query, provider, providerUserID).Scan(
		&fi.ID,
		&fi.UserID,
		&fi.Provider,
		&fi.ProviderUserID,
		&fi.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &fi, nil
}

func (r *PostgresFederatedIdentitiesRepository) CreateFederatedIdentity(userID, provider, providerUserID string) (*domain.FederatedIdentity, error) {
	query := `
	INSERT INTO federated_identities (user_id, provider, provider_user_id)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	var fi domain.FederatedIdentity
	userUUID, err := uuid.Parse(userID)

	if err != nil {
		return nil, err
	}

	fi.UserID = userUUID
	fi.Provider = provider
	fi.ProviderUserID = providerUserID

	err = r.db.QueryRow(query, userID, provider, providerUserID).Scan(
		&fi.ID,
		&fi.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &fi, nil
}
