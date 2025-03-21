package migrate

import (
	"fmt"
	"log"

	"algo/src/shared"
)

type OHLCVHandler interface {
	HandleBatch(shared.InstrumentMaster, []*shared.OHLCV) error
}

var size = 1000

// Migrate processes a CSV file using the provided data provider
func Migrate(source DataSource, handler OHLCVHandler) error {
	chunks, err := ReadCSVInChunks(source.FilePath, size, source.Provider.GetHeaders())
	if err != nil {
		return fmt.Errorf("setting up CSV reader: %w", err)
	}

	for chunk := range chunks {
		var records []*shared.OHLCV

		for _, row := range chunk {
			record, err := source.Provider.MapToOHLCV(row)
			if err != nil {
				log.Printf("Error parsing row: %v", err)
				continue
			}
			records = append(records, record)
		}

		if len(records) > 0 {
			fmt.Printf("%+v\n", records[0])
			if err := handler.HandleBatch(source.InstrumentMaster, records); err != nil {
				log.Printf("Error handling batch: %v", err)
			}
		}
	}

	return nil
}
