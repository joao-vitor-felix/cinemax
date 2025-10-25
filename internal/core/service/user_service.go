package service

import (
	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type UserService struct {
	repo           port.UserRepository
	passwordHasher port.PasswordHasher
}

func NewUserService(repo port.UserRepository, passwordHasher port.PasswordHasher) *UserService {
	return &UserService{
		repo,
		passwordHasher,
	}
}

func (s *UserService) Register(input port.RegisterUserInput) (*domain.User, error) {
	user, err := domain.NewUser(domain.User{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Email:           input.Email,
		Phone:           input.Phone,
		DateOfBirth:     input.DateOfBirth,
		Gender:          input.Gender,
		ProfilePhotoURL: input.ProfilePhotoURL,
	})
	if err != nil {
		return nil, err
	}
	isAvailable, err := s.repo.IsContactInfoAvailable(user.Email, user.Phone)
	if err != nil {
		return nil, err
	}
	if !isAvailable {
		return nil, domain.ContactInfoUnavailableError
	}
	passwordHash, err := s.passwordHasher.Hash([]byte(input.Password))
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(passwordHash)
	s.repo.Create(user)
	//TODO: send email
	return user, nil
}
