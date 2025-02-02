package db

import (
	"database/sql"
	"fmt"
	"log"

	env "algo/src/core/config"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
)

type Clickhouse struct {
	env *env.EnvProv
	db  *sqlx.DB
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

	db, err := sqlx.Open("clickhouse", connectionString)

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

func (clickhouse *Clickhouse) Close() {
	clickhouse.db.Close()
}

func (clickhouse *Clickhouse) Insert(query string, values interface{}) (sql.Result, error) {
	result, err := clickhouse.db.NamedExec(query, values)
	return result, err
}

func (clickhouse *Clickhouse) Select(model any, query string, args ...any) error {
	err := clickhouse.db.Select(model, query, args...)
	return err
}

func (clickhouse *Clickhouse) Update(query string, values interface{}) (sql.Result, error) {
	return clickhouse.db.NamedExec(query, values)
}

func (clickhouse *Clickhouse) Delete(query string, values interface{}) (sql.Result, error) {
	return clickhouse.db.NamedExec(query, values)
}
