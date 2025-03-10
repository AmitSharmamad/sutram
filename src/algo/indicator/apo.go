package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type ApoIndicator struct {
	fast int
	slow int
}

func (indicator *ApoIndicator) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	fast := indicator.fast
	if fast == 0 {
		if f, ok := config["fast"].(int); ok {
			fast = f
		} else {
			fast = 12 // default fast period
		}
	}

	slow := indicator.slow
	if slow == 0 {
		if s, ok := config["slow"].(int); ok {
			slow = s
		} else {
			slow = 26 // default slow period
		}
	}

	if len(prices) < slow {
		return 0.0
	}

	// Calculate fast SMA
	var fastSma float64
	for i := 0; i < fast; i++ {
		fastSma += prices[i]
	}
	fastSma /= float64(fast)

	// Calculate slow SMA
	var slowSma float64
	for i := 0; i < slow; i++ {
		slowSma += prices[i]
	}
	slowSma /= float64(slow)

	// Calculate APO
	apo := fastSma - slowSma

	return util.Round(apo)
}

func (indicator *ApoIndicator) GetDataRange() int32 {
	return int32(indicator.slow)
}

func (indicator *ApoIndicator) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
