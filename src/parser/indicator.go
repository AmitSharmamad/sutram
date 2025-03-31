package parser

import (
	"algo/src/shared"
)

type Indicator struct {
	Type    string         `json:"type"`  /// Indicator / Number
	Name    string         `json:"name"`  /// SMA / EMA / RSI
	Point   string         `json:"point"` /// OHLC
	Value   map[string]any `json:"value"`
	Returns string         `json:"returns"`
}

type Prices []float64

type IndicatorCalculator interface {
	Calculate(prices Prices, config []any, point string) float64
	CalculateAll(data shared.StockData, config map[string]any, point string) []float64
	GetDataRange(config map[string]any) int32
}

func (indicator *Indicator) GetDataRange() int32 {
	calculator := GetIndFactory(indicator.Name)
	return calculator.GetDataRange(indicator.Value)
}

func GetPrices(data shared.StockData, point string) Prices {
	var results Prices

	for _, d := range data {
		var value float64
		switch point {
		case "open":
			value = d.Open
		case "high":
			value = d.High
		case "close":
			value = d.Close
		case "low":
			value = d.Low
		case "volume":
			value = float64(d.Volume)
		default:
			value = d.Close
		}
		results = append(results, value)
	}

	return results
}
