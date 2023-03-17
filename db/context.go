package db

import (
	"context"
	"errors"

	"github.com/belak/toolkit/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	querrierContextKey = internal.ContextKey("querrier")
	ErrNoQuerrier      = errors.New("Querrier missing from context")
)

func WithQuerier(ctx context.Context, q Querier) context.Context {
	return context.WithValue(ctx, querrierContextKey, q)
}

func ExtractQuerrier(ctx context.Context) (Querier, bool) {
	ret := ctx.Value(querrierContextKey)

	if q, ok := ret.(Querier); ok {
		return q, true
	}

	return nil, false
}

func Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	q, ok := ExtractQuerrier(ctx)
	if !ok {
		return pgconn.CommandTag{}, ErrNoQuerrier
	}

	return q.Exec(ctx, query, args...)
}

func Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	q, ok := ExtractQuerrier(ctx)
	if !ok {
		return nil, ErrNoQuerrier
	}

	return q.Query(ctx, query, args...)
}

func QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	q, ok := ExtractQuerrier(ctx)
	if !ok {
		return &erroredRow{}
	}

	return q.QueryRow(ctx, query, args...)
}

// We need a minimal pgx.Row type to return ErrNoQuerier if the querier wasn't
// found in the context
var _ pgx.Row = (*erroredRow)(nil)

type erroredRow struct{}

func (*erroredRow) Scan(targets ...any) error {
	return ErrNoQuerrier
}
