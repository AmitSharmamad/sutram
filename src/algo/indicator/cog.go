package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type COG struct {
	period int
}

func (indicator *COG) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 10 // default period
		}
	}

	if len(prices) < period {
		return 0.0
	}

	var numerator, denominator float64

	// Calculate Center of Gravity
	for i := 0; i < period; i++ {
		price := prices[len(prices)-period+i]
		numerator += price * float64(i+1)
		denominator += price
	}

	if denominator == 0 {
		return 0.0
	}

	// COG = -((Period * Sum(Price * t)) / Sum(Price)) + (Period + 1) / 2
	cog := -((float64(period) * numerator) / denominator) + (float64(period+1) / 2.0)

	return util.Round(cog)
}

func (cog *COG) GetDataRange() int32 {
	return int32(cog.period)
}

func (indicator *COG) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
