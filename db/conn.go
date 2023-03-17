package db

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// Ensure Tx and Conn conform to our interface.
	_ Querier = (pgx.Tx)(nil)
	_ Querier = (*pgx.Conn)(nil)

	// Ensure our interfaces conform to pgxscan.Querier.
	_ pgxscan.Querier = (Querier)(nil)
	_ pgxscan.Querier = (Conn)(nil)
)

// Querier is something that pgxscan can query and get the pgx.Rows from. For
// example, it can be: *pgxpool.Pool, *pgx.Conn, pgx.Tx, or something entirely
// different like our Conn abstraction.
//
// This is a superset of the pgxscan.Querier interface which adds Exec and
// QueryRow.
type Querier interface {
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type Conn interface {
	Querier
	Tx(context.Context, func(context.Context, Querier) error) error
}

type conn struct {
	*pgxpool.Pool
}

func NewDBPool(ctx context.Context, dsn string) (Conn, error) {
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	dbConn, err := dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &conn{dbPool}, nil
}

func (c *conn) Tx(ctx context.Context, f func(context.Context, Querier) error) error {
	tx, err := c.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = f(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func NotFound(err error) bool {
	return pgxscan.NotFound(err)
}
