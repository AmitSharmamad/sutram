package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type AO struct {
	fastPeriod int
	slowPeriod int
}

func (indicator *AO) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	fastPeriod := indicator.fastPeriod
	if fastPeriod == 0 {
		if p, ok := config["fastPeriod"].(int); ok {
			fastPeriod = p
		} else {
			fastPeriod = 5 // default fast period
		}
	}

	slowPeriod := indicator.slowPeriod
	if slowPeriod == 0 {
		if p, ok := config["slowPeriod"].(int); ok {
			slowPeriod = p
		} else {
			slowPeriod = 34 // default slow period
		}
	}

	if len(prices) < slowPeriod {
		return 0.0
	}

	// Calculate SMAs
	fastSMA := &SMA{period: fastPeriod}
	slowSMA := &SMA{period: slowPeriod}

	fastValue := fastSMA.Calculate(prices, nil, point)
	slowValue := slowSMA.Calculate(prices, nil, point)

	// AO = Fast SMA - Slow SMA
	ao := fastValue - slowValue

	return util.Round(ao)
}

func (ao *AO) GetDataRange() int32 {
	if ao.slowPeriod > ao.fastPeriod {
		return int32(ao.slowPeriod)
	}
	return int32(ao.fastPeriod)
}

func (indicator *AO) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
