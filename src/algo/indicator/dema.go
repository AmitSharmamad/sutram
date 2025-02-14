package indicator

import (
	"algo/src/util"
)

type DEMA struct {
	period int
}

func (indicator *DEMA) Calculate(
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

	if len(prices) < period*2 {
		return 0.0
	}

	// Calculate first EMA
	ema1 := &EMA{period: period}
	ema1Value := ema1.Calculate(prices, nil, point)

	// Calculate EMA of EMA
	ema2 := &EMA{period: period}
	ema2Prices := append(make(Prices, 0), ema1Value)
	ema2Value := ema2.Calculate(ema2Prices, nil, point)

	// DEMA = (2 × EMA1) - EMA2
	dema := (2 * ema1Value) - ema2Value

	return util.Round(dema)
}

func (dema *DEMA) GetDataRange() int32 {
	return int32(dema.period * 2)
}
