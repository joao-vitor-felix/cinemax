package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	IsContactInfoAvailable(ctx context.Context, email, phone string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	CreateWithIdentity(ctx context.Context, firstName, lastName, email string, identity *domain.FederatedIdentity) (uuid.UUID, error)
	GetByProviderUserID(ctx context.Context, provider, providerUserID string) (*domain.User, error)
	LinkIdentity(ctx context.Context, userID, provider, providerUserID string) (*domain.FederatedIdentity, error)
}
