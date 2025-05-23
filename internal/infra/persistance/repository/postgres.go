package repository

import (
	"database/sql"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Global variables to store the SQL and SQLx DB connections
var (
	postgresDB *sql.DB
	postgresMu sync.Mutex
	sqlxDB     *sqlx.DB
	sqlxMu     sync.Mutex
)

// ProvidePostgresDB provides a PostgreSQL database connection
func ProvidePostgresDB() (*sql.DB, error) {
	postgresMu.Lock()
	defer postgresMu.Unlock()

	// Return existing connection if available
	if postgresDB != nil {
		return postgresDB, nil
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgresql://admin:admin@localhost:9009/sutram?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Save the connection to the global variable
	postgresDB = db
	return postgresDB, nil
}

// ProvidePostgresSQLXDB provides a PostgreSQL database connection using sqlx
func ProvidePostgresSQLXDB() (*sqlx.DB, error) {
	sqlxMu.Lock()
	defer sqlxMu.Unlock()

	// Return existing connection if available
	if sqlxDB != nil {
		return sqlxDB, nil
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgresql://admin:admin@localhost:9009/sutram?sslmode=disable"
	}

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Save the connection to the global variable
	sqlxDB = db
	return sqlxDB, nil
}

// ProvideSqlxDB provides a SQLx wrapper around the PostgreSQL database connection
func ProvideSqlxDB() *sqlx.DB {
	// Get the standard SQL DB connection
	db, err := ProvidePostgresDB()
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	// Wrap with sqlx
	return sqlx.NewDb(db, "postgres")
}
