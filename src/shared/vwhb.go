package shared

import (
	"algo/src/util"
	"math"
)

type VWHB struct {
	period  int
	stdDevs float64
}

func (indicator *VWHB) Calculate(
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

	stdDevs := indicator.stdDevs
	if stdDevs == 0 {
		if s, ok := config["stdDevs"].(float64); ok {
			stdDevs = s
		} else {
			stdDevs = 2.0 // default standard deviations
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate volume proxy based on price changes
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

	// Calculate Hull MA on volume-weighted prices
	hma := &HMA{period: period}
	middle := hma.Calculate(Prices(vwPrices), nil, point)

	// Calculate standard deviation
	var sumSquares float64
	for i := len(vwPrices) - period; i < len(vwPrices); i++ {
		diff := vwPrices[i] - middle
		sumSquares += diff * diff
	}
	stdDev := math.Sqrt(sumSquares / float64(period))

	// Calculate bands and store them for potential future use
	_ = middle + (stdDev * stdDevs) // upper band
	_ = middle - (stdDev * stdDevs) // lower band

	// Return middle line as primary signal
	return util.Round(middle)
}

func (vwhb *VWHB) GetDataRange() int32 {
	return int32(vwhb.period)
}
