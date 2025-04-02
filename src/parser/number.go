package parser

type IndNumber struct {
	Value float64
}

func (indicator *IndNumber) SliceData(prices Prices) Prices {
	return []float64{}
}

func (indicator *IndNumber) Calculate(
	config []any,
	prices Prices,
) float64 {
	return config[0].(float64)
}

func (number *IndNumber) GetDataRange(config map[string]any) int32 {
	return 0
}
