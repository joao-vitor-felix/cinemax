package port

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
)

type UserRepository interface {
	Create(user *domain.User) (*domain.User, error)
	IsContactInfoAvailable(email, phone string) (bool, error)
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
}
