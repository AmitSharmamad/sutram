package algo

import (
	"algo/src/data"
	"encoding/json"
)

type ExpressionStatement struct {
	Type  string             `json:"type"` /// ExpressionStatement
	Value *LogicalExpression `json:"value"`
}

func NewExpressionStatement(expression string) *ExpressionStatement {
	var input *ExpressionStatement = &ExpressionStatement{}
	err := json.Unmarshal([]byte(expression), input)
	if err != nil {
		panic(err)
	}
	return input
}

func (expression *ExpressionStatement) Eval(data *[]data.Data) {
	println("Evaluating input ", expression)
	if expression.Type == "ExpressionStatement" {
		if expression.Value.Type == "LogicalExpression" {
			expression.Value.Handle(data)
			return
		}
	}
	panic("Incorrect input for expression statement")
}
