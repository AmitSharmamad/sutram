package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var file = "sql/ohlcv_data.csv"
var count = 1_000_000

func main() {
	// Generate and print OHLCV data to CSV
	if err := generateOHLCVDataToCSV(file, count, 1678886400, 1778886400); err != nil {
		panic(err)
	}

	fmt.Println("OHLCV data generated and saved to ohlcv_data.csv")
}

func generateOHLCVDataToCSV(filename string, count int, startTime int64, endTime int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"timestamp", "open", "high", "low", "close", "volume", "created_at", "updated_at"}); err != nil {
		return err
	}

	currentTime := startTime
	for i := 0; i < count; i++ {
		open := rand.Float64() * 1000
		high := open + rand.Float64()*10
		low := open - rand.Float64()*10
		close := low + rand.Float64()*(high-low)
		volume := uint64(rand.Intn(1000000)) // Generate random volume as UInt64

		// Format timestamps without time zone information
		timestamp := time.Unix(currentTime, 0).UTC().Format("2006-01-02 15:04:05")
		now := time.Now().UTC().Format("2006-01-02 15:04:05")

		record := []string{
			timestamp,
			fmt.Sprintf("%.2f", open),
			fmt.Sprintf("%.2f", high),
			fmt.Sprintf("%.2f", low),
			fmt.Sprintf("%.2f", close),
			fmt.Sprintf("%d", volume),
			now,
			now,
		}

		if err := writer.Write(record); err != nil {
			return err
		}

		currentTime += 60
		if currentTime > endTime {
			currentTime = startTime
		}
	}

	return nil
}
