package algo

import (
	"algo/src/data"
	"encoding/json"
)

type IExpression interface {
	Eval(data data.StockData)
}

type Expression struct {
	Type  string             `json:"type"` /// ExpressionStatement
	Value *LogicalExpression `json:"value"`
}

func (*Expression) Init() {}

func NewExpressionStatement(expression string) IExpression {
	var input *Expression = &Expression{}
	err := json.Unmarshal([]byte(expression), input)
	if err != nil {
		panic(err)
	}
	return input
}

func (expression *Expression) Eval(data data.StockData) {
	println("Evaluating input ", expression)
	if expression.Type == "ExpressionStatement" {
		if expression.Value.Type == "LogicalExpression" {
			expression.Value.Handle(data)
			return
		}
	}
	panic("Incorrect input for expression statement")
}
