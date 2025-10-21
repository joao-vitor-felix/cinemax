package port

import "github.com/joao-vitor-felix/cinemax/internal/core/domain"

type UseRepository interface {
	Create(domain.User) (*domain.User, error)
	IsContactInfoAvailable(email, phone string) (bool, error)
}
