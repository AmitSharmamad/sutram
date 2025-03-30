package expression

import (
	"algo/src/shared"
	"encoding/json"
)

type Expression struct {
	Type  string             `json:"type"`
	Value *LogicalExpression `json:"value"`
}

func (*Expression) Init() {}

func NewExpressionStatement(expression string) *Expression {
	expr := &Expression{}
	err := json.Unmarshal([]byte(expression), expr)
	if err != nil {
		panic(err)
	}
	return expr
}

func (expression *Expression) Eval(data shared.StockData) bool {

	if expression.Type == "ExpressionStatement" {
		if expression.Value.Type == "LogicalExpression" {
			return expression.Value.Handle(data)
		}
	}
	return false
}

func (expression *Expression) GetDataRange() int32 {
	if expression.Type == "ExpressionStatement" {
		if expression.Value.Type == "LogicalExpression" {
			return expression.Value.GetDataRange()
		}
	}
	return 0
}
