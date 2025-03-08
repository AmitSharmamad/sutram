package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type CHOP struct {
	period int
}

func (indicator *CHOP) Calculate(
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

	if len(prices) < period+1 {
		return 0.0
	}

	// Calculate ATR sum
	var atrSum float64
	for i := len(prices) - period; i < len(prices); i++ {
		high := prices[i]
		low := prices[i-1]
		close := prices[i-1]

		tr1 := math.Abs(high - low)
		tr2 := math.Abs(high - close)
		tr3 := math.Abs(low - close)
		atrSum += math.Max(tr1, math.Max(tr2, tr3))
	}

	// Find highest high and lowest low
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

	if high-low == 0 {
		return 0.0
	}

	// CHOP = 100 * LOG10( SUM(ATR(1)) / (MaxHi - MinLo) ) / LOG10(period)
	chop := 100 * math.Log10(atrSum/(high-low)) / math.Log10(float64(period))

	return util.Round(chop)
}

func (chop *CHOP) GetDataRange() int32 {
	return int32(chop.period)
}

func (indicator *CHOP) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
