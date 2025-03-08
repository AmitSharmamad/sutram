package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type WilliamsR struct {
	period int
}

func (indicator *WilliamsR) Calculate(
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

	// Find highest high and lowest low over lookback period
	high := prices[len(prices)-period]
	low := prices[len(prices)-period]
	for i := len(prices) - period; i < len(prices); i++ {
		if prices[i] > high {
			high = prices[i]
		}
		if prices[i] < low {
			low = prices[i]
		}
	}

	// Calculate Williams %R
	latest := prices[len(prices)-1]
	if high-low == 0 {
		return 0.0
	}

	// Williams %R = (Highest High - Close)/(Highest High - Lowest Low) * -100
	r := ((high - latest) / (high - low)) * -100

	return util.Round(r)
}

func (wr *WilliamsR) GetDataRange() int32 {
	return int32(wr.period)
}

func (indicator *WilliamsR) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
