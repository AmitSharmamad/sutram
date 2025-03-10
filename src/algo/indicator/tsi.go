package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type TSI struct {
	fastPeriod int
	slowPeriod int
}

func (indicator *TSI) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	fastPeriod := indicator.fastPeriod
	if fastPeriod == 0 {
		if p, ok := config["fastPeriod"].(int); ok {
			fastPeriod = p
		} else {
			fastPeriod = 13 // default fast period
		}
	}

	slowPeriod := indicator.slowPeriod
	if slowPeriod == 0 {
		if p, ok := config["slowPeriod"].(int); ok {
			slowPeriod = p
		} else {
			slowPeriod = 25 // default slow period
		}
	}

	if len(prices) < slowPeriod+1 {
		return 0.0
	}

	// Calculate price changes
	momentum := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		momentum[i-1] = prices[i] - prices[i-1]
	}

	// Double Smoothed Momentum
	ema1 := &EMA{period: slowPeriod}
	firstEMA := ema1.Calculate(Prices(momentum), nil, point)

	ema2 := &EMA{period: fastPeriod}
	smoothedMomentum := ema2.Calculate(Prices{firstEMA}, nil, point)

	// Double Smoothed Absolute Momentum
	absMomentum := make([]float64, len(momentum))
	for i := 0; i < len(momentum); i++ {
		if momentum[i] < 0 {
			absMomentum[i] = -momentum[i]
		} else {
			absMomentum[i] = momentum[i]
		}
	}

	ema3 := &EMA{period: slowPeriod}
	firstAbsEMA := ema3.Calculate(Prices(absMomentum), nil, point)

	ema4 := &EMA{period: fastPeriod}
	smoothedAbsMomentum := ema4.Calculate(Prices{firstAbsEMA}, nil, point)

	if smoothedAbsMomentum == 0 {
		return 0.0
	}

	// Calculate TSI
	tsi := (smoothedMomentum / smoothedAbsMomentum) * 100

	return util.Round(tsi)
}

func (tsi *TSI) GetDataRange() int32 {
	return int32(tsi.slowPeriod + tsi.fastPeriod)
}

func (indicator *TSI) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
