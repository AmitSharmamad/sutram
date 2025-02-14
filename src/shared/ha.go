package shared

import "algo/src/util"

type HA struct {
	period int
}

func (indicator *HA) Calculate(
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

	// Get current and previous prices
	close := prices[len(prices)-1]
	open := prices[len(prices)-2]

	// Calculate Heikin-Ashi values
	haClose := (open + close) / 2

	// Return HA Close as the main signal
	return util.Round(haClose)
}

func (ha *HA) GetDataRange() int32 {
	return int32(ha.period)
}
