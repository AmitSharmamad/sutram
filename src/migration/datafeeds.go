package migration

import (
	"fmt"
	"strconv"
	"time"
)

// DataFeedsProvider implements DataProvider interface for DataFeeds's CSV format
type DataFeedsProvider struct {
	ColumnMap  map[string]string
	DateFormat string
	TimeFormat string
}

// NewDataFeedsProvider creates a new provider with column mapping and time formats
func NewDataFeedsProvider() *DataFeedsProvider {

	columnMap := map[string]string{
		"date":   "Date", // Assuming these are the column names in your CSV
		"time":   "Time", // Adjust these according to your actual CSV headers
		"open":   "Open",
		"high":   "High",
		"low":    "Low",
		"close":  "Close",
		"volume": "Volume",
	}

	return &DataFeedsProvider{
		ColumnMap:  columnMap,
		DateFormat: "2006-01-02",
		TimeFormat: "15:04:05",
	}
}

func (v *DataFeedsProvider) GetHeaders() []string {
	headers := make([]string, 0, len(v.ColumnMap))
	for _, col := range v.ColumnMap {
		headers = append(headers, col)
	}
	return headers
}

func (v *DataFeedsProvider) MapToOHLCV(row map[string]string) (*OHLCV, error) {
	data := &OHLCV{}
	var err error

	// Parse timestamp based on configuration
	if timestamp, hasTimestamp := row[v.ColumnMap["timestamp"]]; hasTimestamp {
		data.Timestamp, err = time.Parse(v.DateFormat, timestamp)
		if err != nil {
			return nil, fmt.Errorf("parsing timestamp: %w", err)
		}
	} else {
		dateStr := row[v.ColumnMap["date"]]
		timeStr := row[v.ColumnMap["time"]]

		date, err := time.Parse(v.DateFormat, dateStr)
		if err != nil {
			return nil, fmt.Errorf("parsing date: %w", err)
		}

		timeVal, err := time.Parse(v.TimeFormat, timeStr)
		if err != nil {
			return nil, fmt.Errorf("parsing time: %w", err)
		}

		data.Timestamp = time.Date(
			date.Year(), date.Month(), date.Day(),
			timeVal.Hour(), timeVal.Minute(), timeVal.Second(), 0,
			date.Location(),
		)
	}

	// Parse OHLCV values
	data.Open, err = strconv.ParseFloat(row[v.ColumnMap["open"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing open: %w", err)
	}

	data.High, err = strconv.ParseFloat(row[v.ColumnMap["high"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing high: %w", err)
	}

	data.Low, err = strconv.ParseFloat(row[v.ColumnMap["low"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing low: %w", err)
	}

	data.Close, err = strconv.ParseFloat(row[v.ColumnMap["close"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing close: %w", err)
	}

	vol, err := strconv.ParseFloat(row[v.ColumnMap["volume"]], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing volume: %w", err)
	}
	data.Volume = int64(vol)

	return data, nil
}
