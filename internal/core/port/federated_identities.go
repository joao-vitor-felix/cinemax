package port

import "github.com/joao-vitor-felix/cinemax/internal/core/domain"

// TODO: a delete method will be necessary down the road for account unlinking
type FederatedIdentitiesRepository interface {
	GetByProviderUserID(provider, providerUserID string) (*domain.FederatedIdentity, error)
	CreateFederatedIdentity(userID, provider, providerUserID string) (*domain.FederatedIdentity, error)
}
