package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	httpPkg "github.com/joao-vitor-felix/cinemax/internal/adapter/http"
	"github.com/joao-vitor-felix/cinemax/internal/database"
	"github.com/joao-vitor-felix/cinemax/internal/factory"
	"github.com/joho/godotenv"
)

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
	// container recebendo dependências e retornando controllers
	container := factory.NewContainer(db)
	r := httpPkg.SetupRoutes(container)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	log.Fatal(srv.ListenAndServe())
}
