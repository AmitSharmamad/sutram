package parser

type High struct {
	offset int
}

func (indicator *High) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

func (indicator *High) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := config[0].(int)
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "HIGH")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

func (indicator *High) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}
