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

	// Create PostgreSQL migration module
	migrationModule, err := migration.NewPostgresModule()
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL migration module: %v", err)
	}

	// Execute the requested action
	switch *action {
	case "up":
		err = migrationModule.Up()
		if err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("✅ All migrations applied successfully")

	case "down":
		err = migrationModule.Down()
		if err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("✅ All migrations rolled back successfully")

	case "steps":
		err = migrationModule.Steps(*steps)
		if err != nil {
			log.Fatalf("Migration steps failed: %v", err)
		}
		fmt.Printf("✅ %d migration steps completed successfully\n", *steps)

	case "force":
		if *version == 0 {
			log.Fatalf("Version must be specified for force action")
		}
		err = migrationModule.Force(*version)
		if err != nil {
			log.Fatalf("Migration force failed: %v", err)
		}
		fmt.Printf("✅ Migration version forced to %d\n", *version)

	case "version":
		currentVersion, dirty, err := migrationModule.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		status := "clean"
		if dirty {
			status = "dirty"
		}
		fmt.Printf("Current migration version: %d (%s)\n", currentVersion, status)

	case "validate":
		err = migrationModule.Validate()
		if err != nil {
			log.Fatalf("Migration validation failed: %v", err)
		}
		fmt.Println("✅ All migrations are valid")

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, steps, force, version, validate")
		os.Exit(1)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDatabase Migration Tool for Sutram Trading Platform\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -action=up                    # Run all pending migrations\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=down                  # Rollback all migrations\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=steps -steps=2        # Run 2 migration steps\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=steps -steps=-1       # Rollback 1 migration\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=force -version=5      # Force migration to version 5\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=version               # Show current migration version\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=validate              # Validate all migrations\n", os.Args[0])
	}
}
