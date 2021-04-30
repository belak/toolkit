package db

import "github.com/jmoiron/sqlx"

type Rows struct {
	r *sqlx.Rows
}

func (r *Rows) Next() bool   { return r.r.Next() }
func (r *Rows) Err() error   { return r.r.Err() }
func (r *Rows) Close() error { return r.r.Close() }
func (r *Rows) Scan(dest interface{}) error {
	return r.r.StructScan(dest)
}
