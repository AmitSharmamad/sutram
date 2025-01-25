package db

import (
	env "algo/src/core/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type Clickhouse struct {
	env *env.EnvProv
	db  *sql.DB
}

func NewClickHouseProv(env *env.EnvProv) *Clickhouse {
	db := &Clickhouse{
		env: env,
	}

	db.Connect()

	return db
}

func (clickhouse *Clickhouse) Connect() {
	envProv := clickhouse.env

	host := envProv.GetProp(env.CLICKHOUSE_HOST)
	port := envProv.GetProp(env.CLICKHOUSE_PORT)

	username := envProv.GetProp(env.CLICKHOUSE_USERNAME)
	password := envProv.GetProp(env.CLICKHOUSE_PASSWORD)

	println(host)
	println(port)
	println(username)
	println(password)

	connectionString := "tcp://" + host + ":" + port + "?username=" + "admin" + "&password=admin"

	println("from inside clickhouse")

	db, err := sql.Open("clickhouse", connectionString)

	clickhouse.db = db

	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}

	// defer db.Close()

	// Ping the database to verify connectivity.
	if err := db.Ping(); err != nil {
		// If there is an error, attempt type assertion to get more details.
		// The clickhouse-go driver returns a *clickhouse.Exception on errors.
		// However, if you're not using the clickhouse-go package explicitly, you can just log the error.
		log.Fatalf("Error connecting to ClickHouse: %v", err)
	}

	fmt.Println("Successfully connected to ClickHouse!")
}

func (ch *Clickhouse) Close() {
	ch.db.Close()
}

func (ch *Clickhouse) Insert(query string, args ...any) (sql.Result, error) {
	return ch.db.Exec(query)
}

func (ch *Clickhouse) SelectOne(query string, args ...any) (any, error) {
	// TODO: do this
	results, err := ch.db.Exec("SELECT COUNT(*) FROM sutram.stock_TESLA;")
	println(results)
	return nil, err
}

func (ch *Clickhouse) Select(query string, args ...any) ([]any, error) {
	// TODO: do this
	_, err := ch.db.Exec("SELECT VERSION();")
	return []any{""}, err
}

func (ch *Clickhouse) Update(query string, args ...any) (sql.Result, error) {
	return ch.db.Exec(query)
}

func (ch *Clickhouse) Delete(query string, args ...any) (sql.Result, error) {
	return ch.db.Exec(query)
}
