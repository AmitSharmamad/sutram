package parser

type Low struct {
	offset int
}

func (indicator *Low) SliceData(prices Prices) Prices {
	if indicator.offset == 0 {
		return prices[len(prices)-1:]
	}
	return prices[len(prices)-1-indicator.offset : len(prices)-indicator.offset]
}

func (indicator *Low) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	offset := config[0].(int)
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "LOW")

	if len(prices) <= offset {
		return 0.0
	}

	return prices[len(prices)-1-offset]
}

func (indicator *Low) GetDataRange(config map[string]any) int32 {
	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}
	return int32(offset + 1)
}
