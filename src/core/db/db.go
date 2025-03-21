package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB string

const (
	CLICKHOUSE DB = "clickhouse"
	POSTGRES   DB = "postgres"
)

type Database interface {
	Connect()
	Close()

	GetDB() *sqlx.DB

	Insert(query string, values interface{}) (sql.Result, error)
	Select(model any, query string, args ...any) error
	Get(model any, query string, args ...any) error
	Update(query string, values interface{}) (sql.Result, error)
	Delete(query string, values interface{}) (sql.Result, error)

	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}
