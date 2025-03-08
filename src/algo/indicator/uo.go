package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type UO struct {
	period1 int
	period2 int
	period3 int
}

func (indicator *UO) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period1 := indicator.period1
	if period1 == 0 {
		if p, ok := config["period1"].(int); ok {
			period1 = p
		} else {
			period1 = 7 // default first period
		}
	}

	period2 := indicator.period2
	if period2 == 0 {
		if p, ok := config["period2"].(int); ok {
			period2 = p
		} else {
			period2 = 14 // default second period
		}
	}

	period3 := indicator.period3
	if period3 == 0 {
		if p, ok := config["period3"].(int); ok {
			period3 = p
		} else {
			period3 = 28 // default third period
		}
	}

	if len(prices) < period3+1 {
		return 0.0
	}

	// Calculate buying pressure (BP) and true range (TR)
	bpSum1, trSum1 := 0.0, 0.0
	bpSum2, trSum2 := 0.0, 0.0
	bpSum3, trSum3 := 0.0, 0.0

	for i := len(prices) - period3; i < len(prices); i++ {
		// Calculate BP = Close - min(Low, Prior Close)
		close := prices[i]
		prevClose := prices[i-1]
		low := math.Min(prices[i], prevClose)
		bp := close - low

		// Calculate TR = max(High, Prior Close) - min(Low, Prior Close)
		high := math.Max(prices[i], prevClose)
		tr := high - low

		// Sum for each period
		if i >= len(prices)-period1 {
			bpSum1 += bp
			trSum1 += tr
		}
		if i >= len(prices)-period2 {
			bpSum2 += bp
			trSum2 += tr
		}
		bpSum3 += bp
		trSum3 += tr
	}

	// Prevent division by zero
	if trSum1 == 0 || trSum2 == 0 || trSum3 == 0 {
		return 0.0
	}

	// Calculate averages with weightings
	avg1 := (bpSum1 / trSum1) * 4
	avg2 := (bpSum2 / trSum2) * 2
	avg3 := (bpSum3 / trSum3) * 1

	// Calculate UO
	uo := (avg1 + avg2 + avg3) * 100 / 7

	return util.Round(uo)
}

func (uo *UO) GetDataRange() int32 {
	maxPeriod := uo.period1
	if uo.period2 > maxPeriod {
		maxPeriod = uo.period2
	}
	if uo.period3 > maxPeriod {
		maxPeriod = uo.period3
	}
	return int32(maxPeriod)
}

func (indicator *UO) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
