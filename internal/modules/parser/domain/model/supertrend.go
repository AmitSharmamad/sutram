package parser

import (
	util "algo/internal/util/math"
	"fmt"
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

	// Create unique indicator keys for different values we need to store
	indicatorKey := fmt.Sprintf("SUPERTREND_%d_%.2f", period, multiplier)
	atrKey := fmt.Sprintf("ATR_%d", period)
	upperBandKey := fmt.Sprintf("UPPERBAND_%d_%.2f", period, multiplier)
	lowerBandKey := fmt.Sprintf("LOWERBAND_%d_%.2f", period, multiplier)
	trendKey := fmt.Sprintf("TREND_%d_%.2f", period, multiplier)

	// Get current timestamp
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()

	// Try to get previous values from Results for optimization
	if prevATR, atrExists := getPreviousValueFromResults(ctx, atrKey, currentTimestamp-1); atrExists {
		if prevUpper, upperExists := getPreviousValueFromResults(ctx, upperBandKey, currentTimestamp-1); upperExists {
			if prevLower, lowerExists := getPreviousValueFromResults(ctx, lowerBandKey, currentTimestamp-1); lowerExists {
				if prevTrend, trendExists := getPreviousValueFromResults(ctx, trendKey, currentTimestamp-1); trendExists {
					if prevSupertrend, stExists := getPreviousValueFromResults(ctx, indicatorKey, currentTimestamp-1); stExists {
						// Use optimized calculation with previous values
						return indicator.calculateOptimized(
							ctx, highs, lows, closes, currentTimestamp,
							prevATR.(float64), prevUpper.(float64), prevLower.(float64),
							prevTrend.(int), prevSupertrend.(float64),
							indicatorKey, atrKey, upperBandKey, lowerBandKey, trendKey,
						)
					}
				}
			}
		}
	}

	// If no previous values exist, calculate from scratch
	return indicator.calculateFromScratch(ctx, highs, lows, closes, currentTimestamp,
		indicatorKey, atrKey, upperBandKey, lowerBandKey, trendKey)
}

// calculateOptimized uses previous values for faster calculation
func (indicator *Supertrend) calculateOptimized(
	ctx *EvaluationContext, highs, lows, closes []float64, currentTimestamp int64,
	prevATR, prevUpper, prevLower float64, prevTrend int, prevSupertrend float64,
	indicatorKey, atrKey, upperBandKey, lowerBandKey, trendKey string,
) float64 {
	// Get the current period data indices
	currentIdx := len(closes) - 1

	if currentIdx < 1 {
		return 0.0
	}

	// Calculate True Range for current period
	tr := max(
		highs[currentIdx]-lows[currentIdx],
		max(
			math.Abs(highs[currentIdx]-closes[currentIdx-1]),
			math.Abs(lows[currentIdx]-closes[currentIdx-1]),
		),
	)

	// Update ATR using previous ATR
	currentATR := (prevATR*(float64(indicator.period)-1) + tr) / float64(indicator.period)

	// Calculate basic bands
	basicUpperBand := (highs[currentIdx]+lows[currentIdx])/2 + indicator.multiplier*currentATR
	basicLowerBand := (highs[currentIdx]+lows[currentIdx])/2 - indicator.multiplier*currentATR

	// Adjust upper band
	var currentUpper float64
	if basicUpperBand < prevUpper || closes[currentIdx-1] > prevUpper {
		currentUpper = basicUpperBand
	} else {
		currentUpper = prevUpper
	}

	// Adjust lower band
	var currentLower float64
	if basicLowerBand > prevLower || closes[currentIdx-1] < prevLower {
		currentLower = basicLowerBand
	} else {
		currentLower = prevLower
	}

	// Determine trend and supertrend value
	var currentSupertrend float64
	var currentTrend int

	if prevSupertrend == prevUpper {
		if closes[currentIdx] <= currentUpper {
			currentSupertrend = currentUpper
			currentTrend = -1
		} else {
			currentSupertrend = currentLower
			currentTrend = 1
		}
	} else { // prevSupertrend == prevLower
		if closes[currentIdx] >= currentLower {
			currentSupertrend = currentLower
			currentTrend = 1
		} else {
			currentSupertrend = currentUpper
			currentTrend = -1
		}
	}

	// Store all values for future calculations
	insertValueToResults(ctx, atrKey, currentTimestamp, currentATR)
	insertValueToResults(ctx, upperBandKey, currentTimestamp, currentUpper)
	insertValueToResults(ctx, lowerBandKey, currentTimestamp, currentLower)
	insertValueToResults(ctx, trendKey, currentTimestamp, currentTrend)
	insertValueToResults(ctx, indicatorKey, currentTimestamp, currentSupertrend)

	return util.Round(currentSupertrend)
}

// calculateFromScratch performs full calculation when no previous values exist
func (indicator *Supertrend) calculateFromScratch(
	ctx *EvaluationContext, highs, lows, closes []float64, currentTimestamp int64,
	indicatorKey, atrKey, upperBandKey, lowerBandKey, trendKey string,
) float64 {
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
	for i := 1; i < indicator.period+1 && i < len(closes); i++ {
		sumTR += tr[i]
	}
	if indicator.period < len(closes) {
		atr[indicator.period] = sumTR / float64(indicator.period)
	}

	// Calculate ATR using the smoothing technique
	for i := indicator.period + 1; i < len(closes); i++ {
		atr[i] = (atr[i-1]*(float64(indicator.period)-1) + tr[i]) / float64(indicator.period)
	}

	// Calculate upper and lower bands
	var upperBand = make([]float64, len(closes))
	var lowerBand = make([]float64, len(closes))
	var supertrend = make([]float64, len(closes))
	var trend = make([]int, len(closes)) // 1 for up, -1 for down

	// Calculate initial bands
	for i := indicator.period; i < len(closes); i++ {
		// Ensure we have valid ATR data
		if i >= len(atr) || atr[i] == 0 {
			continue
		}

		basicUpperBand := (highs[i]+lows[i])/2 + indicator.multiplier*atr[i]
		basicLowerBand := (highs[i]+lows[i])/2 - indicator.multiplier*atr[i]

		// Set initial values
		if i == indicator.period {
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

	// Store the final calculated values for future optimization
	if len(supertrend) > 0 {
		lastIdx := len(supertrend) - 1
		insertValueToResults(ctx, atrKey, currentTimestamp, atr[lastIdx])
		insertValueToResults(ctx, upperBandKey, currentTimestamp, upperBand[lastIdx])
		insertValueToResults(ctx, lowerBandKey, currentTimestamp, lowerBand[lastIdx])
		insertValueToResults(ctx, trendKey, currentTimestamp, trend[lastIdx])
		insertValueToResults(ctx, indicatorKey, currentTimestamp, supertrend[lastIdx])

		return util.Round(supertrend[lastIdx])
	}

	return 0.0
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

// GetLookBack returns the number of historical data points needed for Supertrend calculation
// config[0] = period (float64), config[1] = multiplier (float64), config[2] = offset (float64)
func (indicator *Supertrend) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	period := int(config[0].(float64)) // default: 10
	// config[1] is multiplier (float64) - not needed for lookback calculation
	offset := int(config[2].(float64)) // default: 0

	// For Supertrend calculation, we need:
	// - Period data points for ATR calculation
	// - 1 additional for true range calculation
	// - Additional buffer for stability (using 2*period like SliceData)
	// - Plus any offset
	return int64(2*period + 1 + offset)
}
