package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type HMAChaikin struct {
	period int
}

func (indicator *HMAChaikin) Calculate(
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

	if len(prices) < period {
		return 0.0
	}

	// Calculate square root of period for HMA
	sqrtPeriod := int(math.Sqrt(float64(period)))

	// Calculate money flow multiplier using price changes as proxy for volume
	var mfm float64
	close := prices[len(prices)-1]
	prevClose := prices[len(prices)-2]
	priceRange := math.Abs(close - prevClose)

	if priceRange != 0 {
		mfm = (close - prevClose) / priceRange
	}

	// Calculate weighted prices with money flow
	weightedPrices := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		weightedPrices[i] = prices[i] * (1 + mfm)
	}

	// Calculate WMA with period/2 and full period
	wma1 := &WMA{period: period / 2}
	wma2 := &WMA{period: period}

	// Get WMA values for weighted prices
	wmaPeriodHalf := wma1.Calculate(Prices(weightedPrices), nil, point)
	wmaPeriodFull := wma2.Calculate(Prices(weightedPrices), nil, point)

	// Calculate raw Hull: (2 * WMA(n/2)) - WMA(n)
	rawHull := 2*wmaPeriodHalf - wmaPeriodFull

	// Take WMA of the raw Hull using sqrt(n) period
	wma3 := &WMA{period: sqrtPeriod}
	hmaChaikin := wma3.Calculate(Prices{rawHull}, nil, point)

	return util.Round(hmaChaikin)
}

func (hc *HMAChaikin) GetDataRange() int32 {
	return int32(hc.period)
}

func (indicator *HMAChaikin) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
