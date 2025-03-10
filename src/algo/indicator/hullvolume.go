package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type HullVolume struct {
	period int
}

func (indicator *HullVolume) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 16 // default period
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate square root of period for HMA
	sqrtPeriod := int(math.Sqrt(float64(period)))

	// Calculate volatility-based volume proxy
	volProxy := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		volProxy[i-1] = math.Abs(prices[i] - prices[i-1])
	}

	// Calculate volume-weighted prices
	vwPrices := make([]float64, len(prices))
	for i := 1; i < len(prices); i++ {
		if i < len(volProxy) {
			vwPrices[i] = prices[i] * (1 + volProxy[i-1])
		} else {
			vwPrices[i] = prices[i]
		}
	}

	// Calculate WMA with period/2 and full period
	wma1 := &WMA{period: period / 2}
	wma2 := &WMA{period: period}

	// Get WMA values for volume-weighted prices
	wmaPeriodHalf := wma1.Calculate(Prices(vwPrices), nil, point)
	wmaPeriodFull := wma2.Calculate(Prices(vwPrices), nil, point)

	// Calculate raw Hull: (2 * WMA(n/2)) - WMA(n)
	rawHull := 2*wmaPeriodHalf - wmaPeriodFull

	// Take WMA of the raw Hull using sqrt(n) period
	wma3 := &WMA{period: sqrtPeriod}
	hullVol := wma3.Calculate(Prices{rawHull}, nil, point)

	return util.Round(hullVol)
}

func (hv *HullVolume) GetDataRange() int32 {
	return int32(hv.period)
}

func (indicator *HullVolume) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
