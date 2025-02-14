package shared

import (
	"algo/src/util"
	"math"
)

type BollingerBands struct {
	period  int
	stdDevs float64
}

func (indicator *BollingerBands) Calculate(
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

	stdDevs := indicator.stdDevs
	if stdDevs == 0 {
		if s, ok := config["stdDevs"].(float64); ok {
			stdDevs = s
		} else {
			stdDevs = 2.0 // default standard deviations
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate SMA
	var sum float64
	start := len(prices) - period
	for i := start; i < len(prices); i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)

	// Calculate Standard Deviation
	var variance float64
	for i := start; i < len(prices); i++ {
		diff := prices[i] - sma
		variance += diff * diff
	}
	variance /= float64(period)
	stdDev := math.Sqrt(variance)

	// Calculate Bands - storing them in case we want to expose them later via a different method
	_ = sma + (stdDev * stdDevs) // upper band
	_ = sma - (stdDev * stdDevs) // lower band

	// Return middle band (SMA) by default
	return util.Round(sma)
}

func (bb *BollingerBands) GetDataRange() int32 {
	return int32(bb.period)
}
