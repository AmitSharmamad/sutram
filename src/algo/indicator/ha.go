package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type HA struct {
	period int
}

func (indicator *HA) Calculate(
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

	if len(prices) < period+1 {
		return 0.0
	}

	// Get current and previous prices
	close := prices[len(prices)-1]
	open := prices[len(prices)-2]

	// Calculate Heikin-Ashi values
	haClose := (open + close) / 2

	// Return HA Close as the main signal
	return util.Round(haClose)
}

func (ha *HA) GetDataRange() int32 {
	return int32(ha.period)
}

func (indicator *HA) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
