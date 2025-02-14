package shared

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
