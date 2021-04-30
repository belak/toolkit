package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/belak/toolkit/internal"
)

var querrierContextKey = internal.ContextKey("querrier")
var ErrNoQuerrier = errors.New("Querrier missing from context")

func WithQuerier(ctx context.Context, q Querier) context.Context {
	return context.WithValue(ctx, querrierContextKey, q)
}

func GetQuerrier(ctx context.Context) (Querier, bool) {
	ret := ctx.Value(querrierContextKey)

	if q, ok := ret.(Querier); ok {
		return q, true
	}

	return nil, false
}

func Exec(ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	q, ok := GetQuerrier(ctx)
	if !ok {
		return nil, ErrNoQuerrier
	}

	return q.Exec(ctx, query, params...)
}

func Get(ctx context.Context, dest interface{}, query string, params ...interface{}) error {
	q, ok := GetQuerrier(ctx)
	if !ok {
		return ErrNoQuerrier
	}

	return q.Get(ctx, dest, query, params...)
}

func Select(ctx context.Context, dest interface{}, query string, params ...interface{}) error {
	q, ok := GetQuerrier(ctx)
	if !ok {
		return ErrNoQuerrier
	}

	return q.Select(ctx, dest, query, params...)
}

func Query(ctx context.Context, query string, params ...interface{}) (*Rows, error) {
	q, ok := GetQuerrier(ctx)
	if !ok {
		return nil, ErrNoQuerrier
	}

	return q.Query(ctx, query, params...)
}
