package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
)

type SMA struct {
	period int
	offset int
}

func (indicator *SMA) SliceData(prices Prices) Prices {
	start := len(prices) - indicator.offset - indicator.period
	end := len(prices) - indicator.offset
	return prices[start:end]
}

// generateIndicatorKey generates a unique key for SMA indicator based on its configuration
func (indicator *SMA) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 3 {
		return "", fmt.Errorf("SMA requires 3 parameters: field, period, offset")
	}

	// Extract field parameter
	field, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("SMA field parameter must be a string")
	}

	// Extract period parameter
	periodFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("SMA period parameter must be a number")
	}
	period := int(periodFloat)

	// Extract offset parameter
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("SMA offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	if offset != 0 {
		return fmt.Sprintf("SMA_%s_%d_%d", strings.ToUpper(field), period, offset), nil
	}

	return fmt.Sprintf("SMA_%s_%d", strings.ToUpper(field), period), nil
}

// SMA(field, period, offset)
func (indicator *SMA) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	field := config[0].(string)
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, field)

	if len(prices) < period+offset {
		return 0.0
	}

	// Use go-talib SMA calculation
	smaValues := talib.Sma(prices, period)

	// Get the value at the correct offset position
	targetIndex := len(smaValues) - 1 - offset
	if targetIndex < 0 || targetIndex >= len(smaValues) {
		return 0.0
	}

	return util.Round(smaValues[targetIndex])
}

func (indicator *SMA) GetDataRange(config map[string]any) int32 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(float64); ok {
			period = int(p)
		} else {
			period = 20 // default period from indicators.json
		}
	}

	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(float64); ok {
			offset = int(o)
		} else {
			offset = 0
		}
	}

	return int32(period + offset)
}

// GetLookBack returns the number of historical data points needed for SMA calculation
// config[0] = field (string), config[1] = period (float64), config[2] = offset (float64)
func (indicator *SMA) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is field (string) - not needed for lookback calculation
	// config[1] is period (float64)
	// config[2] is offset (float64)

	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	// For SMA calculation, we need exactly period + offset data points
	// The SMA needs 'period' consecutive data points plus any offset
	return int64(period + offset)
}

// CalculateAll calculates SMA values for all data points
func (indicator *SMA) CalculateAll(
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

	// Use go-talib SMA calculation for all values
	smaValues := talib.Sma(prices, period)

	var filteredResults []IndicatorValues

	// Process SMA values and apply offset
	for i := 0; i < len(smaValues); i++ {
		// Skip NaN values (initial values where SMA can't be calculated)
		if smaValues[i] != smaValues[i] { // NaN check
			continue
		}

		// Calculate the corresponding data index considering offset
		dataIndex := i + offset
		if dataIndex >= len(stockData) {
			break
		}

		// Only include values from startDate onwards
		if !stockData[dataIndex].Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     smaValues[i],
			})
		}
	}

	// values := append(createIndicatorValues(period+offset), filteredResults...)
	values := filteredResults

	return values
}
