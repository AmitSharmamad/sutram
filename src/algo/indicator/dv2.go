package indicator

import (
	"algo/src/util"
)

type DV2 struct {
	emaPeriod int
	rsiPeriod int
}

func (indicator *DV2) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	emaPeriod := indicator.emaPeriod
	if emaPeriod == 0 {
		if p, ok := config["emaPeriod"].(int); ok {
			emaPeriod = p
		} else {
			emaPeriod = 3 // default EMA period
		}
	}

	rsiPeriod := indicator.rsiPeriod
	if rsiPeriod == 0 {
		if p, ok := config["rsiPeriod"].(int); ok {
			rsiPeriod = p
		} else {
			rsiPeriod = 14 // default RSI period
		}
	}

	if len(prices) < rsiPeriod {
		return 0.0
	}

	// Calculate RSI
	rsi := &RSI{period: rsiPeriod}
	rsiValue := rsi.Calculate(prices, nil, point)

	// Calculate EMA of close prices
	ema := &EMA{period: emaPeriod}
	emaValue := ema.Calculate(prices, nil, point)

	// DV2 formula: (RSI - 50) * 0.32 + (Close - EMA) / Close * 100
	dv2 := (rsiValue - 50) * 0.32
	if emaValue != 0 {
		dv2 += (prices[len(prices)-1] - emaValue) / prices[len(prices)-1] * 100
	}

	return util.Round(dv2)
}

func (dv2 *DV2) GetDataRange() int32 {
	if dv2.rsiPeriod > dv2.emaPeriod {
		return int32(dv2.rsiPeriod)
	}
	return int32(dv2.emaPeriod)
}
