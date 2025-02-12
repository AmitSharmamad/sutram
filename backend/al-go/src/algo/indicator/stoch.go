package indicator

import (
	"algo/src/util"
)

type Stochastic struct {
	kPeriod int
	dPeriod int
	smooth  int
}

func (indicator *Stochastic) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	kPeriod := indicator.kPeriod
	if kPeriod == 0 {
		if p, ok := config["kPeriod"].(int); ok {
			kPeriod = p
		} else {
			kPeriod = 14 // default K period
		}
	}

	dPeriod := indicator.dPeriod
	if dPeriod == 0 {
		if p, ok := config["dPeriod"].(int); ok {
			dPeriod = p
		} else {
			dPeriod = 3 // default D period
		}
	}

	smooth := indicator.smooth
	if smooth == 0 {
		if s, ok := config["smooth"].(int); ok {
			smooth = s
		} else {
			smooth = 3 // default smoothing
		}
	}

	if len(prices) < kPeriod {
		return 0.0
	}

	// Find highest high and lowest low over lookback period
	high := prices[len(prices)-kPeriod]
	low := prices[len(prices)-kPeriod]
	for i := len(prices) - kPeriod; i < len(prices); i++ {
		if prices[i] > high {
			high = prices[i]
		}
		if prices[i] < low {
			low = prices[i]
		}
	}

	// Calculate %K
	latest := prices[len(prices)-1]
	if high-low == 0 {
		return 0.0
	}

	k := ((latest - low) / (high - low)) * 100

	// For simplicity, we're returning %K
	// In a full implementation, you might want to calculate %D (SMA of %K)
	return util.Round(k)
}

func (stoch *Stochastic) GetDataRange() int32 {
	return int32(stoch.kPeriod + stoch.dPeriod)
}
