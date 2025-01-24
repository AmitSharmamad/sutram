package main

import (
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

func init() {
}

func main() {

	Print()

	// di.InitialiseDI()

	/*
		println("running ...")

		Print()

		path := "../data/llm.json"

		exp := util.GetFileContent(path)

		util.Validate(exp)

		expression := algo.NewExpressionStatement(exp)

		data := make([]data.Data, 0)

		expression.Eval(&data)

		println()
	*/
}
