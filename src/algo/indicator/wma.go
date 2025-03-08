package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type WMA struct {
	period int
}

func (indicator *WMA) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 14 // default period
		}
	}

	if len(prices) < period {
		return 0.0
	}

	start := len(prices) - period
	weightedSum := 0.0
	weightSum := 0.0

	// Calculate WMA using weighted sum
	for i := 0; i < period; i++ {
		weight := float64(i + 1)
		weightedSum += prices[start+i] * weight
		weightSum += weight
	}

	wma := weightedSum / weightSum
	return util.Round(wma)
}

func (wma *WMA) GetDataRange() int32 {
	return int32(wma.period)
}

func (indicator *WMA) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
