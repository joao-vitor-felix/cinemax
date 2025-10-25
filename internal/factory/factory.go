package factory

import (
	"database/sql"

	"github.com/joao-vitor-felix/cinemax/internal/adapter/auth"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/controller"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/repository"
	"github.com/joao-vitor-felix/cinemax/internal/core/service"
)

func MakeNewUserController(db *sql.DB) *controller.UserController {
	repo := repository.NewPostgresUserRepository(db)
	passwordHasher := auth.NewPasswordHasher()
	userService := service.NewUserService(repo, passwordHasher)
	return controller.NewUserController(userService)
}
