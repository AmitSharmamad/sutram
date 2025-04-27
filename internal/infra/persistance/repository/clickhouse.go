package repository

import (
	"os"
	"sync"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
)

// Global variable to store the ClickHouse DB connection
var (
	clickHouseDB *sqlx.DB
	clickHouseMu sync.Mutex
)

// ProvideClickHouseDB provides a ClickHouse database connection
func ProvideClickHouseDB() (*sqlx.DB, error) {
	clickHouseMu.Lock()
	defer clickHouseMu.Unlock()

	// Return existing connection if available
	if clickHouseDB != nil {
		return clickHouseDB, nil
	}

	conn := os.Getenv("CH_URL")
	if conn == "" {
		conn = "clickhouse://admin:admin@localhost:9000"
	}

	db, err := sqlx.Connect("clickhouse", conn)
	if err != nil {
		return nil, err
	}

	// Save the connection to the global variable
	clickHouseDB = db
	return clickHouseDB, nil
}
