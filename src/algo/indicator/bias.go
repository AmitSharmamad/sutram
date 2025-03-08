package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type BiasIndicator struct {
	period int
}

func (indicator *BiasIndicator) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 26 // default period
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate SMA
	var sma float64
	for i := 0; i < period; i++ {
		sma += prices[i]
	}
	sma /= float64(period)

	// Calculate Bias
	bias := (prices[len(prices)-1] / sma) - 1

	return util.Round(bias * 100)
}

func (indicator *BiasIndicator) GetDataRange() int32 {
	return int32(indicator.period)
}

func (indicator *BiasIndicator) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
