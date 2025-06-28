package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
)

type EMA struct {
	period int
	offset int
}

func (indicator *EMA) SliceData(prices Prices) Prices {
	// We need enough data to calculate EMA properly
	// For EMA calculation we need at least 2*period data points
	// This ensures we have enough data for initial SMA and subsequent EMA calculation
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for EMA indicator based on its configuration
func (indicator *EMA) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 3 {
		return "", fmt.Errorf("EMA requires 3 parameters: field, period, offset")
	}

	// Extract field parameter
	field, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("EMA field parameter must be a string")
	}

	// Extract period parameter
	periodFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("EMA period parameter must be a number")
	}
	period := int(periodFloat)

	// Extract offset parameter (though typically not used in the key)
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("EMA offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - typically we don't include offset in the key
	// as it's used for historical lookback, not for identifying the indicator
	// But if you want to include offset for different EMA calculations with same period:
	if offset != 0 {
		return fmt.Sprintf("EMA_%s_%d_%d", strings.ToUpper(field), period, offset), nil
	}

	return fmt.Sprintf("EMA_%s_%d", strings.ToUpper(field), period), nil
}

func (indicator *EMA) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	field := config[0].(string)
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, field)

	// Check if we have enough data points
	if len(prices) < period+offset {
		return 0.0
	}

	prices = indicator.SliceData(prices)

	// Need at least period data points to calculate
	if len(prices) < period {
		return 0.0
	}

	// Create unique indicator key
	indicatorKey := fmt.Sprintf("EMA_%s_%d", field, period)

	// Get current timestamp (using the last data point's timestamp)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()

	// Use go-talib EMA calculation
	emaValues := talib.Ema(prices, period)

	if len(emaValues) == 0 {
		return 0.0
	}

	// Get the last EMA value
	ema := emaValues[len(emaValues)-1]

	// Store the result for future calculations
	insertValueToResults(ctx, indicatorKey, currentTimestamp, ema)

	return util.Round(ema)
}

func (indicator *EMA) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 20 // default period from indicators.json
	}

	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.period + indicator.offset)
}

func (indicator *EMA) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is field (string)
	// config[1] is period (float64)
	// config[2] is offset (float64)

	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	// For EMA calculation, we need sufficient historical data
	// Using 2*period to ensure we have enough data for proper EMA calculation
	// as used in SliceData method, plus the offset
	return int64(5*period + offset)
}

// CalculateAll calculates EMA values for all data points
func (indicator *EMA) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	field := config[0].(string)
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.offset = offset

	stockData := *ctx.Data

	// Get prices for the entire dataset
	prices := GetPrices(stockData, field)

	// We need at least enough data to calculate initial EMA
	minDataPoints := period + offset
	if len(stockData) < minDataPoints {
		return []IndicatorValues{}
	}

	// Use go-talib EMA calculation
	emaValues := talib.Ema(prices, period)

	if len(emaValues) == 0 {
		return []IndicatorValues{}
	}

	var filteredResults []IndicatorValues

	// Calculate the starting index accounting for the NaN values that talib returns
	// talib.Ema returns NaN for the first (period-1) values
	startIndex := period - 1 + offset

	// Process each EMA value and match it with the corresponding timestamp
	for i, emaValue := range emaValues {
		dataIndex := i + offset

		// Skip if we don't have enough data or if the value is NaN
		if dataIndex >= len(stockData) || emaValue != emaValue { // NaN check
			continue
		}

		// Only include values from startDate onwards
		if !stockData[dataIndex].Timestamp.Before(startDate) && dataIndex >= startIndex {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     util.Round(emaValue),
			})
		}
	}

	return filteredResults
}
