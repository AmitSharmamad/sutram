package migrate

import (
	"fmt"
)

// Module represents the migration module
type Module struct {
	migrator *DataMigrator
}

// NewModule creates a new migration module
func NewModule(accessToken string) *Module {
	return &Module{
		migrator: NewDataMigrator(accessToken),
	}
}

// RunMigration executes the data migration process
func (m *Module) RunMigration() error {
	fmt.Println("Starting data migration process...")
	err := m.migrator.Run()
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	fmt.Println("Data migration completed successfully")
	return nil
}
