package shared

import "algo/src/util"

type SMA struct {
	period int
	offset int
}

func (indicator *SMA) Calculate(
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

	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(int); ok {
			offset = o
		}
	}

	if len(prices) < period+offset {
		return 0.0
	}

	// Calculate SMA
	var sum float64
	start := len(prices) - period - offset
	end := len(prices) - offset

	for i := start; i < end; i++ {
		sum += prices[i]
	}

	return util.Round(sum / float64(period))
}

func (sma *SMA) GetDataRange() int32 {
	return int32(sma.period)
}
