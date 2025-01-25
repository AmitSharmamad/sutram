package algo

import (
	"algo/src/algo/indicator"
	"algo/src/stock"
)

type ILogicalExpression interface {
	Handle(data stock.StockData)
}

type LogicalExpression struct {
	Type     string                         `json:"type"`     /// logicalexpression
	Operator string                         `json:"operator"` /// AND / OR / NA
	Returns  string                         `json:"string"`
	Lhs      *indicator.IndicatorExpression `json:"lhs"`
	Rhs      *indicator.IndicatorExpression `json:"rhs"`
}

func (expression *LogicalExpression) Handle(data stock.StockData) bool {
	println("Handling logical expression")

	lhsValue := expression.Lhs.Handle(data)

	if expression.Operator == "NA" {
		return lhsValue
	}

	rhsValue := expression.Rhs.Handle(data)

	if expression.Operator == "AND" {
		return lhsValue && rhsValue
	} else if expression.Operator == "OR" {
		return lhsValue || rhsValue
	}

	return false
}
