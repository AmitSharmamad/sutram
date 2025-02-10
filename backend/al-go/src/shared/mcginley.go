package shared

import (
	"algo/src/util"
	"math"
)

type McGinley struct {
	period int
}

func (indicator *McGinley) Calculate(
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

	if len(prices) < period {
		return 0.0
	}

	// Initialize McGinley Dynamic with first price
	md := prices[len(prices)-period]

	// Calculate McGinley Dynamic
	for i := len(prices) - period + 1; i < len(prices); i++ {
		price := prices[i]
		// McGinley Dynamic formula: MD = MD[i-1] + (Price - MD[i-1]) / (N * pow(Price/MD[i-1], 4))
		md = md + (price-md)/(float64(period)*math.Pow(price/md, 4))
	}

	return util.Round(md)
}

func (md *McGinley) GetDataRange() int32 {
	return int32(md.period)
}
