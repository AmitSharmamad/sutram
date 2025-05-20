package parser

import (
	util "algo/internal/util/math"
)

type VWAP struct {
	period int
	offset int
}

func (indicator *VWAP) SliceData(prices Prices) Prices {
	// We need enough data to calculate VWAP properly
	start := max(0, len(prices)-(indicator.offset+indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// Helper method to slice all required data
func (indicator *VWAP) sliceAllData(highs, lows, closes, volumes Prices) (Prices, Prices, Prices, Prices) {
	// We need enough data to calculate VWAP properly
	start := max(0, len(closes)-(indicator.offset+indicator.period))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end], volumes[start:end]
}

func (indicator *VWAP) Calculate(
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
	volumes := GetPrices(*ctx.Data, "volume")

	// Check if we have enough data points
	if len(highs) < period+offset || len(lows) < period+offset ||
		len(closes) < period+offset || len(volumes) < period+offset {
		return 0.0
	}

	highs, lows, closes, volumes = indicator.sliceAllData(highs, lows, closes, volumes)

	// Need at least period data points to calculate
	if len(highs) < period || len(lows) < period ||
		len(closes) < period || len(volumes) < period {
		return 0.0
	}

	// Calculate typical price
	var typicalPrices = make([]float64, len(closes))
	for i := 0; i < len(closes); i++ {
		typicalPrices[i] = (highs[i] + lows[i] + closes[i]) / 3.0
	}

	// Calculate VWAP
	var cumulativeTPV float64 // Typical Price * Volume
	var cumulativeVolume float64
	var vwap = make([]float64, len(closes))

	// Calculate for the lookback period
	for i := len(closes) - period; i < len(closes); i++ {
		cumulativeTPV += typicalPrices[i] * volumes[i]
		cumulativeVolume += volumes[i]

		// Avoid division by zero
		if cumulativeVolume > 0 {
			vwap[i] = cumulativeTPV / cumulativeVolume
		} else {
			vwap[i] = 0.0
		}
	}

	// Return the last calculated VWAP value
	return util.Round(vwap[len(vwap)-1])
}

func (indicator *VWAP) GetDataRange(config map[string]any) int32 {
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

	return int32(indicator.period + indicator.offset)
}
