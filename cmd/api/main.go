package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	httpPkg "github.com/joao-vitor-felix/cinemax/internal/adapter/http"
	"github.com/joao-vitor-felix/cinemax/internal/database"
	"github.com/joao-vitor-felix/cinemax/internal/factory"
	"github.com/joho/godotenv"
)

//	@title			Cinemax API
//	@version		1.0
//	@description	This is the API documentation for the Cinemax application.

//	@host		localhost:8080
//	@BasePath	/

//	@accept		json
//	@produce	json

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	var port int
	var migrate bool
	flag.BoolVar(&migrate, "migrate", false, "run database migrations")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()
	db := database.OpenPool()
	defer db.Close()
	if migrate {
		err := database.RunMigrations(db)
		if err != nil {
			log.Fatal(err)
		}
	}
	container := factory.NewContainer(db)
	r := httpPkg.SetupRoutes(container)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	slog.Info("Server has started", "port", port)
	log.Fatal(srv.ListenAndServe())
}
