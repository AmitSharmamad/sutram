package shared

type ILogicalExpression interface {
	Handle(data StockData)
}

type LogicalExpression struct {
	Type     string               `json:"type"`     /// logicalexpression
	Operator string               `json:"operator"` /// AND / OR / NA
	Returns  string               `json:"string"`
	Lhs      *IndicatorExpression `json:"lhs"`
	Rhs      *IndicatorExpression `json:"rhs"`
}

func (expression *LogicalExpression) Handle(data StockData) bool {
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
