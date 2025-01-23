package indicator

import (
	"algo/src/algo/operator"
	"algo/src/data"
)

type IndicatorCalculator interface {
	Calculate(data *[]data.Data, config map[string]any, point string) float64
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

func (expression *IndicatorExpression) Handle(data *[]data.Data) bool {
	println("Handling indicator expression")

	lhsValue := expression.Lhs.handle(data)
	rhsValue := expression.Rhs.handle(data)

	return operator.Get(expression.Operator).Operate(
		lhsValue, rhsValue,
	)
}

func (indicator *Indicator) handle(data *[]data.Data) float64 {

	println("Handling indicator", indicator)

	config := make(map[string]any, 0)

	calculator := Get(indicator.Name)

	return (calculator).Calculate(data, config, indicator.Point)
}
