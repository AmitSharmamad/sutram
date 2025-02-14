package indicator

import (
	"algo/src/util"
	"math"
)

type QQE struct {
	rsiPeriod  int
	sfPeriod   int // Smoothing factor period
	multFactor float64
}

func (indicator *QQE) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	rsiPeriod := indicator.rsiPeriod
	if rsiPeriod == 0 {
		if p, ok := config["rsiPeriod"].(int); ok {
			rsiPeriod = p
		} else {
			rsiPeriod = 14 // default RSI period
		}
	}

	sfPeriod := indicator.sfPeriod
	if sfPeriod == 0 {
		if p, ok := config["sfPeriod"].(int); ok {
			sfPeriod = p
		} else {
			sfPeriod = 5 // default smoothing period
		}
	}

	multFactor := indicator.multFactor
	if multFactor == 0 {
		if m, ok := config["multFactor"].(float64); ok {
			multFactor = m
		} else {
			multFactor = 4.236 // default multiplier
		}
	}

	if len(prices) < rsiPeriod*2 {
		return 0.0
	}

	// Calculate RSI
	rsi := &RSI{period: rsiPeriod}
	rsiValue := rsi.Calculate(prices, nil, point)

	// Calculate first smoothing using EMA
	ema1 := &EMA{period: sfPeriod}
	smoothRSI := ema1.Calculate(Prices{rsiValue}, nil, point)

	// Calculate second smoothing
	ema2 := &EMA{period: sfPeriod}
	doubleSmoothed := ema2.Calculate(Prices{smoothRSI}, nil, point)

	// Calculate TR (True Range) of smoothed RSI
	tr := math.Abs(smoothRSI - doubleSmoothed)

	// Calculate ATR of TR
	atr := &ATR{period: sfPeriod}
	atrValue := atr.Calculate(Prices{tr}, nil, point)

	// Calculate QQE
	qqe := doubleSmoothed + (multFactor * atrValue)

	// Normalize to 0-100 range
	if qqe > 100 {
		qqe = 100
	} else if qqe < 0 {
		qqe = 0
	}

	return util.Round(qqe)
}

func (qqe *QQE) GetDataRange() int32 {
	return int32(qqe.rsiPeriod * 2)
}
