package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	if os.Getenv("APP_ENV") == "development" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file: ", err)
		}
	}
	var port int
	var migrate bool
	var migrateOnly bool
	flag.BoolVar(&migrate, "migrate", false, "run database migrations")
	flag.BoolVar(&migrateOnly, "migrate-only", false, "run database migrations and exit")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()
	db := database.OpenPool()

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if migrate || migrateOnly {
		err := database.RunMigrations(db)
		if err != nil {
			log.Fatal(err)
		}
		if migrateOnly {
			return
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Server has started", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	slog.Info("Server stopped gracefully")
}
