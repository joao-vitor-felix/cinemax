package usecase

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type SignUpUseCase struct {
	userRepo       port.UserRepository
	passwordHasher port.PasswordHasher
}

func NewSignUpUseCase(userRepo port.UserRepository, passwordHasher port.PasswordHasher) port.SignUpUseCase {
	return &SignUpUseCase{
		userRepo,
		passwordHasher,
	}
}

func (s *SignUpUseCase) Execute(ctx context.Context, input port.SignUpInput) (*domain.User, error) {
	user, err := domain.NewUser(
		input.FirstName,
		input.LastName,
		input.Email,
		input.Phone,
		input.DateOfBirth,
		input.Gender,
	)
	if err != nil {
		return nil, err
	}
	isAvailable, err := s.userRepo.IsContactInfoAvailable(ctx, user.Email, *user.Phone)
	if err != nil {
		return nil, err
	}
	if !isAvailable {
		return nil, domain.ContactInfoUnavailableError
	}
	passwordHash, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = &passwordHash
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}
