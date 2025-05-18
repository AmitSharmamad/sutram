package parser

type Close struct {
	offset int
}

// SliceData returns the closing prices based on offset
func (indicator *Close) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

// Calculate returns the closing price at the specified offset
// Close(offset)
func (indicator *Close) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := int(config[0].(float64))
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "CLOSE")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

// GetDataRange returns the required number of data points
func (indicator *Close) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}