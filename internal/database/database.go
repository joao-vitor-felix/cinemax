package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	m "github.com/joao-vitor-felix/cinemax/migrations"
	"github.com/pressly/goose/v3"
)

func OpenPool() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("pgx", dsn)
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

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(m.MigrationFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	tableName := os.Getenv("GOOSE_TABLE")
	if tableName != "" {
		goose.SetTableName(tableName)
	}
	if err := goose.Up(db, "."); err != nil {
		return err
	}
	slog.Info("Migrations applied")
	return nil
}
