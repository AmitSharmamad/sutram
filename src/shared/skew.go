package shared

import (
	"algo/src/util"
	"math"
)

type Skew struct {
	period int
}

func (indicator *Skew) Calculate(
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

	// Calculate mean
	var sum float64
	start := len(prices) - period
	for i := start; i < len(prices); i++ {
		sum += prices[i]
	}
	mean := sum / float64(period)

	// Calculate standard deviation and third moment
	var variance, thirdMoment float64
	for i := start; i < len(prices); i++ {
		dev := prices[i] - mean
		variance += dev * dev
		thirdMoment += dev * dev * dev
	}
	variance /= float64(period)
	thirdMoment /= float64(period)

	stdDev := math.Sqrt(variance)
	if stdDev == 0 {
		return 0.0
	}

	// Calculate skewness
	skew := thirdMoment / math.Pow(stdDev, 3)

	return util.Round(skew)
}

func (skew *Skew) GetDataRange() int32 {
	return int32(skew.period)
}
