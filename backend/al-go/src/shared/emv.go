package shared

import (
	"algo/src/util"
	"math"
)

type EMV struct {
	period int
}

func (indicator *EMV) Calculate(
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

	// Since we don't have actual volume data, we'll use a simplified version
	var emv float64
	for i := len(prices) - period; i < len(prices); i++ {
		// Calculate midpoint move
		currentMid := (prices[i] + prices[i-1]) / 2
		prevMid := (prices[i-1] + prices[i-2]) / 2
		move := currentMid - prevMid

		// Calculate box ratio (normally would use volume)
		// Using price difference as a proxy for box ratio
		boxRatio := (prices[i] - prices[i-1]) / 1000000 // Normalized value

		if boxRatio != 0 {
			// Distance / Box Ratio
			emv += move / math.Abs(boxRatio)
		}
	}

	// Return average EMV over the period
	return util.Round(emv / float64(period))
}

func (emv *EMV) GetDataRange() int32 {
	return int32(emv.period)
}
