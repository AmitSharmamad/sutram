package parser

import (
	"fmt"
	"time"
)

type Volume struct {
	offset int
}

func (indicator *Volume) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

// generateIndicatorKey generates a unique key for Volume indicator based on its configuration
func (indicator *Volume) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 1 {
		return "", fmt.Errorf("Volume requires 1 parameter: offset")
	}

	// Extract offset parameter
	offsetFloat, ok := config[0].(int)
	if !ok {
		// Try float64 conversion as fallback
		if offsetFloatAlt, ok := config[0].(float64); ok {
			offsetFloat = int(offsetFloatAlt)
		} else {
			return "", fmt.Errorf("Volume offset parameter must be a number")
		}
	}
	offset := offsetFloat

	// Generate the key - include offset only if non-zero
	if offset != 0 {
		return fmt.Sprintf("VOLUME_%d", offset), nil
	}

	return "VOLUME", nil
}

func (indicator *Volume) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := config[0].(int)
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "VOLUME")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

func (indicator *Volume) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}

// GetLookBack returns the number of historical data points needed for Volume calculation
// config[0] = offset (int)
func (indicator *Volume) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	offset := config[0].(float64) // default: 0

	// For Volume calculation, we need:
	// - 1 data point for current volume
	// - Plus any offset to access historical volume values
	return int64(offset + 1)
}

// CalculateAll implements the IndicatorCalculator interface
func (indicator *Volume) CalculateAll(
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
			Value:     float64(data.Volume),
		})
	}

	return result
}
