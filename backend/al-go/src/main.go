package main

import (
	"algo/src/algo"
	"algo/src/algo/indicator"
	"algo/src/algo/operator"
	"algo/src/data"
	"algo/src/util"
	"fmt"
	"os"
)

func Print() {
	if len(os.Args) > 1 {
		fmt.Println("Arguments provided:")
		for i, arg := range os.Args[1:] {
			fmt.Printf("Arg %d: %s\n", i+1, arg)
		}
	} else {
		fmt.Println("No arguments provided.")
	}
}

func main() {

	Print()

	Init()

	path := "../data/many-indicators.json"

	exp := util.GetFileContent(path)

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
