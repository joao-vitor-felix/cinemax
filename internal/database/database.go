package database

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	m "github.com/joao-vitor-felix/cinemax/migrations"
	"github.com/pressly/goose/v3"
)

func OpenPool() *sql.DB {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	slog.Info("Connected to database")
	return db
}

func RunMigrations(db *sql.DB, dir string) error {
	if dir == "" {
		dir = "../../migrations"
	}
	goose.SetBaseFS(m.MigrationFS)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, dir); err != nil {
		panic(err)
	}
	return nil
}
