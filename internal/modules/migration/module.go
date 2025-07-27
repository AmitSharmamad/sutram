package migration

import (
	"fmt"
	"os"
	"path/filepath"

	"algo/internal/infra/persistance/repository"
)

// Module represents the migration module
type Module struct {
	service *MigrationService
}

// NewModule creates a new migration module
func NewModule() (*Module, error) {
	// Get database connection
	db, err := repository.ProvidePostgresDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Get absolute path to migrations directory
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	migrationsDir := filepath.Join(workingDir, "db", "migrations", "postgres")

	// Create migration service
	service := NewMigrationService(db, migrationsDir)

	return &Module{
		service: service,
	}, nil
}

// GetService returns the migration service
func (m *Module) GetService() *MigrationService {
	return m.service
}

// Up runs all available migrations
func (m *Module) Up() error {
	return m.service.Up()
}

// Down runs all down migrations
func (m *Module) Down() error {
	return m.service.Down()
}

// Steps runs n migration steps
func (m *Module) Steps(n int) error {
	return m.service.Steps(n)
}

// Force sets the migration version without running migrations
func (m *Module) Force(version int) error {
	return m.service.Force(version)
}

// Version returns the current migration version
func (m *Module) Version() (uint, bool, error) {
	return m.service.Version()
}

// Validate checks if all migrations are valid
func (m *Module) Validate() error {
	return m.service.Validate()
}
