package indicator

import "algo/src/shared"

type IndNumber struct {
	Value float64
}

func (indicator *IndNumber) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	println("Executing Number")

	println(prices)
	println(config)
	println(point)

	return 0.0
}

func (number *IndNumber) GetDataRange() int32 {
	return 0
}

func (indicator *IndNumber) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
