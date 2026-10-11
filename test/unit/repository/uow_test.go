package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joao-vitor-felix/cinemax/internal/adapter/repository"
	"github.com/stretchr/testify/require"
)

func TestPostgresUnitOfWork(t *testing.T) {
	t.Run("should execute successfully and commit the transaction", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectCommit()

		uow := repository.NewPostgresUnitOfWork(db)

		err = uow.Do(context.Background(), func(ctx context.Context) error {
			// Do nothing, just return success
			return nil
		})

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should rollback transaction when callback returns an error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectRollback()

		uow := repository.NewPostgresUnitOfWork(db)

		expectedErr := errors.New("callback failed")

		err = uow.Do(context.Background(), func(ctx context.Context) error {
			return expectedErr
		})

		require.ErrorIs(t, err, expectedErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error if BeginTx fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		expectedErr := errors.New("begin failed")
		mock.ExpectBegin().WillReturnError(expectedErr)

		uow := repository.NewPostgresUnitOfWork(db)

		err = uow.Do(context.Background(), func(ctx context.Context) error {
			return nil
		})

		require.ErrorIs(t, err, expectedErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error if Commit fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectBegin()
		expectedErr := errors.New("commit failed")
		mock.ExpectCommit().WillReturnError(expectedErr)

		uow := repository.NewPostgresUnitOfWork(db)

		err = uow.Do(context.Background(), func(ctx context.Context) error {
			return nil
		})

		require.ErrorIs(t, err, expectedErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
