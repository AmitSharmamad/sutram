package indicator

import (
	"algo/src/util"
	"math"
)

type WaveTrend struct {
	channel1 int
	channel2 int
}

func (indicator *WaveTrend) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	channel1 := indicator.channel1
	if channel1 == 0 {
		if c, ok := config["channel1"].(int); ok {
			channel1 = c
		} else {
			channel1 = 10 // default channel1 period
		}
	}

	channel2 := indicator.channel2
	if channel2 == 0 {
		if c, ok := config["channel2"].(int); ok {
			channel2 = c
		} else {
			channel2 = 21 // default channel2 period
		}
	}

	if len(prices) < channel2 {
		return 0.0
	}

	// Calculate HLC/3 (typical price)
	typicalPrice := prices[len(prices)-1]

	// Calculate First EMA and First EMA of Absolute Values
	ema1 := &EMA{period: channel1}
	ap := ema1.Calculate(Prices{typicalPrice}, nil, point)

	// Calculate difference between price and EMA
	diff := typicalPrice - ap

	// Calculate absolute difference
	absDiff := math.Abs(diff)

	// Calculate Second EMA of Absolute Values
	ema2 := &EMA{period: channel1}
	avgAbsDiff := ema2.Calculate(Prices{absDiff}, nil, point)

	// Avoid division by zero
	if avgAbsDiff == 0 {
		return 0.0
	}

	// Calculate ci (Coefficient Index)
	ci := diff / (0.015 * avgAbsDiff)

	// Calculate Wave Trend using second smoothing
	ema3 := &EMA{period: channel2}
	wt := ema3.Calculate(Prices{ci}, nil, point)

	return util.Round(wt)
}

func (wt *WaveTrend) GetDataRange() int32 {
	if wt.channel2 > wt.channel1 {
		return int32(wt.channel2)
	}
	return int32(wt.channel1)
}
