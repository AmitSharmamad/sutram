package indicator

import (
	"algo/src/util"
	"math"
)

type STC struct {
	fastLength  int
	slowLength  int
	cycleLength int
}

func (indicator *STC) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	fastLength := indicator.fastLength
	if fastLength == 0 {
		if f, ok := config["fastLength"].(int); ok {
			fastLength = f
		} else {
			fastLength = 23 // default fast length
		}
	}

	slowLength := indicator.slowLength
	if slowLength == 0 {
		if s, ok := config["slowLength"].(int); ok {
			slowLength = s
		} else {
			slowLength = 50 // default slow length
		}
	}

	cycleLength := indicator.cycleLength
	if cycleLength == 0 {
		if c, ok := config["cycleLength"].(int); ok {
			cycleLength = c
		} else {
			cycleLength = 10 // default cycle length
		}
	}

	if len(prices) < slowLength+cycleLength {
		return 0.0
	}

	// Calculate MACD-like values using EMAs
	fastMA := &EMA{period: fastLength}
	slowMA := &EMA{period: slowLength}

	fastValue := fastMA.Calculate(prices, nil, point)
	slowValue := slowMA.Calculate(prices, nil, point)
	macd := fastValue - slowValue

	// Calculate Stochastic of MACD
	var highestMACD, lowestMACD float64
	highestMACD = macd
	lowestMACD = macd

	// Find highest and lowest MACD values over cycle length
	for i := 0; i < cycleLength; i++ {
		if macd > highestMACD {
			highestMACD = macd
		}
		if macd < lowestMACD {
			lowestMACD = macd
		}
	}

	// Calculate Fast K
	var fastK float64
	if highestMACD != lowestMACD {
		fastK = ((macd - lowestMACD) / (highestMACD - lowestMACD)) * 100
	}

	// Smooth Fast K to get Fast D (first smoothing)
	ema1 := &EMA{period: cycleLength}
	fastD := ema1.Calculate(Prices{fastK}, nil, point)

	// Second smoothing
	ema2 := &EMA{period: cycleLength}
	stc := ema2.Calculate(Prices{fastD}, nil, point)

	// Ensure output is between 0 and 100
	stc = math.Max(0, math.Min(100, stc))

	return util.Round(stc)
}

func (stc *STC) GetDataRange() int32 {
	maxPeriod := stc.slowLength
	if stc.cycleLength > maxPeriod {
		maxPeriod = stc.cycleLength
	}
	return int32(maxPeriod + stc.cycleLength)
}
