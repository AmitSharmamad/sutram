package indicator

import (
	"algo/src/util"
	"math"
)

type Vortex struct {
	period int
}

func (indicator *Vortex) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 14 // default period
		}
	}

	if len(prices) < period+1 {
		return 0.0
	}

	// Calculate +VM and -VM
	var plusVM, minusVM float64
	high := prices[len(prices)-1] // Using price as high
	low := prices[len(prices)-2]  // Using previous price as low
	prevHigh := prices[len(prices)-2]
	prevLow := prices[len(prices)-3]

	plusVM = math.Abs(high - prevLow)
	minusVM = math.Abs(low - prevHigh)

	// Calculate True Range
	tr1 := math.Abs(high - low)
	tr2 := math.Abs(high - prevHigh)
	tr3 := math.Abs(low - prevLow)
	trueRange := math.Max(tr1, math.Max(tr2, tr3))

	// Sum up the values over the period
	var sumPlusVM, sumMinusVM, sumTR float64
	for i := 0; i < period; i++ {
		sumPlusVM += plusVM
		sumMinusVM += minusVM
		sumTR += trueRange
	}

	// Calculate VI+ and VI-
	viPlus := sumPlusVM / sumTR
	viMinus := sumMinusVM / sumTR

	// Return the difference between VI+ and VI-
	return util.Round(viPlus - viMinus)
}

func (v *Vortex) GetDataRange() int32 {
	return int32(v.period)
}
