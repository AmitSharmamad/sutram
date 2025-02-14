package shared

import "algo/src/util"

type TRIX struct {
	period int
}

func (indicator *TRIX) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 15 // default period
		}
	}

	if len(prices) < period*3+1 {
		return 0.0
	}

	// Calculate first EMA
	ema1 := &EMA{period: period}
	ema1Value := ema1.Calculate(prices, nil, point)

	// Calculate second EMA
	ema2 := &EMA{period: period}
	ema2Prices := append(make(Prices, 0), ema1Value)
	ema2Value := ema2.Calculate(ema2Prices, nil, point)

	// Calculate third EMA
	ema3 := &EMA{period: period}
	ema3Prices := append(make(Prices, 0), ema2Value)
	currentEMA3 := ema3.Calculate(ema3Prices, nil, point)

	// Calculate previous triple EMA for ROC
	prevPrices := prices[:len(prices)-1]
	prevEma1Value := ema1.Calculate(prevPrices, nil, point)
	prevEma2Prices := append(make(Prices, 0), prevEma1Value)
	prevEma2Value := ema2.Calculate(prevEma2Prices, nil, point)
	prevEma3Prices := append(make(Prices, 0), prevEma2Value)
	prevEMA3 := ema3.Calculate(prevEma3Prices, nil, point)

	// Calculate TRIX
	if prevEMA3 == 0 {
		return 0.0
	}

	trix := ((currentEMA3 - prevEMA3) / prevEMA3) * 100

	return util.Round(trix)
}

func (trix *TRIX) GetDataRange() int32 {
	return int32(trix.period * 3)
}
