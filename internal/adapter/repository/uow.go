package repository

import (
	"context"
	"database/sql"
)

type txKey struct{}

type PostgresUnitOfWork struct {
	db *sql.DB
}

func NewPostgresUnitOfWork(db *sql.DB) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{db: db}
}

func (uow *PostgresUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

type executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func getExecutor(ctx context.Context, db *sql.DB) executor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return db
}
