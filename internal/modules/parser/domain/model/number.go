package parser

import "time"

type IndNumber struct {
	Value float64
}

func (indicator *IndNumber) SliceData(prices Prices) Prices {
	return []float64{}
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
