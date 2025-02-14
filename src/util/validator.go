package util

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func Validate(expression string) bool {
	data := GetFileContent("../data/schema.json")

	schema := gojsonschema.NewStringLoader(string(data))

	println(schema)

	value := gojsonschema.NewStringLoader(expression)

	result, err := gojsonschema.Validate(schema, value)
	if err != nil {
		panic(err)
	}

	if !result.Valid() {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			// fmt.Printf("- %s\n", desc)
			fmt.Printf("  - %s\n", desc)
		}
	}

	return result.Valid()
}
