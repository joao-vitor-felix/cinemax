package factory

import (
	"database/sql"
	"os"

	"github.com/joao-vitor-felix/cinemax/internal/adapter/auth"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/controller"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/repository"
	"github.com/joao-vitor-felix/cinemax/internal/core/service"
)

type Container struct {
	SignUpController           *controller.SignUpController
	SignInController           *controller.SignInController
	RefreshTokenController     *controller.RefreshTokenController
	SignInWithGoogleController *controller.SignInGoogleController
}

func NewContainer(db *sql.DB) *Container {
	passwordHasherAdapter := auth.NewPasswordHasherAdapter()
	tokenIssuerAdapter := auth.NewTokenIssuerAdapter()
	googleOAuthAdapter := auth.NewGoogleOAuthAdapter(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URI"),
	)

	userRepo := repository.NewPostgresUserRepository(db)
	refreshTokenRepo := repository.NewPostgresRefreshTokenRepository(db)

	signUpService := service.NewSignUpService(userRepo, passwordHasherAdapter)
	signInService := service.NewSignInService(userRepo, passwordHasherAdapter, tokenIssuerAdapter, refreshTokenRepo)
	refreshTokenService := service.NewRefreshTokenService(refreshTokenRepo, userRepo, tokenIssuerAdapter)
	signInGoogleService := service.NewSignInGoogleService(googleOAuthAdapter, userRepo, tokenIssuerAdapter, refreshTokenRepo)

	return &Container{
		SignUpController:           controller.NewSignUpController(signUpService),
		SignInController:           controller.NewSignInController(signInService),
		RefreshTokenController:     controller.NewRefreshTokenController(refreshTokenService),
		SignInWithGoogleController: controller.NewSignInGoogleController(signInGoogleService),
	}
}
