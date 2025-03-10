package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type PM struct {
	shortPeriod int
	longPeriod  int
}

func (indicator *PM) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	shortPeriod := indicator.shortPeriod
	if shortPeriod == 0 {
		if p, ok := config["shortPeriod"].(int); ok {
			shortPeriod = p
		} else {
			shortPeriod = 10 // default short period
		}
	}

	longPeriod := indicator.longPeriod
	if longPeriod == 0 {
		if p, ok := config["longPeriod"].(int); ok {
			longPeriod = p
		} else {
			longPeriod = 20 // default long period
		}
	}

	if len(prices) < longPeriod {
		return 0.0
	}

	// Calculate short-term momentum
	shortMom := (prices[len(prices)-1] - prices[len(prices)-shortPeriod]) / prices[len(prices)-shortPeriod] * 100

	// Calculate long-term momentum
	longMom := (prices[len(prices)-1] - prices[len(prices)-longPeriod]) / prices[len(prices)-longPeriod] * 100

	// Calculate PM oscillator
	pm := (shortMom + longMom) / 2

	// Apply smoothing using EMA
	ema := &EMA{period: 5} // Short smoothing period
	smoothedPM := ema.Calculate(Prices{pm}, nil, point)

	return util.Round(smoothedPM)
}

func (pm *PM) GetDataRange() int32 {
	if pm.longPeriod > pm.shortPeriod {
		return int32(pm.longPeriod)
	}
	return int32(pm.shortPeriod)
}

func (indicator *PM) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
