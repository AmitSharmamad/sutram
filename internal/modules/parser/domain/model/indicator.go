package parser

import (
	instrument "algo/internal/modules/instrument/domain/model"
	"time"
)

type Indicator struct {
	Type    string         `json:"type"`  /// Indicator / Number
	Name    string         `json:"name"`  /// SMA / EMA / RSI
	Point   string         `json:"point"` /// OHLC
	Value   map[string]any `json:"value"`
	Returns string         `json:"returns"`
}

type IndicatorValues struct {
	Timestamp time.Time
	Value     float64
}

type Prices []float64

type IndicatorCalculator interface {
	GetLookBack(config []any) int64
	GenerateIndicatorKey(config []any, ctx *EvaluationContext) (string, error)
	Calculate(config []any, ctx *EvaluationContext) float64
	CalculateAll(config []any, ctx *EvaluationContext, startDate time.Time) []IndicatorValues
	SliceData(prices Prices) Prices
	GetDataRange(config map[string]any) int32
}

func (indicator *Indicator) GetDataRange() int32 {
	calculator := GetIndFactory(indicator.Name)
	return calculator.GetDataRange(indicator.Value)
}

func GetPrices(data instrument.StockData, point string) Prices {
	var results Prices

	for _, d := range data {
		var value float64
		switch point {
		case "OPEN":
			value = d.Open
		case "HIGH":
			value = d.High
		case "CLOSE":
			value = d.Close
		case "LOW":
			value = d.Low
		case "VOLUME":
			value = float64(d.Volume)
		default:
			value = d.Close
		}
		results = append(results, value)
	}

	return results
}

// getPreviousValueFromResults retrieves a previous value from the Results map
func getPreviousValueFromResults(ctx *EvaluationContext, indicatorName string, timestamp int64) (interface{}, bool) {
	if ctx.Results == nil {
		return nil, false
	}

	if timestampResults, exists := ctx.Results[timestamp]; exists {
		if value, found := timestampResults[indicatorName]; found {
			return value, true
		}
	}

	return nil, false
}

// insertValueToResults stores a calculated value in the Results map
func insertValueToResults(ctx *EvaluationContext, indicatorName string, timestamp int64, value interface{}) {
	if ctx.Results == nil {
		ctx.Results = make(map[int64]map[string]interface{})
	}

	if _, exists := ctx.Results[timestamp]; !exists {
		ctx.Results[timestamp] = make(map[string]interface{})
	}

	ctx.Results[timestamp][indicatorName] = value
}
