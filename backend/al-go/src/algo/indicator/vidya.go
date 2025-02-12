package indicator

import (
	"algo/src/util"
	"math"
)

type VIDYA struct {
	period int
}

func (indicator *VIDYA) Calculate(
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

	if len(prices) < period*2 {
		return 0.0
	}

	// Calculate Momentum
	var upSum, downSum float64
	for i := len(prices) - period; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			upSum += change
		} else {
			downSum -= change
		}
	}

	// Calculate CMO (Chande Momentum Oscillator)
	cmo := math.Abs((upSum - downSum) / (upSum + downSum))

	// Calculate alpha
	alpha := 2.0 / (float64(period) + 1.0)

	// Calculate VIDYA
	previousVIDYA := prices[len(prices)-2]
	currentPrice := prices[len(prices)-1]

	vidya := previousVIDYA + alpha*cmo*(currentPrice-previousVIDYA)

	return util.Round(vidya)
}

func (v *VIDYA) GetDataRange() int32 {
	return int32(v.period * 2)
}
