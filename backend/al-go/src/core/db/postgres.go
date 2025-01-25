package db

import "database/sql"

type Postgres struct{}

func NewPostgresProv() *Postgres {
	return &Postgres{}
}

func (*Postgres) Connect() {
}

func (*Postgres) Close() {
}

func (*Postgres) Insert(query string, args ...any) (sql.Result, error) {
	return nil, nil
}

func (*Postgres) Select(model []any, query string, args ...any) ([]any, error) {
	return nil, nil
}

func (*Postgres) Update(query string, args ...any) (sql.Result, error) {
	return nil, nil
}

func (*Postgres) Delete(query string, args ...any) (sql.Result, error) {
	return nil, nil
}
