package shared

import (
	"algo/src/util"
	"math"
)

type CCI struct {
	period int
}

func (indicator *CCI) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 20 // default period
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate typical price (we're using the price directly as we don't have high/low)
	typicalPrice := prices[len(prices)-1]

	// Calculate SMA of typical price
	var sum float64
	start := len(prices) - period
	for i := start; i < len(prices); i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)

	// Calculate Mean Deviation
	var meanDev float64
	for i := start; i < len(prices); i++ {
		meanDev += math.Abs(prices[i] - sma)
	}
	meanDev /= float64(period)

	// Calculate CCI
	if meanDev == 0 {
		return 0
	}

	cci := (typicalPrice - sma) / (0.015 * meanDev)

	return util.Round(cci)
}

func (cci *CCI) GetDataRange() int32 {
	return int32(cci.period)
}
