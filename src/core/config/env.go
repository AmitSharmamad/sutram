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
