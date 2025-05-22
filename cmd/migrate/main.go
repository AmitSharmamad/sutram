package main

import (
	"flag"
	"fmt"
	"os"

	"algo/internal/modules/migrate"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Define command-line flags
	accessToken := flag.String("token", os.Getenv("DHAN_ACCESS_TOKEN"), "Access token for Dhan API")
	flag.Parse()

	if *accessToken == "" {
		fmt.Println("Error: Access token is required. Set it using the -token flag or DHAN_ACCESS_TOKEN environment variable.")
		os.Exit(1)
	}

	// Create and run migration module
	module := migrate.NewModule(*accessToken)
	err := module.RunMigration()
	if err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		os.Exit(1)
	}
}
