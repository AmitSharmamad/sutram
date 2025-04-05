package parser

import (
	"algo/src/util"
)

type EMA struct {
	period int
	offset int
}

func (indicator *EMA) SliceData(prices Prices) Prices {
	start := len(prices) - indicator.offset - indicator.period
	end := len(prices) - indicator.offset
	return prices[start:end]
}

func (indicator *EMA) Calculate(
	config []any,
	prices Prices,
) float64 {

	period := int(config[1].(float64))
	indicator.period = period
	indicator.offset = int(config[2].(float64))

	// println(point, period, offset)

	prices = indicator.SliceData(prices)

	if len(prices) < period {
		return 0.0
	}

	// Calculate multiplier
	multiplier := 2.0 / float64(period+1)

	// Calculate initial SMA
	var sma float64
	for i := 0; i < period; i++ {
		sma += prices[i]
	}
	sma /= float64(period)

	// Calculate EMA
	ema := sma
	for i := period; i < len(prices); i++ {
		ema = (prices[i]-ema)*multiplier + ema
	}

	return util.Round(ema)
}

func (indicator *EMA) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 14 // default period
	}

	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.period + indicator.offset)
}
