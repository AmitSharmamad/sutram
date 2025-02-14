package indicator

import (
	"algo/src/util"
)

type VWAP struct {
	period int
}

func (indicator *VWAP) Calculate(
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

	// Since we don't have actual volume data in the current implementation,
	// we'll use a simplified version that assumes equal volume weights
	// In a real implementation, you would multiply price by actual volume

	var sumPV float64
	var sumV float64

	// Calculate sum of Price * Volume and sum of Volume
	for i := len(prices) - period; i < len(prices); i++ {
		price := prices[i]
		volume := 1.0 // In real implementation, this would be actual volume
		sumPV += price * volume
		sumV += volume
	}

	if sumV == 0 {
		return 0.0
	}

	// VWAP = Sum(Price * Volume) / Sum(Volume)
	vwap := sumPV / sumV

	return util.Round(vwap)
}

func (vwap *VWAP) GetDataRange() int32 {
	return int32(vwap.period)
}
