package migration

import (
	"fmt"
	"log"
)

type OHLCVHandler interface {
	HandleBatch([]*OHLCV) error
}

// Migrate processes a CSV file using the provided data provider
func Migrate(source DataSource, handler OHLCVHandler) error {
	chunks, err := ReadCSVInChunks(source.FilePath, 1000, source.Provider.GetHeaders())
	if err != nil {
		return fmt.Errorf("setting up CSV reader: %w", err)
	}

	for chunk := range chunks {
		var records []*OHLCV

		for _, row := range chunk {
			record, err := source.Provider.MapToOHLCV(row)
			if err != nil {
				log.Printf("Error parsing row: %v", err)
				continue
			}
			records = append(records, record)
		}

		if len(records) > 0 {
			if err := handler.HandleBatch(records); err != nil {
				log.Printf("Error handling batch: %v", err)
			}
		}
	}

	return nil
}
