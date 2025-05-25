package parser

import (
	util "algo/internal/util/math"
)

type MFI struct {
	period   int
	maType   string
	maPeriod int
	offset   int
}

func (indicator *MFI) SliceData(prices Prices) Prices {
	// We need enough data to calculate MFI properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// Helper method to slice all required data
func (indicator *MFI) sliceAllData(highs, lows, closes, volumes Prices) (Prices, Prices, Prices, Prices) {
	// We need enough data to calculate MFI properly
	start := max(0, len(closes)-(indicator.offset+2*indicator.period))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end], volumes[start:end]
}

func (indicator *MFI) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	period := int(config[0].(float64))
	maType := config[1].(string)
	maPeriod := int(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.period = period
	indicator.maType = maType
	indicator.maPeriod = maPeriod
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")
	volumes := GetPrices(*ctx.Data, "volume")

	// Check if we have enough data points
	if len(highs) < period+offset+1 || len(lows) < period+offset+1 ||
		len(closes) < period+offset+1 || len(volumes) < period+offset+1 {
		return 0.0
	}

	highs, lows, closes, volumes = indicator.sliceAllData(highs, lows, closes, volumes)

	// Need at least period+1 data points to calculate
	if len(highs) < period+1 || len(lows) < period+1 ||
		len(closes) < period+1 || len(volumes) < period+1 {
		return 0.0
	}

	// Calculate typical price
	var typicalPrices = make([]float64, len(closes))
	for i := 0; i < len(closes); i++ {
		typicalPrices[i] = (highs[i] + lows[i] + closes[i]) / 3.0
	}

	// Calculate raw money flow
	var rawMoneyFlow = make([]float64, len(closes))
	for i := 0; i < len(closes); i++ {
		rawMoneyFlow[i] = typicalPrices[i] * volumes[i]
	}

	// Calculate positive and negative money flow
	var positiveMoneyFlow = make([]float64, len(closes))
	var negativeMoneyFlow = make([]float64, len(closes))

	for i := 1; i < len(closes); i++ {
		if typicalPrices[i] > typicalPrices[i-1] {
			positiveMoneyFlow[i] = rawMoneyFlow[i]
			negativeMoneyFlow[i] = 0.0
		} else if typicalPrices[i] < typicalPrices[i-1] {
			positiveMoneyFlow[i] = 0.0
			negativeMoneyFlow[i] = rawMoneyFlow[i]
		} else {
			positiveMoneyFlow[i] = 0.0
			negativeMoneyFlow[i] = 0.0
		}
	}

	// Calculate money flow ratio and MFI
	var mfi = make([]float64, len(closes))

	for i := period; i < len(closes); i++ {
		var sumPositiveFlow float64
		var sumNegativeFlow float64

		for j := i - period + 1; j <= i; j++ {
			sumPositiveFlow += positiveMoneyFlow[j]
			sumNegativeFlow += negativeMoneyFlow[j]
		}

		// Avoid division by zero
		if sumNegativeFlow == 0 {
			mfi[i] = 100.0
		} else {
			moneyFlowRatio := sumPositiveFlow / sumNegativeFlow
			mfi[i] = 100.0 - (100.0 / (1.0 + moneyFlowRatio))
		}
	}

	// Apply moving average smoothing if maPeriod > 0
	if indicator.maPeriod > 0 {
		// Extract valid MFI values (skip initial zeros)
		validMfi := mfi[period:]

		// Apply the specified moving average type
		smoothedMfi := util.CalculateMovingAverage(validMfi, indicator.maType, indicator.maPeriod)

		// Return the last value from the smoothed MFI
		if len(smoothedMfi) > 0 {
			return util.Round(smoothedMfi[len(smoothedMfi)-1])
		}
	}

	// Return the last calculated MFI value if no smoothing applied
	return util.Round(mfi[len(mfi)-1])
}

func (indicator *MFI) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 20 // default period
	}

	if t, ok := config["type"].(string); ok {
		indicator.maType = t
	} else {
		indicator.maType = "SIMPLE" // default type
	}

	if mp, ok := config["maPeriod"].(float64); ok {
		indicator.maPeriod = int(mp)
	} else {
		indicator.maPeriod = 20 // default MA period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// We need period+1 data points (due to the use of previous values) + offset
	return int32(indicator.period + 1 + indicator.offset)
}
