package migration

import (
	"fmt"
	"os"
	"path/filepath"

	"algo/internal/infra/persistance/repository"
)

// Module represents the migration module
type Module struct {
	pgService *MigrationService
	chService *ClickHouseMigrationService
}

// NewModule creates a new migration module with both PostgreSQL and ClickHouse support
func NewModule() (*Module, error) {
	// Get PostgreSQL database connection
	pgDB, err := repository.ProvidePostgresDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get PostgreSQL database connection: %w", err)
	}

	// Get ClickHouse database connection
	chSQLXDB, err := repository.ProvideClickHouseDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get ClickHouse database connection: %w", err)
	}

	// Convert sqlx.DB to sql.DB for ClickHouse
	chDB := chSQLXDB.DB

	// Get absolute path to migrations directory
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	pgMigrationsDir := filepath.Join(workingDir, "db", "migrations", "postgres")
	chMigrationsDir := filepath.Join(workingDir, "db", "migrations", "clickhouse")

	// Create migration services
	pgService := NewMigrationService(pgDB, pgMigrationsDir)
	chService := NewClickHouseMigrationService(chDB, chMigrationsDir)

	return &Module{
		pgService: pgService,
		chService: chService,
	}, nil
}

// NewPostgresModule creates a new migration module for PostgreSQL only
func NewPostgresModule() (*Module, error) {
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
		pgService: service,
	}, nil
}

// NewClickHouseModule creates a new migration module for ClickHouse only
func NewClickHouseModule() (*Module, error) {
	// Get ClickHouse database connection
	chSQLXDB, err := repository.ProvideClickHouseDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get ClickHouse database connection: %w", err)
	}

	// Convert sqlx.DB to sql.DB for ClickHouse
	chDB := chSQLXDB.DB

	// Get absolute path to migrations directory
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	chMigrationsDir := filepath.Join(workingDir, "db", "migrations", "clickhouse")

	// Create migration service
	chService := NewClickHouseMigrationService(chDB, chMigrationsDir)

	return &Module{
		chService: chService,
	}, nil
}

// GetPostgresService returns the PostgreSQL migration service
func (m *Module) GetPostgresService() *MigrationService {
	return m.pgService
}

// GetClickHouseService returns the ClickHouse migration service
func (m *Module) GetClickHouseService() *ClickHouseMigrationService {
	return m.chService
}

// PostgreSQL migration methods for backward compatibility
func (m *Module) Up() error {
	if m.pgService == nil {
		return fmt.Errorf("PostgreSQL migration service not initialized")
	}
	return m.pgService.Up()
}

func (m *Module) Down() error {
	if m.pgService == nil {
		return fmt.Errorf("PostgreSQL migration service not initialized")
	}
	return m.pgService.Down()
}

func (m *Module) Steps(n int) error {
	if m.pgService == nil {
		return fmt.Errorf("PostgreSQL migration service not initialized")
	}
	return m.pgService.Steps(n)
}

func (m *Module) Force(version int) error {
	if m.pgService == nil {
		return fmt.Errorf("PostgreSQL migration service not initialized")
	}
	return m.pgService.Force(version)
}

func (m *Module) Version() (uint, bool, error) {
	if m.pgService == nil {
		return 0, false, fmt.Errorf("PostgreSQL migration service not initialized")
	}
	return m.pgService.Version()
}

func (m *Module) Validate() error {
	if m.pgService == nil {
		return fmt.Errorf("PostgreSQL migration service not initialized")
	}
	return m.pgService.Validate()
}

// ClickHouse migration methods
func (m *Module) UpClickHouse() error {
	if m.chService == nil {
		return fmt.Errorf("ClickHouse migration service not initialized")
	}
	return m.chService.Up()
}

func (m *Module) DownClickHouse() error {
	if m.chService == nil {
		return fmt.Errorf("ClickHouse migration service not initialized")
	}
	return m.chService.Down()
}

func (m *Module) StepsClickHouse(n int) error {
	if m.chService == nil {
		return fmt.Errorf("ClickHouse migration service not initialized")
	}
	return m.chService.Steps(n)
}

func (m *Module) ForceClickHouse(version int) error {
	if m.chService == nil {
		return fmt.Errorf("ClickHouse migration service not initialized")
	}
	return m.chService.Force(version)
}

func (m *Module) VersionClickHouse() (uint, bool, error) {
	if m.chService == nil {
		return 0, false, fmt.Errorf("ClickHouse migration service not initialized")
	}
	return m.chService.Version()
}

func (m *Module) ValidateClickHouse() error {
	if m.chService == nil {
		return fmt.Errorf("ClickHouse migration service not initialized")
	}
	return m.chService.Validate()
}
