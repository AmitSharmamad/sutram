package parser

type Open struct {
	offset int
}

func (indicator *Open) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

func (indicator *Open) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := int(config[0].(float64))
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "OPEN")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

func (indicator *Open) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}

// GetLookBack returns the number of historical data points needed for Open price calculation
// config[0] = offset (float64)
func (indicator *Open) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	offset := int(config[0].(float64)) // default: 0

	// For Open price calculation, we need:
	// - 1 data point for current open
	// - Plus any offset to access historical open values
	return int64(offset + 1)
}
