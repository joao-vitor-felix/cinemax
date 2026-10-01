package usecase

import (
	"context"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
	"github.com/joao-vitor-felix/cinemax/internal/core/util"
)

type ForgotPasswordUseCase struct {
	passwordResetRepo port.PasswordResetRepository
	userRepo          port.UserRepository
}

func NewPasswordResetUseCase(
	passwordResetRepo port.PasswordResetRepository,
	userRepo port.UserRepository,
) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{
		passwordResetRepo: passwordResetRepo,
		userRepo:          userRepo,
	}
}

func (u *ForgotPasswordUseCase) Execute(ctx context.Context, input port.ForgotPasswordInput) error {
	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	if user == nil {
		return domain.NotFoundError("user")
	}

	token, err := util.GenerateSecureTokenCustom(64)

	if err != nil {
		return err
	}

	_, err = u.passwordResetRepo.SaveToken(ctx, token, user.ID.String())

	if err != nil {
		return err
	}

	//TODO: send e-mail
	return nil
}
