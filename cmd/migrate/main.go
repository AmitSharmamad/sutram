package main

import (
	"fmt"
	"log"

	"algo/src/migrate"
	"algo/src/shared"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	container := BuildContainer(engine)

	// Configure DataFeeds provider

	// Initialize provider with date/time formats
	provider := migrate.NewDataFeedsProvider()

	// Create data source configuration
	source := migrate.DataSource{
		Provider: provider,
		FilePath: "data/historical",
		InstrumentMaster: shared.InstrumentMaster{
			Exchange: shared.MCX,
			Segment:  shared.FNO,
		},
	}

	// Start migration
	if err := migrate.Migrate(source, container.dbHandler); err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	fmt.Println("Migration completed successfully")
}
