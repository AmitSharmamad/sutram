package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type ElderRay struct {
	period int
}

func (indicator *ElderRay) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 13 // default period
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate EMA
	ema := &EMA{period: period}
	emaValue := ema.Calculate(prices, nil, point)

	// Get current high and low
	high := prices[len(prices)-1]
	low := prices[len(prices)-2]

	// Calculate Bull and Bear Power
	bullPower := high - emaValue
	bearPower := low - emaValue

	// Return the difference between bull and bear power as a composite indicator
	return util.Round(bullPower - bearPower)
}

func (er *ElderRay) GetDataRange() int32 {
	return int32(er.period)
}

func (indicator *ElderRay) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
