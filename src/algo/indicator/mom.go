package indicator

import (
	"algo/src/util"
)

type MOM struct {
	period int
}

func (indicator *MOM) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 10 // default period
		}
	}

	if len(prices) < period+1 {
		return 0.0
	}

	// Simple momentum calculation: current price minus price 'period' days ago
	momentum := prices[len(prices)-1] - prices[len(prices)-1-period]

	return util.Round(momentum)
}

func (mom *MOM) GetDataRange() int32 {
	return int32(mom.period)
}
