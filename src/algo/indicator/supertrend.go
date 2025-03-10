package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type SuperTrend struct {
	period     int
	multiplier float64
}

func (indicator *SuperTrend) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 7 // default period
		}
	}

	multiplier := indicator.multiplier
	if multiplier == 0 {
		if m, ok := config["multiplier"].(float64); ok {
			multiplier = m
		} else {
			multiplier = 3.0 // default multiplier
		}
	}

	if len(prices) < period+1 {
		return 0.0
	}

	// Calculate ATR
	atr := &ATR{period: period}
	atrValue := atr.Calculate(prices, nil, point)

	// Calculate basic upper and lower bands
	medianPrice := prices[len(prices)-1]
	upperBand := medianPrice + (multiplier * atrValue)
	lowerBand := medianPrice - (multiplier * atrValue)

	// For simplicity, we return the upper band if price is above median, lower band if below
	if prices[len(prices)-1] > medianPrice {
		return util.Round(upperBand)
	}
	return util.Round(lowerBand)
}

func (st *SuperTrend) GetDataRange() int32 {
	return int32(st.period)
}

func (indicator *SuperTrend) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
