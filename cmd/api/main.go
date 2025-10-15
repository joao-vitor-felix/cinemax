package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	httpAdpt "github.com/joao-vitor-felix/cinemax/internal/adapters/http"
	"github.com/joao-vitor-felix/cinemax/internal/database"
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
	// bundle := locales.NewBundle()
	r := httpAdpt.SetupRoutes()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	log.Fatal(srv.ListenAndServe())
}
