package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type VPT struct {
	period int
}

func (indicator *VPT) Calculate(
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

	// Since we don't have volume data in the current implementation,
	// we'll use a simplified version with assumed volume
	var vpt float64
	previousClose := prices[len(prices)-2]
	currentClose := prices[len(prices)-1]

	// VPT = Previous VPT + Volume * ((Today's Close - Yesterday's Close) / Yesterday's Close)
	priceChange := (currentClose - previousClose) / previousClose
	assumedVolume := 1.0 // In real implementation, this would be actual volume
	vpt += assumedVolume * priceChange

	// Smooth VPT with EMA
	ema := &EMA{period: period}
	smoothedVPT := ema.Calculate(Prices{vpt}, nil, point)

	return util.Round(smoothedVPT)
}

func (vpt *VPT) GetDataRange() int32 {
	return int32(vpt.period)
}

func (indicator *VPT) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
