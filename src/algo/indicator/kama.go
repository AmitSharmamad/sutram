package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type KAMA struct {
	period int
	fast   int
	slow   int
}

func (indicator *KAMA) Calculate(
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

	fast := indicator.fast
	if fast == 0 {
		if f, ok := config["fast"].(int); ok {
			fast = f
		} else {
			fast = 2 // default fast EMA period
		}
	}

	slow := indicator.slow
	if slow == 0 {
		if s, ok := config["slow"].(int); ok {
			slow = s
		} else {
			slow = 30 // default slow EMA period
		}
	}

	if len(prices) < period+1 {
		return 0.0
	}

	// Calculate price change and volatility
	change := math.Abs(prices[len(prices)-1] - prices[len(prices)-period])

	var volatility float64
	for i := len(prices) - period; i < len(prices)-1; i++ {
		volatility += math.Abs(prices[i+1] - prices[i])
	}

	// Calculate efficiency ratio
	er := 0.0
	if volatility != 0 {
		er = change / volatility
	}

	// Calculate smoothing constant
	fastSC := 2.0 / (float64(fast) + 1.0)
	slowSC := 2.0 / (float64(slow) + 1.0)
	sc := math.Pow(er*(fastSC-slowSC)+slowSC, 2)

	// Calculate KAMA
	var kama float64
	if len(prices) > period {
		prevKama := prices[len(prices)-2] // Use previous price as initial KAMA
		kama = prevKama + sc*(prices[len(prices)-1]-prevKama)
	} else {
		kama = prices[len(prices)-1]
	}

	return util.Round(kama)
}

func (kama *KAMA) GetDataRange() int32 {
	return int32(kama.period)
}

func (indicator *KAMA) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
