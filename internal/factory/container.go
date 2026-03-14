package factory

import (
	"database/sql"

	"github.com/joao-vitor-felix/cinemax/internal/adapter/auth"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/controller"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/repository"
	"github.com/joao-vitor-felix/cinemax/internal/core/service"
)

type Container struct {
	UserController         *controller.UserController
	SignInController       *controller.SignInController
	RefreshTokenController *controller.RefreshTokenController
}

func NewContainer(db *sql.DB) *Container {
	// Adapters
	passwordHasher := auth.NewPasswordHasherAdapter()
	tokenIssuer := auth.NewTokenIssuerAdapter()

	// Repositories
	userRepo := repository.NewPostgresUserRepository(db)
	refreshTokenRepo := repository.NewPostgresRefreshTokenRepository(db)

	// Services
	userService := service.NewUserService(userRepo, passwordHasher)
	signInService := service.NewSignInService(userRepo, passwordHasher, tokenIssuer, refreshTokenRepo)
	refreshTokenService := service.NewRefreshTokenService(refreshTokenRepo, userRepo, tokenIssuer)

	return &Container{
		UserController:         controller.NewUserController(userService),
		SignInController:       controller.NewSignInController(signInService),
		RefreshTokenController: controller.NewRefreshTokenController(refreshTokenService),
	}
}
