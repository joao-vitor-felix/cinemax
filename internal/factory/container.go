package factory

import (
	"database/sql"
	"os"

	"github.com/joao-vitor-felix/cinemax/internal/adapter/auth"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/controller"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/repository"
	"github.com/joao-vitor-felix/cinemax/internal/core/usecase"
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

	signUpUseCase := usecase.NewSignUpUseCase(userRepo, passwordHasherAdapter)
	signInUseCase := usecase.NewSignInUseCase(userRepo, passwordHasherAdapter, tokenIssuerAdapter, refreshTokenRepo)
	refreshTokenUseCase := usecase.NewRefreshTokenUseCase(refreshTokenRepo, userRepo, tokenIssuerAdapter)
	signInGoogleUseCase := usecase.NewSignInGoogleUseCase(googleOAuthAdapter, userRepo, tokenIssuerAdapter, refreshTokenRepo)

	return &Container{
		SignUpController:           controller.NewSignUpController(signUpUseCase),
		SignInController:           controller.NewSignInController(signInUseCase),
		RefreshTokenController:     controller.NewRefreshTokenController(refreshTokenUseCase),
		SignInWithGoogleController: controller.NewSignInGoogleController(signInGoogleUseCase),
	}
}
