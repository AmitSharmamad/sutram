package migrate

import "algo/src/shared"

// DataProvider defines the contract for any data source provider
type DataProvider interface {
	// GetHeaders returns the required headers for validation
	GetHeaders() []string

	// MapToOHLCV converts a row of data to OHLCV struct
	MapToOHLCV(row map[string]string) (*shared.OHLCV, error)
}

// DataSource represents configuration for a data source
type DataSource struct {
	Provider         DataProvider
	FilePath         string // This could be generalized further for non-file sources
	InstrumentMaster shared.InstrumentMaster
}
