package db

import "database/sql"

const (
	CLICKHOUSE Db = "clickhouse"
	POSTGRES   Db = "postgres"
)

type Database interface {
	Connect()
	Close()

	Insert(query string, args ...any) (sql.Result, error)
	Select(query string, args ...any) ([]any, error)
	Update(query string, args ...any) (sql.Result, error)
	Delete(query string, args ...any) (sql.Result, error)
}

type Db string
