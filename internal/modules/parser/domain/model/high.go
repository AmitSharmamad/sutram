package parser

import "time"

type High struct {
	offset int
}

func (indicator *High) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

func (indicator *High) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := int(config[0].(float64))
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "HIGH")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

func (indicator *High) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}

// GetLookBack returns the number of historical data points needed for High price calculation
// config[0] = offset (float64)
func (indicator *High) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	offset := int(config[0].(float64)) // default: 0

	// For High price calculation, we need:
	// - 1 data point for current high
	// - Plus any offset to access historical high values
	return int64(offset + 1)
}

// CalculateAll implements the IndicatorCalculator interface
func (indicator *High) CalculateAll(
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
			Value:     data.High,
		})
	}

	return result
}
