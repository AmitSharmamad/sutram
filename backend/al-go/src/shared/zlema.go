package shared

import "algo/src/util"

type ZLEMA struct {
	period int
}

func (indicator *ZLEMA) Calculate(
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

	if len(prices) < period {
		return 0.0
	}

	// Calculate the lag
	lag := (period - 1) / 2

	// Calculate error correction
	var zlema float64
	if len(prices) > lag {
		// Get the current price and the lagged price
		currentPrice := prices[len(prices)-1]
		laggedPrice := prices[len(prices)-1-lag]

		// Calculate the zero-lag error correction
		errorCorrection := currentPrice + (currentPrice - laggedPrice)

		// Calculate alpha
		alpha := 2.0 / float64(period+1)

		// Calculate ZLEMA
		previousZLEMA := prices[len(prices)-2] // Use previous price as initial ZLEMA
		zlema = alpha*errorCorrection + (1-alpha)*previousZLEMA
	} else {
		zlema = prices[len(prices)-1]
	}

	return util.Round(zlema)
}

func (zlema *ZLEMA) GetDataRange() int32 {
	return int32(zlema.period)
}
