package migrate

import (
	"fmt"
	"time"
)

// Module represents the migration module
type Module struct {
	migrator                 *DataMigrator
	instrumentMasterMigrator *InstrumentMasterMigrator
}

// NewModule creates a new migration module
func NewModule(accessToken string) *Module {
	instrumentMasterMigrator, err := NewInstrumentMasterMigrator()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize instrument master migrator: %v\n", err)
	}

	return &Module{
		migrator:                 NewDataMigrator(accessToken),
		instrumentMasterMigrator: instrumentMasterMigrator,
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

// MigrateInstrumentMasterFromCSV migrates instrument master data from CSV file
func (m *Module) MigrateInstrumentMasterFromCSV(filePath string, truncateFirst bool) error {
	if m.instrumentMasterMigrator == nil {
		return fmt.Errorf("instrument master migrator not initialized")
	}

	overallStartTime := time.Now()
	fmt.Printf("=== Instrument Master Migration Started ===\n")

	if truncateFirst {
		fmt.Println("Truncating existing instrument master data...")
		err := m.instrumentMasterMigrator.TruncateInstrumentMaster()
		if err != nil {
			return fmt.Errorf("failed to truncate instrument master: %w", err)
		}
	}

	fmt.Printf("Starting instrument master migration from CSV file: %s\n", filePath)
	err := m.instrumentMasterMigrator.MigrateFromCSV(filePath)
	if err != nil {
		return fmt.Errorf("instrument master migration failed: %w", err)
	}

	overallDuration := time.Since(overallStartTime)
	fmt.Printf("=== Instrument Master Migration Completed Successfully ===\n")
	fmt.Printf("Overall migration time (including truncation): %s\n", overallDuration.String())
	return nil
}
