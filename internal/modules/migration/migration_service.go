package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// MigrationService handles database migrations
type MigrationService struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrationService creates a new migration service
func NewMigrationService(db *sql.DB, migrationsDir string) *MigrationService {
	return &MigrationService{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Up runs all available migrations
func (m *MigrationService) Up() error {
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
func (m *MigrationService) Down() error {
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
func (m *MigrationService) Steps(n int) error {
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
func (m *MigrationService) Force(version int) error {
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
func (m *MigrationService) Version() (uint, bool, error) {
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
func (m *MigrationService) Validate() error {
	migrator, err := m.createMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer migrator.Close()

	// Check if migrations directory exists
	if _, err := os.Stat(m.migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", m.migrationsDir)
	}

	// Attempt to get current version to validate migration setup
	_, _, err = migrator.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	return nil
}

// createMigrate creates a new migrate instance
func (m *MigrationService) createMigrate() (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.migrationsDir),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return migrator, nil
}
