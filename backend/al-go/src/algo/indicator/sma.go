package indicator

import "algo/src/data"

type SMA struct {
	period int
	offset int
}

func (indicator *SMA) Calculate(data *[]data.Data, config map[string]any, point string) float64 {
	println("Executing SMA")

	println(data)
	println(config)
	println(point)

	return 0.0
}
