package indicator

import "algo/src/stock"

type SMA struct {
	period int
	offset int
}

func NewSMA() *SMA {
	return &SMA{}
}

func (indicator *SMA) Calculate(data stock.StockData, config map[string]any, point string) float64 {
	println("Executing SMA")

	println(data)
	println(config)
	println(point)

	return 0.0
}
