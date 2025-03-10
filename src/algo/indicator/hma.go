package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type HMA struct {
	period int
}

func (indicator *HMA) Calculate(
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

	// Calculate period for intermediate WMAs
	sqrtPeriod := int(math.Sqrt(float64(period)))

	// Calculate WMA with period/2 and full period
	wma1 := &WMA{period: period / 2}
	wma2 := &WMA{period: period}

	// Get WMA values
	wmaPeriodHalf := wma1.Calculate(prices, nil, point)
	wmaPeriodFull := wma2.Calculate(prices, nil, point)

	// Calculate raw HMA: (2 * WMA(n/2)) - WMA(n)
	rawHMA := 2*wmaPeriodHalf - wmaPeriodFull

	// Take WMA of the raw HMA using sqrt(n) period
	wma3 := &WMA{period: sqrtPeriod}
	hma := wma3.Calculate(Prices{rawHMA}, nil, point)

	return util.Round(hma)
}

func (hma *HMA) GetDataRange() int32 {
	return int32(hma.period)
}

func (indicator *HMA) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
