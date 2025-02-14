package indicator

import (
	"algo/src/util"
	"math"
)

type RAVI struct {
	shortPeriod int
	longPeriod  int
}

func (indicator *RAVI) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	shortPeriod := indicator.shortPeriod
	if shortPeriod == 0 {
		if p, ok := config["shortPeriod"].(int); ok {
			shortPeriod = p
		} else {
			shortPeriod = 7 // default short period
		}
	}

	longPeriod := indicator.longPeriod
	if longPeriod == 0 {
		if p, ok := config["longPeriod"].(int); ok {
			longPeriod = p
		} else {
			longPeriod = 65 // default long period
		}
	}

	if len(prices) < longPeriod {
		return 0.0
	}

	// Calculate short-term SMA
	shortSMA := &SMA{period: shortPeriod}
	shortValue := shortSMA.Calculate(prices, nil, point)

	// Calculate long-term SMA
	longSMA := &SMA{period: longPeriod}
	longValue := longSMA.Calculate(prices, nil, point)

	if longValue == 0 {
		return 0.0
	}

	// RAVI = |Short SMA - Long SMA| / Long SMA * 100
	ravi := math.Abs(shortValue-longValue) / longValue * 100

	return util.Round(ravi)
}

func (ravi *RAVI) GetDataRange() int32 {
	if ravi.longPeriod > ravi.shortPeriod {
		return int32(ravi.longPeriod)
	}
	return int32(ravi.shortPeriod)
}
