package repository

import (
	"database/sql"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

// Global variable to store the PostgreSQL DB connection
var (
	postgresDB *sql.DB
	postgresMu sync.Mutex
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
