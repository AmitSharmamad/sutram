package migration

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// OHLCV represents a single candlestick data point
type OHLCV struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
}

// CSVMapper maps CSV column names to their respective field names
type CSVMapper struct {
	TimestampCol string // Combined timestamp column (optional)
	DateCol      string // Separate date column (optional)
	TimeCol      string // Separate time column (optional)
	OpenCol      string
	HighCol      string
	LowCol       string
	CloseCol     string
	VolumeCol    string
	TimeFormat   string // Format for combined timestamp
	DateFormat   string // Format for separate date
	TimeFormat24 string // Format for separate time
}

// ParseCSVRow converts a row of string data to OHLCV struct
func ParseCSVRow(row map[string]string, mapper CSVMapper) (*OHLCV, error) {
	data := &OHLCV{}
	var err error

	// Parse timestamp based on whether we have combined or separate date/time columns
	if mapper.TimestampCol != "" {
		// Use combined timestamp
		ts, err := time.Parse(mapper.TimeFormat, row[mapper.TimestampCol])
		if err != nil {
			return nil, fmt.Errorf("parsing timestamp: %w", err)
		}
		data.Timestamp = ts
	} else if mapper.DateCol != "" && mapper.TimeCol != "" {
		// Use separate date and time columns
		dateStr := row[mapper.DateCol]
		timeStr := row[mapper.TimeCol]

		// Parse date
		date, err := time.Parse(mapper.DateFormat, dateStr)
		if err != nil {
			return nil, fmt.Errorf("parsing date: %w", err)
		}

		// Parse time
		timeVal, err := time.Parse(mapper.TimeFormat24, timeStr)
		if err != nil {
			return nil, fmt.Errorf("parsing time: %w", err)
		}

		// Combine date and time
		data.Timestamp = time.Date(
			date.Year(), date.Month(), date.Day(),
			timeVal.Hour(), timeVal.Minute(), timeVal.Second(), 0,
			date.Location(),
		)
	} else {
		return nil, fmt.Errorf("either TimestampCol or both DateCol and TimeCol must be provided")
	}

	// Parse OHLCV values
	data.Open, err = strconv.ParseFloat(row[mapper.OpenCol], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing open: %w", err)
	}

	data.High, err = strconv.ParseFloat(row[mapper.HighCol], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing high: %w", err)
	}

	data.Low, err = strconv.ParseFloat(row[mapper.LowCol], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing low: %w", err)
	}

	data.Close, err = strconv.ParseFloat(row[mapper.CloseCol], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing close: %w", err)
	}

	vol, err := strconv.ParseFloat(row[mapper.VolumeCol], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing volume: %w", err)
	}
	data.Volume = int64(vol)

	return data, nil
}

// ReadCSVInChunks reads a CSV file in chunks and processes each chunk
func ReadCSVInChunks(filePath string, chunkSize int, headers []string) (chan []map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	reader := csv.NewReader(file)
	// Read and validate headers
	fileHeaders, err := reader.Read()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("reading headers: %w", err)
	}

	// Create header map for quick lookups
	headerMap := make(map[string]int)
	for i, h := range fileHeaders {
		headerMap[h] = i
	}

	// Validate required headers exist
	for _, required := range headers {
		if _, exists := headerMap[required]; !exists {
			file.Close()
			return nil, fmt.Errorf("required header not found: %s", required)
		}
	}

	chunks := make(chan []map[string]string)

	go func() {
		defer file.Close()
		defer close(chunks)

		chunk := make([]map[string]string, 0, chunkSize)

		for {
			record, err := reader.Read()
			if err == io.EOF {
				if len(chunk) > 0 {
					chunks <- chunk
				}
				break
			}
			if err != nil {
				// Log error but continue processing
				fmt.Printf("Error reading row: %v\n", err)
				continue
			}

			row := make(map[string]string)
			for header, idx := range headerMap {
				row[header] = record[idx]
			}

			chunk = append(chunk, row)
			if len(chunk) >= chunkSize {
				chunks <- chunk
				chunk = make([]map[string]string, 0, chunkSize)
			}
		}
	}()

	return chunks, nil
}
