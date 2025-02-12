package expression

import (
	"algo/src/algo/indicator"
	"algo/src/shared"
	"algo/src/util"
)

type ILogicalExpression interface {
	Handle(data shared.StockData)
}

type LogicalExpression struct {
	Type     string                         `json:"type"`     /// logicalexpression
	Operator string                         `json:"operator"` /// AND / OR / NA
	Returns  string                         `json:"string"`
	Lhs      *indicator.IndicatorExpression `json:"lhs"`
	Rhs      *indicator.IndicatorExpression `json:"rhs"`
}

func (expression *LogicalExpression) Handle(data shared.StockData) bool {
	println("Handling logical expression")

	lhsValue := expression.Lhs.Handle(data)

	if expression.Operator == "NA" {
		return lhsValue
	}

	rhsValue := expression.Rhs.Handle(data)

	switch expression.Operator {
	case "AND":
		return lhsValue && rhsValue
	case "OR":
		return lhsValue || rhsValue
	default:
		panic("invalid operator")
	}
}

func (expression *LogicalExpression) GetDataRange() int32 {
	lhsValue := expression.Lhs.GetDataRange()

	if expression.Operator == "NA" {
		return lhsValue
	}

	rhsValue := expression.Rhs.GetDataRange()

	return util.Max(lhsValue, rhsValue)
}
