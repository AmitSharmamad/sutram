package parser

import (
	"fmt"
	"time"
)

type Open struct {
	offset int
}

func (indicator *Open) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

// GenerateIndicatorKey generates a unique key for Open indicator based on its configuration
func (indicator *Open) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 1 {
		return "", fmt.Errorf("Open requires 1 parameter: offset")
	}

	// Extract offset parameter
	offsetFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("Open offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	if offset != 0 {
		return fmt.Sprintf("OPEN_%d", offset), nil
	}

	return "OPEN", nil
}

func (indicator *Open) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := int(config[0].(float64))
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "OPEN")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

func (indicator *Open) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}

// GetLookBack returns the number of historical data points needed for Open price calculation
// config[0] = offset (float64)
func (indicator *Open) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	offset := int(config[0].(float64)) // default: 0

	// For Open price calculation, we need:
	// - 1 data point for current open
	// - Plus any offset to access historical open values
	return int64(offset + 1)
}

// CalculateAll implements the IndicatorCalculator interface
func (indicator *Open) CalculateAll(
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
			Value:     data.Open,
		})
	}

	return result
}
