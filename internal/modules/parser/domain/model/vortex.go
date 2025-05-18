package parser

import (
	util "algo/internal/util/math"
	"math"
)

type Vortex struct {
	period int
	offset int
}

func (indicator *Vortex) SliceData(prices Prices) Prices {
	// We need enough data to calculate Vortex indicator properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// Internal helper method to slice all required data
func (indicator *Vortex) sliceAllData(highs, lows, closes Prices) (Prices, Prices, Prices) {
	// We need enough data to calculate Vortex indicator properly
	start := max(0, len(closes)-(indicator.offset+2*indicator.period))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end]
}

func (indicator *Vortex) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	period := int(config[0].(float64))
	offset := int(config[1].(float64))

	indicator.period = period
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")

	// Check if we have enough data points
	if len(highs) < period+offset+1 || len(lows) < period+offset+1 || len(closes) < period+offset+1 {
		return 0.0
	}

	highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least period+1 data points to calculate
	if len(highs) < period+1 || len(lows) < period+1 || len(closes) < period+1 {
		return 0.0
	}

	// Calculate VM+ and VM- for each bar
	var vmPlus = make([]float64, len(closes))
	var vmMinus = make([]float64, len(closes))
	var tr = make([]float64, len(closes))

	for i := 1; i < len(closes); i++ {
		// VM+ = |Current High - Prior Low|
		vmPlus[i] = math.Abs(highs[i] - lows[i-1])

		// VM- = |Current Low - Prior High|
		vmMinus[i] = math.Abs(lows[i] - highs[i-1])

		// Calculate True Range
		tr[i] = max(
			highs[i]-lows[i],
			max(
				math.Abs(highs[i]-closes[i-1]),
				math.Abs(lows[i]-closes[i-1]),
			),
		)
	}

	// Calculate Vortex Indicator
	var viPlus = make([]float64, len(closes))
	var viMinus = make([]float64, len(closes))

	for i := period; i < len(closes); i++ {
		// Calculate sum of VM+ for period
		var sumVMPlus float64
		for j := i - period + 1; j <= i; j++ {
			sumVMPlus += vmPlus[j]
		}

		// Calculate sum of VM- for period
		var sumVMMinus float64
		for j := i - period + 1; j <= i; j++ {
			sumVMMinus += vmMinus[j]
		}

		// Calculate sum of TR for period
		var sumTR float64
		for j := i - period + 1; j <= i; j++ {
			sumTR += tr[j]
		}

		// Calculate VI+ and VI-
		viPlus[i] = sumVMPlus / sumTR
		viMinus[i] = sumVMMinus / sumTR
	}

	// Return the last calculated VI+ value
	return util.Round(viPlus[len(viPlus)-1])
}

func (indicator *Vortex) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 14 // default period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// We need period+1 data points (due to the use of previous values) + offset
	return int32(indicator.period + 1 + indicator.offset)
}
