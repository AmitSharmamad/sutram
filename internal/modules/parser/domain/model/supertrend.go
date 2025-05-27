package parser

import (
	util "algo/internal/util/math"
	"math"
)

type Supertrend struct {
	period     int
	multiplier float64
	offset     int
}

func (indicator *Supertrend) SliceData(prices Prices) Prices {
	// We need enough data to calculate Supertrend properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// Internal helper method to slice all required data
func (indicator *Supertrend) sliceAllData(highs, lows, closes Prices) (Prices, Prices, Prices) {
	// We need enough data to calculate Supertrend properly
	start := max(0, len(closes)-(indicator.offset+2*indicator.period))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end]
}

func (indicator *Supertrend) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	period := int(config[0].(float64))
	multiplier := float64(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.multiplier = multiplier
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")

	// Check if we have enough data points
	if len(highs) < period+offset || len(lows) < period+offset || len(closes) < period+offset {
		return 0.0
	}

	highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least period+1 data points to calculate Supertrend properly
	if len(highs) <= period || len(lows) <= period || len(closes) <= period {
		return 0.0
	}

	// Calculate ATR (Average True Range)
	var atr = make([]float64, len(closes))
	var tr = make([]float64, len(closes))

	// Calculate True Range
	for i := 1; i < len(closes); i++ {
		tr[i] = max(
			highs[i]-lows[i],
			max(
				math.Abs(highs[i]-closes[i-1]),
				math.Abs(lows[i]-closes[i-1]),
			),
		)
	}

	// Calculate initial ATR (simple average of TR for the first 'period' days)
	var sumTR float64
	for i := 1; i < period+1 && i < len(closes); i++ {
		sumTR += tr[i]
	}
	if period < len(closes) {
		atr[period] = sumTR / float64(period)
	}

	// Calculate ATR using the smoothing technique
	for i := period + 1; i < len(closes); i++ {
		atr[i] = (atr[i-1]*(float64(period)-1) + tr[i]) / float64(period)
	}

	// Calculate upper and lower bands
	var upperBand = make([]float64, len(closes))
	var lowerBand = make([]float64, len(closes))
	var supertrend = make([]float64, len(closes))
	var trend = make([]int, len(closes)) // 1 for up, -1 for down

	// Calculate initial bands
	for i := period; i < len(closes); i++ {
		// Ensure we have valid ATR data
		if i >= len(atr) || atr[i] == 0 {
			continue
		}

		basicUpperBand := (highs[i]+lows[i])/2 + indicator.multiplier*atr[i]
		basicLowerBand := (highs[i]+lows[i])/2 - indicator.multiplier*atr[i]

		// Set initial values
		if i == period {
			upperBand[i] = basicUpperBand
			lowerBand[i] = basicLowerBand
			if closes[i] <= upperBand[i] {
				supertrend[i] = lowerBand[i]
				trend[i] = -1
			} else {
				supertrend[i] = upperBand[i]
				trend[i] = 1
			}
			continue
		}

		// Adjust upper band
		if basicUpperBand < upperBand[i-1] || closes[i-1] > upperBand[i-1] {
			upperBand[i] = basicUpperBand
		} else {
			upperBand[i] = upperBand[i-1]
		}

		// Adjust lower band
		if basicLowerBand > lowerBand[i-1] || closes[i-1] < lowerBand[i-1] {
			lowerBand[i] = basicLowerBand
		} else {
			lowerBand[i] = lowerBand[i-1]
		}

		// Determine trend and supertrend value
		if supertrend[i-1] == upperBand[i-1] {
			if closes[i] <= upperBand[i] {
				supertrend[i] = upperBand[i]
				trend[i] = -1
			} else {
				supertrend[i] = lowerBand[i]
				trend[i] = 1
			}
		} else { // supertrend[i-1] == lowerBand[i-1]
			if closes[i] >= lowerBand[i] {
				supertrend[i] = lowerBand[i]
				trend[i] = 1
			} else {
				supertrend[i] = upperBand[i]
				trend[i] = -1
			}
		}
	}

	// Return the last calculated supertrend value
	if len(supertrend) == 0 {
		return 0.0
	}
	return util.Round(supertrend[len(supertrend)-1])
}

func (indicator *Supertrend) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 7 // default period from indicators.json
	}

	if m, ok := config["multiplier"].(float64); ok {
		indicator.multiplier = m
	} else {
		indicator.multiplier = 3.0 // default multiplier
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.period + indicator.offset)
}
