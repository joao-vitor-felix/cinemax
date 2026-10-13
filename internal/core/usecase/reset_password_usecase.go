package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/joao-vitor-felix/cinemax/internal/core/domain"
	"github.com/joao-vitor-felix/cinemax/internal/core/port"
)

type ResetPasswordUseCase struct {
	passwordResetRepo port.PasswordResetRepository
	userRepo          port.UserRepository
	passwordHasher    port.PasswordHasher
	uow               port.UnitOfWork
}

func NewResetPasswordUseCase(
	passwordResetRepo port.PasswordResetRepository,
	userRepo port.UserRepository,
	passwordHasher port.PasswordHasher,
	uow port.UnitOfWork,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		passwordResetRepo: passwordResetRepo,
		userRepo:          userRepo,
		passwordHasher:    passwordHasher,
		uow:               uow,
	}
}

func (u *ResetPasswordUseCase) Execute(ctx context.Context, input port.ResetPasswordInput) error {
	hasher := sha256.New()
	hasher.Write([]byte(input.Token))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	resetToken, err := u.passwordResetRepo.GetByToken(ctx, hashedToken)
	if err != nil {
		return err
	}
	if resetToken == nil {
		return domain.NotFoundError("token")
	}

	if err := resetToken.Validate(); err != nil {
		return err
	}

	hashedPassword, err := u.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return err
	}

	err = u.uow.Do(ctx, func(txCtx context.Context) error {
		if err := u.userRepo.UpdatePassword(txCtx, resetToken.UserID, hashedPassword); err != nil {
			return err
		}

		if err := u.passwordResetRepo.MarkTokenAsUsed(txCtx, hashedToken); err != nil {
			return err
		}

		// TODO: invalidate active sessions

		return nil
	})

	if err != nil {
		return err
	}

	//TODO: email user

	return nil
}
