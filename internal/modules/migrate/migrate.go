package migrate

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// SecurityData represents the input configuration for fetching security data
type SecurityData struct {
	SecurityID      string `json:"securityId"`
	ExchangeSegment string `json:"exchangeSegment"`
	Instrument      string `json:"instrument"`
	Interval        string `json:"interval"`
	OI              bool   `json:"oi"`
	FromDate        string `json:"fromDate"`
	ToDate          string `json:"toDate"`
}

// ChartResponse represents the API response structure
type ChartResponse struct {
	Open         []float64 `json:"open"`
	High         []float64 `json:"high"`
	Low          []float64 `json:"low"`
	Close        []float64 `json:"close"`
	Volume       []float64 `json:"volume"`
	Timestamp    []float64 `json:"timestamp"`
	OpenInterest []int     `json:"open_interest"`
}

// DataMigrator handles the migration of historical data
type DataMigrator struct {
	securities  []SecurityData
	apiURL      string
	accessToken string
}

// NewDataMigrator creates a new instance of DataMigrator
func NewDataMigrator(accessToken string) *DataMigrator {
	// Hardcoded securities data
	securities := []SecurityData{
		{
			SecurityID:      "2885",
			ExchangeSegment: "NSE_EQ",
			Instrument:      "EQUITY",
			Interval:        "1",
			OI:              false,
			FromDate:        "2024-09-11 09:30:00",
			ToDate:          "2024-09-15 13:00:00",
		},
		{
			SecurityID:      "1333",
			ExchangeSegment: "NSE_EQ",
			Instrument:      "EQUITY",
			Interval:        "1",
			OI:              false,
			FromDate:        "2024-09-11 09:30:00",
			ToDate:          "2024-09-15 13:00:00",
		},
		{
			SecurityID:      "3045",
			ExchangeSegment: "NSE_EQ",
			Instrument:      "EQUITY",
			Interval:        "1",
			OI:              false,
			FromDate:        "2024-09-11 09:30:00",
			ToDate:          "2024-09-15 13:00:00",
		},
		{
			SecurityID:      "1594",
			ExchangeSegment: "NSE_EQ",
			Instrument:      "EQUITY",
			Interval:        "1",
			OI:              false,
			FromDate:        "2024-09-11 09:30:00",
			ToDate:          "2024-09-15 13:00:00",
		},
		{
			SecurityID:      "4963",
			ExchangeSegment: "NSE_EQ",
			Instrument:      "EQUITY",
			Interval:        "1",
			OI:              false,
			FromDate:        "2024-09-11 09:30:00",
			ToDate:          "2024-09-15 13:00:00",
		},
	}

	return &DataMigrator{
		securities:  securities,
		apiURL:      "https://api.dhan.co/v2/charts/intraday",
		accessToken: accessToken,
	}
}

// Run executes the migration process for all securities
func (m *DataMigrator) Run() error {
	for i, security := range m.securities {
		fmt.Printf("Processing security %d of %d: %s (%s)\n", i+1, len(m.securities), security.SecurityID, security.ExchangeSegment)
		err := m.processSecurity(security)
		if err != nil {
			return fmt.Errorf("error processing security %s: %w", security.SecurityID, err)
		}
	}
	return nil
}

// processSecurity handles the data fetching and file writing for a single security
func (m *DataMigrator) processSecurity(security SecurityData) error {
	// Start from January 2020
	startDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Now()

	fmt.Printf("Processing security %s from %s to %s\n", security.SecurityID, startDate.Format("2006-01-02"), currentDate.Format("2006-01-02"))

	var yearData [][]string
	currentYear := startDate.Year()

	// Create header row
	headers := []string{
		"exchange_security_id", "open", "high", "low", "close",
		"volume", "timestamp", "created_at", "updated_at",
	}
	yearData = append(yearData, headers)

	// Process data in 2-month chunks
	processDate := startDate
	for processDate.Before(currentDate) {
		endDate := processDate.AddDate(0, 2, 0) // Add 2 months
		if endDate.After(currentDate) {
			endDate = currentDate
		}

		// If we crossed a year boundary, write the accumulated data and reset
		if processDate.Year() > currentYear {
			err := m.writeCSV(security.SecurityID, currentYear, yearData)
			if err != nil {
				return fmt.Errorf("error writing CSV for year %d: %w", currentYear, err)
			}

			// Reset for the new year
			currentYear = processDate.Year()
			yearData = [][]string{headers}
		}

		fmt.Printf("Fetching data for period: %s to %s\n", processDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

		// Update dates in the security data
		securityCopy := security
		securityCopy.FromDate = processDate.Format("2006-01-02 15:04:05")
		securityCopy.ToDate = endDate.Format("2006-01-02 15:04:05")

		// Fetch data for this period
		data, err := m.fetchData(securityCopy)
		if err != nil {
			return fmt.Errorf("error fetching data: %w", err)
		}

		// Check if we got any data
		if len(data.Timestamp) == 0 {
			fmt.Printf("Warning: No data returned for period %s to %s\n", processDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		} else {
			// Process the response data
			chunkData := m.processResponseData(security.SecurityID, data)
			yearData = append(yearData, chunkData...)
			fmt.Printf("Received %d records for period\n", len(data.Timestamp))
		}

		// Move to the next period
		processDate = endDate.AddDate(0, 0, 1) // Start from the next day

		// Sleep to avoid overwhelming the API
		time.Sleep(1 * time.Second)
	}

	// Write any remaining data for the current year
	if len(yearData) > 1 { // If we have more than just headers
		err := m.writeCSV(security.SecurityID, currentYear, yearData)
		if err != nil {
			return fmt.Errorf("error writing final CSV for year %d: %w", currentYear, err)
		}
	}

	return nil
}

// fetchData makes the API request to fetch chart data
func (m *DataMigrator) fetchData(security SecurityData) (*ChartResponse, error) {
	maxRetries := 3
	retryDelay := 2 * time.Second

	var (
		chartResp *ChartResponse
		err       error
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		chartResp, err = m.doFetchRequest(security)
		if err == nil {
			return chartResp, nil
		}

		fmt.Printf("Attempt %d/%d failed: %v\n", attempt, maxRetries, err)

		if attempt < maxRetries {
			fmt.Printf("Retrying in %v...\n", retryDelay)
			time.Sleep(retryDelay)
			// Increase delay for next attempt
			retryDelay *= 2
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
}

// doFetchRequest performs the actual API request
func (m *DataMigrator) doFetchRequest(security SecurityData) (*ChartResponse, error) {
	reqBody, err := json.Marshal(security)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", m.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access-token", m.accessToken)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var chartResp ChartResponse
	err = json.NewDecoder(resp.Body).Decode(&chartResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &chartResp, nil
}

// processResponseData converts the API response into CSV rows
func (m *DataMigrator) processResponseData(exchangeSegment string, resp *ChartResponse) [][]string {
	var rows [][]string
	now := time.Now().Format("2006-01-02 15:04:05")

	for i := 0; i < len(resp.Timestamp); i++ {
		// Convert timestamp to formatted date
		t := time.Unix(int64(resp.Timestamp[i]), 0).Format("2006-01-02 15:04:05")

		// Create a row with all the data
		row := []string{
			exchangeSegment,
			fmt.Sprintf("%.2f", resp.Open[i]),
			fmt.Sprintf("%.2f", resp.High[i]),
			fmt.Sprintf("%.2f", resp.Low[i]),
			fmt.Sprintf("%.2f", resp.Close[i]),
			fmt.Sprintf("%.2f", resp.Volume[i]),
			t,
			now, // created_at
			now, // updated_at
		}

		rows = append(rows, row)
	}

	return rows
}

// writeCSV writes the accumulated data to a CSV file
func (m *DataMigrator) writeCSV(securityID string, year int, data [][]string) error {
	// Create directory if it doesn't exist
	outputDir := "data"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("instrument_data_%s_%d.csv", securityID, year))
	fmt.Printf("Writing data to file: %s\n", filename)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing row to CSV: %w", err)
		}
	}

	fmt.Printf("Successfully wrote %d records to %s\n", len(data)-1, filename)
	return nil
}
