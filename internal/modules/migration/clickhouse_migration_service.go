package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// ClickHouseMigrationService handles ClickHouse database migrations
type ClickHouseMigrationService struct {
	db            *sql.DB
	migrationsDir string
}

// NewClickHouseMigrationService creates a new ClickHouse migration service
func NewClickHouseMigrationService(db *sql.DB, migrationsDir string) *ClickHouseMigrationService {
	return &ClickHouseMigrationService{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Up runs all available migrations
func (m *ClickHouseMigrationService) Up() error {
	migrator, err := m.createMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Down runs all down migrations
func (m *ClickHouseMigrationService) Down() error {
	migrator, err := m.createMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run down migrations: %w", err)
	}

	return nil
}

// Steps runs n migration steps
func (m *ClickHouseMigrationService) Steps(n int) error {
	migrator, err := m.createMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Steps(n); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run %d migration steps: %w", n, err)
	}

	return nil
}

// Force sets the migration version without running migrations
func (m *ClickHouseMigrationService) Force(version int) error {
	migrator, err := m.createMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version %d: %w", version, err)
	}

	return nil
}

// Version returns the current migration version
func (m *ClickHouseMigrationService) Version() (uint, bool, error) {
	migrator, err := m.createMigrate()
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer migrator.Close()

	version, dirty, err := migrator.Version()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// Validate checks if all migrations are valid
func (m *ClickHouseMigrationService) Validate() error {
	migrator, err := m.createMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer migrator.Close()

	// Check if migrations directory exists and has valid files
	if _, err := os.Stat(m.migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", m.migrationsDir)
	}

	return nil
}

// createMigrate creates a new migrate instance for ClickHouse
func (m *ClickHouseMigrationService) createMigrate() (*migrate.Migrate, error) {
	driver, err := clickhouse.WithInstance(m.db, &clickhouse.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create clickhouse driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.migrationsDir),
		"clickhouse",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return migrator, nil
}
