package db

import (
	"database/sql"
	"fmt"
	"log"

	env "algo/src/core/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	env *env.EnvProv
	db  *sqlx.DB
}

func NewPostgresProv(env *env.EnvProv) *Postgres {
	db := &Postgres{
		env: env,
	}

	db.Connect()

	return db
}

func (postgres *Postgres) GetDB() *sqlx.DB {
	return postgres.db
}

func (postgres *Postgres) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return postgres.db.NamedQuery(query, arg)
}

func (postgres *Postgres) Connect() {
	envProv := postgres.env

	host := envProv.GetProp(env.POSTGRES_HOST)
	port := envProv.GetProp(env.POSTGRES_PORT)
	dbname := envProv.GetProp(env.POSTGRES_DBNAME)
	username := envProv.GetProp(env.POSTGRES_USERNAME)
	password := envProv.GetProp(env.POSTGRES_PASSWORD)

	// PostgreSQL connection string format
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname,
	)

	db, err := sqlx.Open("postgres", connectionString)

	postgres.db = db

	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}

	// Ping the database to verify connectivity
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
}

func (postgres *Postgres) Close() {
	if postgres.db != nil {
		postgres.db.Close()
	}
}

func (postgres *Postgres) Insert(query string, values interface{}) (sql.Result, error) {
	result, err := postgres.db.NamedExec(query, values)
	return result, err
}

func (postgres *Postgres) Select(model any, query string, args ...any) error {
	err := postgres.db.Select(model, query, args...)
	return err
}

func (postgres *Postgres) Get(model any, query string, args ...any) error {
	err := postgres.db.Get(model, query, args...)
	return err
}

func (postgres *Postgres) Update(query string, values interface{}) (sql.Result, error) {
	return postgres.db.NamedExec(query, values)
}

func (postgres *Postgres) Delete(query string, values interface{}) (sql.Result, error) {
	return postgres.db.NamedExec(query, values)
}
