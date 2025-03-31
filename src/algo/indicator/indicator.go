package indicator

import (
	"algo/src/algo/operator"
	"algo/src/shared"
	"algo/src/util"
)

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
	Calculate(prices Prices, config []any, point string) float64
	CalculateAll(data shared.StockData, config map[string]any, point string) []float64
	GetDataRange(config map[string]any) int32
}

func (expression *IndicatorExpression) Handle(data shared.StockData) bool {

	lhsValue := expression.Lhs.Handle(data)
	rhsValue := expression.Rhs.Handle(data)

	return operator.GetOperator(expression.Operator).Operate(
		lhsValue, rhsValue,
	)
}

func (expression *IndicatorExpression) GetDataRange() int32 {
	lhsValue := expression.Lhs.GetDataRange()
	rhsValue := expression.Rhs.GetDataRange()
	return util.Max(lhsValue, rhsValue)
}

func (indicator *Indicator) Handle(data shared.StockData) float64 {

	calculator := GetIndFactory(indicator.Name)

	prices := GetPrices(data, indicator.Point)
	config := []any{}

	return (calculator).Calculate(prices, config, indicator.Point)
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
