package parser

type Volume struct {
	offset int
}

func (indicator *Volume) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

func (indicator *Volume) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := config[0].(int)
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "VOLUME")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

func (indicator *Volume) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}

// GetLookBack returns the number of historical data points needed for Volume calculation
// config[0] = offset (int)
func (indicator *Volume) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	offset := config[0].(int) // default: 0

	// For Volume calculation, we need:
	// - 1 data point for current volume
	// - Plus any offset to access historical volume values
	return int64(offset + 1)
}
