package shared

import (
	"encoding/json"
)

type Expression struct {
	Type  string             `json:"type"`
	Value *LogicalExpression `json:"value"`
}

func (*Expression) Init() {}

func NewExpressionStatement(expression string) *Expression {
	input := &Expression{}
	err := json.Unmarshal([]byte(expression), input)
	if err != nil {
		panic(err)
	}
	return input
}

func (expression *Expression) Eval(data StockData) bool {
	println("Evaluating input ", expression)
	if expression.Type == "ExpressionStatement" {
		if expression.Value.Type == "LogicalExpression" {
			return expression.Value.Handle(data)
		}
	}
	return false
}
