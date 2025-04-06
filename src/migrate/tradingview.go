package migrate

import (
	"fmt"
	"strconv"
	"time"

	"algo/src/shared"
)

// TradingViewProvider implements DataProvider interface for TradingView's data format
type TradingViewProvider struct {
	ColumnMap  map[string]string
	DateFormat string
}

// NewTradingViewProvider creates a new provider with column mapping and time formats
func NewTradingViewProvider() *TradingViewProvider {
	columnMap := map[string]string{
		"timestamp": "time", // TradingView uses ISO 8601 format
		"open":      "open",
		"high":      "high",
		"low":       "low",
		"close":     "close",
		"symbol":    "symbol",
	}

	return &TradingViewProvider{
		ColumnMap:  columnMap,
		DateFormat: time.RFC3339, // Use RFC3339 for ISO 8601 parsing
	}
}

// GetHeaders returns the required headers for validation
func (tv *TradingViewProvider) GetHeaders() []string {
	headers := make([]string, 0, len(tv.ColumnMap))
	for _, col := range tv.ColumnMap {
		headers = append(headers, col)
	}
	return headers
}

// MapToOHLCV converts a row of data to OHLCV struct
func (tv *TradingViewProvider) MapToOHLCV(row map[string]string) (*shared.OHLCV, error) {
	data := &shared.OHLCV{}
	var err error

	// Parse ISO 8601 timestamp
	timestampStr := row[tv.ColumnMap["timestamp"]]
	data.Timestamp, err = time.Parse(tv.DateFormat, timestampStr)
	if err != nil {
		return nil, fmt.Errorf("parsing timestamp: %w", err)
	}

	// Set instrument/symbol
	data.Instrument = row[tv.ColumnMap["symbol"]]

	// Parse OHLC values
	data.Open, err = strconv.ParseFloat(row[tv.ColumnMap["open"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing open: %w", err)
	}

	data.High, err = strconv.ParseFloat(row[tv.ColumnMap["high"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing high: %w", err)
	}

	data.Low, err = strconv.ParseFloat(row[tv.ColumnMap["low"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing low: %w", err)
	}

	data.Close, err = strconv.ParseFloat(row[tv.ColumnMap["close"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing close: %w", err)
	}

	// Set volume to 0 since it's not present in TradingView data
	data.Volume = 0

	return data, nil
}
