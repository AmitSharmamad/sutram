package shared

import "algo/src/util"

type RMA struct {
	period   int
	numLines int
}

func (indicator *RMA) Calculate(
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

	numLines := indicator.numLines
	if numLines == 0 {
		if n, ok := config["numLines"].(int); ok {
			numLines = n
		} else {
			numLines = 8 // default number of lines
		}
	}

	if len(prices) < period*numLines {
		return 0.0
	}

	// Calculate multiple EMAs with increasing periods
	var sumEMA float64
	for i := 1; i <= numLines; i++ {
		currentPeriod := period * i
		ema := &EMA{period: currentPeriod}
		emaValue := ema.Calculate(prices, nil, point)
		sumEMA += emaValue
	}

	// Return average of all EMAs
	rainbowMA := sumEMA / float64(numLines)

	return util.Round(rainbowMA)
}

func (rma *RMA) GetDataRange() int32 {
	return int32(rma.period * rma.numLines)
}
