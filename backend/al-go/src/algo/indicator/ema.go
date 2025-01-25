package indicator

import "algo/src/stock"

type EMA struct {
	period int
	offset int
}

func (indicator *EMA) Calculate(data stock.StockData, config map[string]any, point string) float64 {
	println("Executing EMA")

	println(data)
	println(config)
	println(point)

	return 0.0
}
