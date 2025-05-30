package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"
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

	// Calculate SMA
	var sum float64
	start := len(prices) - period - offset
	end := len(prices) - offset

	for i := start; i < end; i++ {
		sum += prices[i]
	}

	return util.Round(sum / float64(period))
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

	var filteredResults []IndicatorValues

	// Calculate SMA for each data point where we have enough historical data
	for i := period + offset; i < len(stockData); i++ {
		// Calculate SMA for current position
		var sum float64
		start := i - period - offset
		end := i - offset

		for j := start; j < end; j++ {
			sum += prices[j]
		}

		smaValue := util.Round(sum / float64(period))

		// Only include values from startDate onwards
		if !stockData[i].Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[i].Timestamp,
				Value:     smaValue,
			})
		}
	}

	values := append(createIndicatorValues(period+offset), filteredResults...)

	return values
}
