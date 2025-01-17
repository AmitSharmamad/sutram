package indicator

import "algo/src/data"

type IndNumber struct {
	Value float64
}

func (indicator *IndNumber) Calculate(data *[]data.Data, config map[string]any, point string) float64 {
	println("Executing Number")

	println(data)
	println(config)
	println(point)

	return 0.0
}
