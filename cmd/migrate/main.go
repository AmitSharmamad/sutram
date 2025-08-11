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
	migrationType := flag.String("type", "historical", "Migration type: 'historical' for price data, 'instrument-master' for instrument master data, or 'generate-sample' to create sample Excel file")
	csvFile := flag.String("csv", "", "Path to CSV file containing instrument master data (required for instrument-master type)")
	truncateFirst := flag.Bool("truncate", false, "Truncate existing instrument master data before migration (only for instrument-master type)")
	flag.Parse()

	// Create migration module
	module := migrate.NewModule(*accessToken)

	switch *migrationType {
	case "historical":
		if *accessToken == "" {
			fmt.Println("Error: Access token is required for historical data migration. Set it using the -token flag or DHAN_ACCESS_TOKEN environment variable.")
			os.Exit(1)
		}

		fmt.Println("Running historical data migration...")
		err := module.RunMigration()
		if err != nil {
			fmt.Printf("Historical data migration failed: %v\n", err)
			os.Exit(1)
		}

	case "instrument-master":
		if *csvFile == "" {
			fmt.Println("Error: CSV file path is required for instrument master migration. Use the -csv flag.")
			flag.Usage()
			os.Exit(1)
		}

		// Check if file exists
		if _, err := os.Stat(*csvFile); os.IsNotExist(err) {
			fmt.Printf("Error: CSV file does not exist: %s\n", *csvFile)
			os.Exit(1)
		}

		fmt.Printf("Running instrument master migration from CSV file: %s\n", *csvFile)
		err := module.MigrateInstrumentMasterFromCSV(*csvFile, *truncateFirst)
		if err != nil {
			fmt.Printf("Instrument master migration failed: %v\n", err)
			os.Exit(1)
		}

	case "generate-sample":
		fmt.Println("Generating sample CSV file for testing...")
		// Note: Sample generation functionality removed for CSV-only version
		// You can create a test CSV manually or use head command on existing files
		fmt.Println("To create a test sample, use: head -100 your_file.csv > test_sample.csv")
		return

	default:
		fmt.Printf("Error: Unknown migration type: %s\n", *migrationType)
		fmt.Println("Supported types: 'historical', 'instrument-master', 'generate-sample'")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Migration completed successfully!")
}
