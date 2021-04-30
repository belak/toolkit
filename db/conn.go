package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var _ Conn = (*dbImpl)(nil)
var _ Tx = (*txImpl)(nil)

// Conn is a stripped down convenience wrapper around sqlx (which is, in-turn, a
// wrapper around database/sql) meant to make DB access more convenient and
// easier to fall into the pit of success.
type Conn interface {
	Querier

	// Begin is an alternative way of starting a transaction if you need more
	// control over it. It is strongly recommended that you use `defer
	// tx.Rollback()` directly after acquiring the transaction to ensure it is
	// closed even in the case of a panic.
	Begin(context.Context) (Tx, error)

	// Tx is wraps a transaction. If no error is returned from the callback, the
	// transaction will be committed, otherwise it will be aborted.
	Tx(context.Context, func(context.Context, Querier) error) error
}

// Querier is the common interface between a DB and a Tx.
type Querier interface {
	// Exec is meant for running queries which don't read data, such as INSERTs
	// or UPDATES.
	Exec(ctx context.Context, query string, params ...interface{}) (sql.Result, error)

	// Get is meant for SELECT queries returning a single row. It will return an
	// error if no rows are returned.
	Get(ctx context.Context, dest interface{}, query string, params ...interface{}) error

	// Select is meant for SELECT queries returning an unknown number of rows
	// but which should easily fit into memory.
	Select(ctx context.Context, dest interface{}, query string, params ...interface{}) error

	// Query is meant for SELECT queries which return an unknown number of rows
	// and may not all fit into memory. The Rows object can be used to scan the
	// results one row at a time.
	Query(ctx context.Context, query string, params ...interface{}) (*Rows, error)
}

func Open(driver string, connInfo string) (Conn, error) {
	db, err := sqlx.Open(driver, connInfo)
	if err != nil {
		return nil, err
	}

	return &dbImpl{db}, nil
}

type dbImpl struct {
	db *sqlx.DB
}

func (db *dbImpl) Begin(ctx context.Context) (Tx, error) {
	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &txImpl{tx}, nil
}

func (db *dbImpl) Tx(ctx context.Context, cb func(context.Context, Querier) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = cb(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *dbImpl) Exec(ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, params...)
}

func (db *dbImpl) Get(ctx context.Context, dest interface{}, query string, params ...interface{}) error {
	return db.db.GetContext(ctx, dest, query, params...)
}

func (db *dbImpl) Select(ctx context.Context, dest interface{}, query string, params ...interface{}) error {
	return db.db.SelectContext(ctx, dest, query, params...)
}

func (db *dbImpl) Query(ctx context.Context, query string, params ...interface{}) (*Rows, error) {
	rows, err := db.db.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return &Rows{rows}, nil
}
