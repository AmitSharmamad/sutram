package indicator

import "algo/src/data"

type EMA struct {
	period int
	offset int
}

func (indicator *EMA) Calculate(data *[]data.Data, config map[string]any, point string) float64 {
	println("Executing EMA")

	println(data)
	println(config)
	println(point)

	return 0.0
}
