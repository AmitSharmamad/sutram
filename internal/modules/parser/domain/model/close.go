package parser

import (
	"fmt"
	"time"
)

type Close struct {
	offset int
}

// SliceData returns the closing prices based on offset
func (indicator *Close) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

// GenerateIndicatorKey generates a unique key for Close indicator based on its configuration
func (indicator *Close) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 1 {
		return "", fmt.Errorf("Close requires 1 parameter: offset")
	}

	// Extract offset parameter
	offsetFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("Close offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	if offset != 0 {
		return fmt.Sprintf("CLOSE_%d", offset), nil
	}

	return "CLOSE", nil
}

// Calculate returns the closing price at the specified offset
// Close(offset)
func (indicator *Close) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := int(config[0].(float64))
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "CLOSE")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

// CalculateAll implements the IndicatorCalculator interface
func (indicator *Close) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	offset := int(config[0].(float64))
	indicator.offset = offset

	// Get all stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	var result []IndicatorValues

	// Filter data starting from startDate
	for _, data := range stockData {
		if data.Timestamp.Before(startDate) {
			continue
		}

		// Calculate indicator value for this data point
		result = append(result, IndicatorValues{
			Timestamp: data.Timestamp,
			Value:     data.Close,
		})
	}

	return result
}

// GetDataRange returns the required number of data points
func (indicator *Close) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}

// GetLookBack returns the number of historical data points needed for Close price calculation
// config[0] = offset (float64)
func (indicator *Close) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	offset := int(config[0].(float64)) // default: 0

	// For Close price calculation, we need:
	// - 1 data point for current close
	// - Plus any offset to access historical close values
	return int64(offset + 1)
}
