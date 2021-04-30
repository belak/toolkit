package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Tx is a Querier representing a transaction rather than a connection. The
// first call to either Rollback or Commit will end the transaction any calls
// afterwords will be dropped and an error returned.
type Tx interface {
	Querier

	// Rollback will roll back the transaction.
	Rollback() error

	// Commit will commit the transaction.
	Commit() error
}

type txImpl struct {
	tx *sqlx.Tx
}

func (tx *txImpl) Exec(ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, params...)
}

func (tx *txImpl) Get(ctx context.Context, dest interface{}, query string, params ...interface{}) error {
	return tx.tx.GetContext(ctx, dest, query, params...)
}

func (tx *txImpl) Select(ctx context.Context, dest interface{}, query string, params ...interface{}) error {
	return tx.tx.SelectContext(ctx, dest, query, params...)
}

func (tx *txImpl) Query(ctx context.Context, query string, params ...interface{}) (*Rows, error) {
	rows, err := tx.tx.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return &Rows{rows}, nil
}

func (tx *txImpl) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *txImpl) Commit() error {
	return tx.tx.Commit()
}
