package port

import "github.com/joao-vitor-felix/cinemax/internal/core/domain"

type UserRepository interface {
	Create(user *domain.User) (*domain.User, error)
	IsContactInfoAvailable(email, phone string) (bool, error)
}

type UserService interface {
	Register(user *domain.User) (*domain.User, error)
}
