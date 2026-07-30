package port

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	IsContactInfoAvailable(ctx context.Context, email, phone string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	CreateWithIdentity(ctx context.Context, user *domain.User, identity *domain.FederatedIdentity) error
}
