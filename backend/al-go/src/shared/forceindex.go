package shared

import (
	"algo/src/util"
	"math"
)

type ForceIndex struct {
	period int
}

func (indicator *ForceIndex) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 13 // default period
		}
	}

	if len(prices) < period+1 {
		return 0.0
	}

	// Calculate raw Force Index
	// Since we don't have volume data, we'll use price changes as a proxy for volume
	priceChange := prices[len(prices)-1] - prices[len(prices)-2]
	volumeProxy := math.Abs(priceChange) // Using price change magnitude as volume proxy
	rawForce := priceChange * volumeProxy

	// Smooth the Force Index with EMA
	ema := &EMA{period: period}
	smoothedForce := ema.Calculate(Prices{rawForce}, nil, point)

	return util.Round(smoothedForce)
}

func (fi *ForceIndex) GetDataRange() int32 {
	return int32(fi.period)
}
