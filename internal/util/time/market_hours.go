package time

import (
	"time"

	strategyDto "algo/internal/modules/strategy/application/dto"
)

// MarketHoursUtil provides utility functions for market hours and time window validation
type MarketHoursUtil struct{}

// NewMarketHoursUtil creates a new instance of MarketHoursUtil
func NewMarketHoursUtil() *MarketHoursUtil {
	return &MarketHoursUtil{}
}

// GetDefaultMarketHours returns the default trading hours for different exchanges
func (m *MarketHoursUtil) GetDefaultMarketHours(exchange string) (startHour, startMinute, endHour, endMinute int) {
	switch exchange {
	case "NSE", "BSE":
		// Indian equity markets: 9:15 AM to 3:30 PM IST
		return 9, 15, 15, 30
	case "MCX":
		// MCX commodity market: 9:00 AM to 11:30 PM IST
		return 9, 0, 23, 30
	case "NCDEX":
		// NCDEX commodity market: 10:00 AM to 5:00 PM IST
		return 10, 0, 17, 0
	default:
		// Default to NSE timings
		return 9, 15, 15, 30
	}
}

// IsWithinEntryTimeWindow validates if the current time is within the entry time window
func (m *MarketHoursUtil) IsWithinEntryTimeWindow(
	currentTimestamp time.Time,
	strategy *strategyDto.StrategyResponse,
	exchange string,
) bool {
	// Extract time of day from the current timestamp
	currentTimeOfDay := time.Date(0, 1, 1,
		currentTimestamp.Hour(), currentTimestamp.Minute(), currentTimestamp.Second(), 0, time.UTC)

	var entryStartTime, entryEndTime time.Time

	// Determine entry start time
	if strategy.EntryStartTime != nil {
		// Use provided entry start time
		entryStartTime = time.Date(0, 1, 1,
			strategy.EntryStartTime.Hour, strategy.EntryStartTime.Minute,
			strategy.EntryStartTime.Second, 0, time.UTC)
	} else {
		// Use default market start time if not specified
		startHour, startMinute, _, _ := m.GetDefaultMarketHours(exchange)
		entryStartTime = time.Date(0, 1, 1, startHour, startMinute, 0, 0, time.UTC)
	}

	// Determine entry end time
	if strategy.EntryEndTime != nil {
		// Use provided entry end time
		entryEndTime = time.Date(0, 1, 1,
			strategy.EntryEndTime.Hour, strategy.EntryEndTime.Minute,
			strategy.EntryEndTime.Second, 0, time.UTC)
	} else {
		// Use default market end time if not specified
		_, _, endHour, endMinute := m.GetDefaultMarketHours(exchange)
		entryEndTime = time.Date(0, 1, 1, endHour, endMinute, 0, 0, time.UTC)
	}

	// Check if current time is within the entry window
	// Current time must be >= start time AND <= end time
	return (currentTimeOfDay.Equal(entryStartTime) || currentTimeOfDay.After(entryStartTime)) &&
		(currentTimeOfDay.Equal(entryEndTime) || currentTimeOfDay.Before(entryEndTime))
}

// IsMISEndOfDayHit checks if the current time is at or after the MIS end-of-day cutoff time
func (m *MarketHoursUtil) IsMISEndOfDayHit(timestamp time.Time, exchange string) bool {
	// NOTE: MIS close is 15 minutes prior to market close
	endHour := 15
	endMinute := 15

	if exchange == "MCX" {
		// MCX has different timings, so we need to adjust accordingly
		// MCX trading ends at 11:30 PM IST
		endHour = 23
		endMinute = 30
	}

	// Extract the time from the provided timestamp
	currentTime := timestamp

	// Create a time representing the MIS end-of-day cutoff time on the same date
	eodTime := time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		endHour,   // End hour
		endMinute, // End minute
		0,         // 0 seconds
		0,         // 0 nanoseconds
		currentTime.Location(),
	)

	// Return true if current time is at or after the MIS end-of-day cutoff time
	return currentTime.Equal(eodTime) || currentTime.After(eodTime)
}
