package factory

import (
	"database/sql"

	"github.com/joao-vitor-felix/cinemax/internal/adapter/controller"
)

type Container struct {
	UserController *controller.UserController
}

func NewContainer(db *sql.DB) *Container {
	return &Container{
		UserController: MakeNewUserController(db),
	}
}
