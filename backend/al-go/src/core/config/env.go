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
)

func GetProp(property EnvProp) string {
	return os.Getenv(string(property))
}

func Init() {

	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	println("PORT: " + GetProp(PORT))

}
