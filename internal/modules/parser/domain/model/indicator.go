package parser

import instrument "algo/internal/modules/instrument/domain/model"

type Indicator struct {
	Type    string         `json:"type"`  /// Indicator / Number
	Name    string         `json:"name"`  /// SMA / EMA / RSI
	Point   string         `json:"point"` /// OHLC
	Value   map[string]any `json:"value"`
	Returns string         `json:"returns"`
}

type Prices []float64

type IndicatorCalculator interface {
	Calculate(config []any, ctx *EvaluationContext) float64
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