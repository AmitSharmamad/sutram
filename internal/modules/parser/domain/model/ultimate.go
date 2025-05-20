package parser

import (
	util "algo/internal/util/math"
)

type UltimateOscillator struct {
	shortPeriod int
	midPeriod   int
	longPeriod  int
	weightShort float64
	weightMid   float64
	weightLong  float64
	offset      int
}

func (indicator *UltimateOscillator) SliceData(prices Prices) Prices {
	// We need enough data to calculate Ultimate Oscillator properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.longPeriod))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// Helper method to slice all required data
func (indicator *UltimateOscillator) sliceAllData(highs, lows, closes Prices) (Prices, Prices, Prices) {
	// We need enough data to calculate Ultimate Oscillator properly
	start := max(0, len(closes)-(indicator.offset+2*indicator.longPeriod))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end]
}

func (indicator *UltimateOscillator) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	shortPeriod := int(config[0].(float64))
	midPeriod := int(config[1].(float64))
	longPeriod := int(config[2].(float64))
	weightShort := float64(config[3].(float64))
	weightMid := float64(config[4].(float64))
	weightLong := float64(config[5].(float64))
	offset := int(config[6].(float64))

	indicator.shortPeriod = shortPeriod
	indicator.midPeriod = midPeriod
	indicator.longPeriod = longPeriod
	indicator.weightShort = weightShort
	indicator.weightMid = weightMid
	indicator.weightLong = weightLong
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")

	// Check if we have enough data points
	if len(highs) < longPeriod+offset+1 || len(lows) < longPeriod+offset+1 || len(closes) < longPeriod+offset+1 {
		return 0.0
	}

	highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least longPeriod+1 data points to calculate
	if len(highs) < longPeriod+1 || len(lows) < longPeriod+1 || len(closes) < longPeriod+1 {
		return 0.0
	}

	// Calculate buying pressure (BP) and true range (TR)
	var bp = make([]float64, len(closes))
	var tr = make([]float64, len(closes))

	for i := 1; i < len(closes); i++ {
		// Buying Pressure = Close - Minimum(Low, Previous Close)
		bp[i] = closes[i] - min(lows[i], closes[i-1])

		// True Range = Maximum(High, Previous Close) - Minimum(Low, Previous Close)
		tr[i] = max(highs[i], closes[i-1]) - min(lows[i], closes[i-1])
	}

	// Calculate Average7, Average14, and Average28
	var sumBPShort, sumTRShort float64
	var sumBPMid, sumTRMid float64
	var sumBPLong, sumTRLong float64

	// Calculate UO
	var uo = make([]float64, len(closes))

	for i := indicator.longPeriod; i < len(closes); i++ {
		// Reset sums for each period
		sumBPShort = 0
		sumTRShort = 0
		sumBPMid = 0
		sumTRMid = 0
		sumBPLong = 0
		sumTRLong = 0

		// Calculate short period sums
		for j := i - indicator.shortPeriod + 1; j <= i; j++ {
			sumBPShort += bp[j]
			sumTRShort += tr[j]
		}

		// Calculate mid period sums
		for j := i - indicator.midPeriod + 1; j <= i; j++ {
			sumBPMid += bp[j]
			sumTRMid += tr[j]
		}

		// Calculate long period sums
		for j := i - indicator.longPeriod + 1; j <= i; j++ {
			sumBPLong += bp[j]
			sumTRLong += tr[j]
		}

		// Calculate averages
		var average7, average14, average28 float64

		// Avoid division by zero
		if sumTRShort > 0 {
			average7 = sumBPShort / sumTRShort
		}

		if sumTRMid > 0 {
			average14 = sumBPMid / sumTRMid
		}

		if sumTRLong > 0 {
			average28 = sumBPLong / sumTRLong
		}

		// Calculate Ultimate Oscillator
		totalWeight := indicator.weightShort + indicator.weightMid + indicator.weightLong
		uo[i] = 100 * ((indicator.weightShort*average7 +
			indicator.weightMid*average14 +
			indicator.weightLong*average28) / totalWeight)
	}

	// Return the last calculated Ultimate Oscillator value
	return util.Round(uo[len(uo)-1])
}

func (indicator *UltimateOscillator) GetDataRange(config map[string]any) int32 {
	if p, ok := config["shortPeriod"].(float64); ok {
		indicator.shortPeriod = int(p)
	} else {
		indicator.shortPeriod = 7 // default short period
	}

	if p, ok := config["midPeriod"].(float64); ok {
		indicator.midPeriod = int(p)
	} else {
		indicator.midPeriod = 14 // default mid period
	}

	if p, ok := config["longPeriod"].(float64); ok {
		indicator.longPeriod = int(p)
	} else {
		indicator.longPeriod = 28 // default long period
	}

	if w, ok := config["weightShort"].(float64); ok {
		indicator.weightShort = w
	} else {
		indicator.weightShort = 4.0 // default weight for short period
	}

	if w, ok := config["weightMid"].(float64); ok {
		indicator.weightMid = w
	} else {
		indicator.weightMid = 2.0 // default weight for mid period
	}

	if w, ok := config["weightLong"].(float64); ok {
		indicator.weightLong = w
	} else {
		indicator.weightLong = 1.0 // default weight for long period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// We need longPeriod+1 data points (due to the use of previous values) + offset
	return int32(indicator.longPeriod + 1 + indicator.offset)
}
