package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvProp string

const (
	PORT                EnvProp = "PORT"
	CLICKHOUSE_HOST     EnvProp = "CLICKHOUSE_HOST"
	CLICKHOUSE_PORT     EnvProp = "CLICKHOUSE_PORT"
	CLICKHOUSE_USERNAME EnvProp = "CLICKHOUSE_USERNAME"
	CLICKHOUSE_PASSWORD EnvProp = "CLICKHOUSE_PASSWORD"

	POSTGRES_HOST     EnvProp = "POSTGRES_HOST"
	POSTGRES_PORT     EnvProp = "POSTGRES_PORT"
	POSTGRES_DBNAME   EnvProp = "POSTGRES_DBNAME"
	POSTGRES_USERNAME EnvProp = "POSTGRES_USERNAME"
	POSTGRES_PASSWORD EnvProp = "POSTGRES_PASSWORD"

	JWT_SIGNING_SECRET EnvProp = "JWT_SIGNING_SECRET"

	GOOGLE_CLIENT_ID     EnvProp = "GOOGLE_CLIENT_ID"
	GOOGLE_CLIENT_SECRET EnvProp = "GOOGLE_CLIENT_SECRET"
	GOOGLE_REDIRECT_URL  EnvProp = "GOOGLE_REDIRECT_URL"

	// Broker environment variables
	DHAN_CLIENT_ID      EnvProp = "DHAN_CLIENT_ID"
	DHAN_ACCESS_TOKEN   EnvProp = "DHAN_ACCESS_TOKEN"
	DHAN_BASE_URL       EnvProp = "DHAN_BASE_URL"
	DHAN_DHAN_CLIENT_ID EnvProp = "DHAN_DHAN_CLIENT_ID"

	// Zerodha environment variables
	ZERODHA_CLIENT_ID    EnvProp = "ZERODHA_CLIENT_ID"
	ZERODHA_ACCESS_TOKEN EnvProp = "ZERODHA_ACCESS_TOKEN"
	ZERODHA_BASE_URL     EnvProp = "ZERODHA_BASE_URL"
)

type EnvProv struct{}

func NewEnvProv() *EnvProv {
	prov := &EnvProv{}
	prov.Init()
	return prov
}

func (*EnvProv) GetProp(property EnvProp) string {
	return os.Getenv(string(property))
}

func (env *EnvProv) Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	println("PORT: " + env.GetProp(PORT))
}
