package shared

import "algo/src/util"

type TEMA struct {
	period int
}

func (indicator *TEMA) Calculate(
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

	if len(prices) < period*3 {
		return 0.0
	}

	// Calculate first EMA
	ema1 := &EMA{period: period}
	ema1Value := ema1.Calculate(prices, nil, point)

	// Calculate EMA of EMA (second EMA)
	ema2 := &EMA{period: period}
	ema2Prices := append(make(Prices, 0), ema1Value)
	ema2Value := ema2.Calculate(ema2Prices, nil, point)

	// Calculate EMA of EMA of EMA (third EMA)
	ema3 := &EMA{period: period}
	ema3Prices := append(make(Prices, 0), ema2Value)
	ema3Value := ema3.Calculate(ema3Prices, nil, point)

	// TEMA = (3 × EMA1) - (3 × EMA2) + EMA3
	tema := (3 * ema1Value) - (3 * ema2Value) + ema3Value

	return util.Round(tema)
}

func (tema *TEMA) GetDataRange() int32 {
	return int32(tema.period * 3)
}
