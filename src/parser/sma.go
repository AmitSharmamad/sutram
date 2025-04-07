package parser

import (
	"algo/src/util"
)

type SMA struct {
	period int
	offset int
}

func (indicator *SMA) SliceData(prices Prices) Prices {
	start := len(prices) - indicator.offset - indicator.period
	end := len(prices) - indicator.offset
	return prices[start:end]
}

// SMA(OHLCV, period, offset)
func (indicator *SMA) Calculate(
	config []any,
	prices Prices,
) float64 {
	point := config[0].(string)
	period := config[1].(int)
	offset := config[2].(int)

	println(point, period, offset)

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

func (indicator *SMA) GetDataRange(config map[string]any) int32 {
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
		if p, ok := config["period"].(int); ok {
			offset = p
		} else {
			offset = 0
		}
	}

	return int32(period + offset)
}
