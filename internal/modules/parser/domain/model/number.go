package parser

type IndNumber struct {
	Value float64
}

func (indicator *IndNumber) SliceData(prices Prices) Prices {
	return []float64{}
}

func (indicator *IndNumber) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	return config[0].(float64)
}

func (number *IndNumber) GetDataRange(config map[string]any) int32 {
	return 0
}

// GetLookBack returns the number of historical data points needed for Number calculation
// Number indicator doesn't need any historical data as it returns a constant value
func (number *IndNumber) GetLookBack(config []any) int64 {
	// Number indicator returns a constant value and doesn't need any historical data
	return 0
}
