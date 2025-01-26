package indicator

import (
	"algo/src/algo/operator"
	"algo/src/stock"
)

type Prices []float64

type IndicatorCalculator interface {
	Calculate(prices Prices, config map[string]any, point string) float64
}

type IndicatorExpression struct {
	Type     string     `json:"type"`     /// IndicatorExpression
	Operator string     `json:"operator"` /// CrossesAbove / CrossesBelow
	Returns  string     `json:"returns"`
	Lhs      *Indicator `json:"lhs"`
	Rhs      *Indicator `json:"rhs"`
}

type Indicator struct {
	Type    string         `json:"type"`  /// Indicator / Number
	Name    string         `json:"name"`  /// SMA / EMA / RSI
	Point   string         `json:"point"` /// OHLC
	Value   map[string]any `json:"value"`
	Returns string         `json:"returns"`
}

type Prices []float64

type IndicatorCalculator interface {
	Calculate(prices Prices, config map[string]any, point string) float64
}

func (expression *IndicatorExpression) Handle(data stock.StockData) bool {
	println("Handling indicator expression")

	lhsValue := expression.Lhs.Handle(data)
	rhsValue := expression.Rhs.Handle(data)

	return operator.Get(expression.Operator).Operate(
		lhsValue, rhsValue,
	)
}

func (indicator *Indicator) handle(data stock.StockData) float64 {
	println("Handling indicator", indicator)

	config := make(map[string]any, 0)

	calculator := Get(indicator.Name)

	prices := GetPrices(data, indicator.Point)

	return (calculator).Calculate(prices, config, indicator.Point)
}

func GetPrices(data stock.StockData, point string) Prices {
	var results Prices

	for _, d := range data {
		var value float64
		switch point {
		case "open":
			value = d.Close
		case "high":
			value = d.Open
		case "close":
			value = d.Close
		case "low":
			value = d.Open
		case "volume":
			value = d.Open
		default:
			value = d.Close
		}
		results = append(results, value)
	}

	return results
}
