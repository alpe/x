// Package pg provides abstraction over jmoiron/sqlx package.
package pg

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/net/context"
)

// Getter is generic interface for getting single entity
type Getter interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

// Getter is generic interface for getting multiple enties
type Selector interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

// Getter is generic interface for executing SQL query with no result
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// CastErr inspect given error and replace generic SQL error with easier to
// compare equivalent.
//
// See http://www.postgresql.org/docs/current/static/errcodes-appendix.html
func CastErr(err error) error {
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23505":
			return ErrConflict
		case "23503":
			return ErrForeignKeyViolation
		}
	}
	return err
}

var (
	// ErrNotFound is returned when result was expected but not returned
	ErrNotFound = errors.New("not found")

	// ErrConflict is returned when database query cannot be fulfilled because
	// of constraint conflict
	ErrConflict = errors.New("conflict")

	// ErrForeignKeyViolation is returned when a insert or update on table
	// violates foreign key constraint
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

// WithDB return context with given database instance bind to it. Use DB(ctx)
// to get database back.
func WithDB(ctx context.Context, db *sql.DB) context.Context {
	dbx := sqlx.NewDb(db, "postgres")
	return context.WithValue(ctx, "pg:db", &sqlxdb{dbx})
}

// DB return database instance carried by given context.
func DB(ctx context.Context) Database {
	db := ctx.Value("pg:db")
	if db == nil {
		panic("missing database in context")
	}
	return db.(Database)
}

// sqlxdb wraps sqlx.DB structure and provides custom function notations that
// can be easily mocked. This wrapper is required, because of sqlx.DB's Beginx
// method notation
type sqlxdb struct {
	dbx *sqlx.DB
}

func (x *sqlxdb) Beginx() (Connection, error) {
	return x.dbx.Beginx()
}

func (x *sqlxdb) Get(dest interface{}, query string, args ...interface{}) error {
	return x.dbx.Get(dest, query, args...)
}

func (x *sqlxdb) Select(dest interface{}, query string, args ...interface{}) error {
	return x.dbx.Select(dest, query, args...)
}

func (x *sqlxdb) Exec(query string, args ...interface{}) (sql.Result, error) {
	return x.dbx.Exec(query, args...)
}

type Database interface {
	Beginx() (Connection, error)
	Getter
	Selector
	Execer
}

type Connection interface {
	Getter
	Selector
	Execer
	Rollback() error
	Commit() error
}
