package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"algo/internal/modules/migration"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Define command-line flags
	var (
		action  = flag.String("action", "up", "Migration action: up, down, steps, force, version, validate")
		steps   = flag.Int("steps", 1, "Number of steps for 'steps' action")
		version = flag.Int("version", 0, "Version number for 'force' action")
	)
	flag.Parse()

	// Create ClickHouse migration module
	migrationModule, err := migration.NewClickHouseModule()
	if err != nil {
		log.Fatalf("Failed to create ClickHouse migration module: %v", err)
	}

	// Execute the requested action
	switch *action {
	case "up":
		err = migrationModule.UpClickHouse()
		if err != nil {
			log.Fatalf("ClickHouse migration up failed: %v", err)
		}
		fmt.Println("✅ All ClickHouse migrations applied successfully")

	case "down":
		err = migrationModule.DownClickHouse()
		if err != nil {
			log.Fatalf("ClickHouse migration down failed: %v", err)
		}
		fmt.Println("✅ All ClickHouse migrations rolled back successfully")

	case "steps":
		err = migrationModule.StepsClickHouse(*steps)
		if err != nil {
			log.Fatalf("ClickHouse migration steps failed: %v", err)
		}
		if *steps > 0 {
			fmt.Printf("✅ Successfully applied %d ClickHouse migration step(s)\n", *steps)
		} else {
			fmt.Printf("✅ Successfully rolled back %d ClickHouse migration step(s)\n", -*steps)
		}

	case "force":
		if *version < 0 {
			log.Fatal("Version must be a non-negative integer")
		}
		err = migrationModule.ForceClickHouse(*version)
		if err != nil {
			log.Fatalf("ClickHouse migration force failed: %v", err)
		}
		fmt.Printf("✅ Successfully forced ClickHouse migration to version %d\n", *version)

	case "version":
		currentVersion, dirty, err := migrationModule.VersionClickHouse()
		if err != nil {
			log.Fatalf("Failed to get ClickHouse migration version: %v", err)
		}
		if dirty {
			fmt.Printf("📊 Current ClickHouse migration version: %d (dirty)\n", currentVersion)
		} else {
			fmt.Printf("📊 Current ClickHouse migration version: %d\n", currentVersion)
		}

	case "validate":
		err = migrationModule.ValidateClickHouse()
		if err != nil {
			log.Fatalf("ClickHouse migration validation failed: %v", err)
		}
		fmt.Println("✅ All ClickHouse migrations are valid")

	default:
		log.Fatalf("Unknown action: %s. Valid actions are: up, down, steps, force, version, validate", *action)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nClickHouse Database Migration Tool for Sutram Trading Platform\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -action=up                    # Run all pending ClickHouse migrations\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=down                  # Rollback all ClickHouse migrations\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=steps -steps=2        # Run 2 ClickHouse migration steps\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=steps -steps=-1       # Rollback 1 ClickHouse migration\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=force -version=5      # Force ClickHouse migration to version 5\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=version               # Show current ClickHouse migration version\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=validate              # Validate all ClickHouse migrations\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNote: This tool specifically handles ClickHouse database migrations.\n")
		fmt.Fprintf(os.Stderr, "For PostgreSQL migrations, use the db-migrate command.\n")
	}
}
