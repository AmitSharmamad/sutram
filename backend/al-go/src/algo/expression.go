package algo

import (
	"encoding/json"

	"algo/src/shared"
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

func (expression *Expression) Eval(data shared.StockData) {
	println("Evaluating input ", expression)
	if expression.Type == "ExpressionStatement" {
		if expression.Value.Type == "LogicalExpression" {
			expression.Value.Handle(data)
			return
		}
	}
	panic("Incorrect input for expression statement")
}
