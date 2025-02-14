package shared

import (
	"algo/src/util"
	"math"
)

type ATR struct {
	period int
}

func (indicator *ATR) Calculate(
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

	// Calculate True Range series
	trueRanges := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		high := prices[i]
		low := prices[i-1]
		close := prices[i-1]

		tr1 := math.Abs(high - low)
		tr2 := math.Abs(high - close)
		tr3 := math.Abs(low - close)
		trueRanges[i-1] = math.Max(tr1, math.Max(tr2, tr3))
	}

	// Calculate initial ATR as simple average
	var sum float64
	for i := 0; i < period; i++ {
		sum += trueRanges[i]
	}
	atr := sum / float64(period)

	// Calculate subsequent ATR values using smoothing
	multiplier := 1.0 / float64(period)
	for i := period; i < len(trueRanges); i++ {
		atr = (atr*float64(period-1) + trueRanges[i]) * multiplier
	}

	return util.Round(atr)
}

func (atr *ATR) GetDataRange() int32 {
	return int32(atr.period)
}
