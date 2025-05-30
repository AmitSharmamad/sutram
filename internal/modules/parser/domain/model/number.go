package parser

import (
	"fmt"
	"time"
)

type IndNumber struct {
	Value float64
}

func (indicator *IndNumber) SliceData(prices Prices) Prices {
	return []float64{}
}

// GenerateIndicatorKey generates a unique key for IndNumber indicator based on its configuration
func (indicator *IndNumber) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 1 {
		return "", fmt.Errorf("IndNumber requires 1 parameter: value")
	}

	// Extract value parameter
	value, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("IndNumber value parameter must be a number")
	}

	// Generate the key - numbers are constants, so use the value
	return fmt.Sprintf("NUMBER_%.6f", value), nil
}

func (indicator *IndNumber) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	return config[0].(float64)
}

func (indicator *IndNumber) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	value := config[0].(float64)

	// Get all stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	var result []IndicatorValues

	// Return the same value for all data points from startDate onwards
	for _, data := range stockData {
		if !data.Timestamp.Before(startDate) {
			result = append(result, IndicatorValues{
				Timestamp: data.Timestamp,
				Value:     value,
			})
		}
	}

	return result
}

func (number *IndNumber) GetDataRange(config map[string]any) int32 {
	return 0
}

// GetLookBack returns the number of historical data points needed for Number calculation
// Number indicator doesn't need any historical data as it returns a constant value
func (number *IndNumber) GetLookBack(config []any) int64 {
	// Number indicator returns a constant value and doesn't need any historical data
	return 0
}
