package main

import (
	"algo/src/algo"
	"algo/src/algo/indicator"
	"algo/src/algo/operator"
	"algo/src/data"
	"algo/src/util"
)

func main() {

	Init()

	exp := util.GetFileContent("data/many-indicators.json")

	util.Validate(exp)

	expression := algo.NewExpressionStatement(exp)

	data := make([]data.Data, 0)

	expression.Eval(&data)

	println()
}

func Init() {
	indicator.Init()
	operator.Init()
}
