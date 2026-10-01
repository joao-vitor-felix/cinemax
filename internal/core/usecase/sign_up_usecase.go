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

func NewSignUpUseCase(userRepo port.UserRepository, passwordHasher port.PasswordHasher) *SignUpUseCase {
	return &SignUpUseCase{
		userRepo,
		passwordHasher,
	}
}

func (u *SignUpUseCase) Execute(ctx context.Context, input port.SignUpInput) (*domain.User, error) {
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
	isAvailable, err := u.userRepo.IsContactInfoAvailable(ctx, user.Email, *user.Phone)
	if err != nil {
		return nil, err
	}
	if !isAvailable {
		return nil, domain.ContactInfoUnavailableError
	}
	passwordHash, err := u.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = &passwordHash
	createdUser, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}
