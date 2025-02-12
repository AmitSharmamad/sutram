package indicator

import (
	"algo/src/util"
)

type DPO struct {
	period int
}

func (indicator *DPO) Calculate(
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

	if len(prices) < period+1 {
		return 0.0
	}

	// DPO = Price - SMA(Price, period/2 + 1) shifted period/2 + 1 bars backward
	shift := period/2 + 1

	// Calculate SMA for the shifted period
	var sum float64
	start := len(prices) - period - shift
	for i := start; i < start+period; i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)

	// Current price minus the shifted SMA
	dpo := prices[len(prices)-1] - sma

	return util.Round(dpo)
}

func (dpo *DPO) GetDataRange() int32 {
	return int32(dpo.period * 2)
}
