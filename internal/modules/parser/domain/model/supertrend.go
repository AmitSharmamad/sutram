package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"math"
	"time"
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

// generateIndicatorKey generates a unique key for Supertrend indicator based on its configuration
func (indicator *Supertrend) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 3 {
		return "", fmt.Errorf("Supertrend requires 3 parameters: period, multiplier, offset")
	}

	// Extract period parameter
	periodFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("Supertrend period parameter must be a number")
	}
	period := int(periodFloat)

	// Extract multiplier parameter
	multiplier, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("Supertrend multiplier parameter must be a number")
	}

	// Extract offset parameter
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("Supertrend offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	baseKey := fmt.Sprintf("SUPERTREND_%d_%.1f", period, multiplier)

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
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

	highs := GetPrices(*ctx.Data, "HIGH")
	lows := GetPrices(*ctx.Data, "LOW")
	closes := GetPrices(*ctx.Data, "CLOSE")

	// Check if we have enough data points
	if len(highs) < period+offset || len(lows) < period+offset || len(closes) < period+offset {
		return 0.0
	}

	// highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least period+1 data points to calculate Supertrend properly
	// if len(highs) <= period || len(lows) <= period || len(closes) <= period {
	// 	return 0.0
	// }

	// Calculate Supertrend from scratch
	return indicator.calculateSupertrend(highs, lows, closes)
}

// calculateSupertrend performs the complete Supertrend calculation
// This implementation follows the pandas_ta algorithm exactly
// Key improvements made to match Python pandas_ta:
// 1. Uses direction tracking (dir) to determine trend changes
// 2. Applies band adjustments only when direction doesn't change
// 3. Sets supertrend value based on direction (bullish=lowerband, bearish=upperband)
// 4. Initializes direction to bullish (1) at start
func (indicator *Supertrend) calculateSupertrend(highs, lows, closes []float64) float64 {
	m := len(closes)
	if m < indicator.period+1 {
		return 0.0
	}

	// Calculate ATR (Average True Range)
	var atr = make([]float64, m)
	var tr = make([]float64, m)

	// Calculate True Range
	for i := 1; i < m; i++ {
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
	for i := 1; i < indicator.period+1 && i < m; i++ {
		sumTR += tr[i]
	}
	if indicator.period < m {
		atr[indicator.period] = sumTR / float64(indicator.period)
	}

	// Calculate ATR using the smoothing technique
	for i := indicator.period + 1; i < m; i++ {
		atr[i] = (atr[i-1]*(float64(indicator.period)-1) + tr[i]) / float64(indicator.period)
	}

	// Initialize arrays following pandas_ta algorithm
	dir := make([]int, m)
	trend := make([]float64, m)

	// Initialize direction to 1 (bullish) as per pandas_ta
	for i := 0; i < m; i++ {
		dir[i] = 1
	}

	// Calculate HL2 (median price) and bands
	hl2 := make([]float64, m)
	upperband := make([]float64, m)
	lowerband := make([]float64, m)

	for i := 0; i < m; i++ {
		hl2[i] = (highs[i] + lows[i]) / 2.0
		if i >= indicator.period && atr[i] > 0 {
			matr := indicator.multiplier * atr[i]
			upperband[i] = hl2[i] + matr
			lowerband[i] = hl2[i] - matr
		}
	}

	// Main supertrend calculation loop following pandas_ta algorithm
	for i := 1; i < m; i++ {
		if i < indicator.period {
			continue
		}

		// Determine direction based on close vs previous bands
		if closes[i] > upperband[i-1] {
			dir[i] = 1
		} else if closes[i] < lowerband[i-1] {
			dir[i] = -1
		} else {
			dir[i] = dir[i-1]
			// Adjust bands based on direction
			if dir[i] > 0 && lowerband[i] < lowerband[i-1] {
				lowerband[i] = lowerband[i-1]
			}
			if dir[i] < 0 && upperband[i] > upperband[i-1] {
				upperband[i] = upperband[i-1]
			}
		}

		// Set supertrend value based on direction
		if dir[i] > 0 {
			trend[i] = lowerband[i]
		} else {
			trend[i] = upperband[i]
		}
	}

	// Return the final calculated supertrend value
	if m > 0 {
		lastIdx := m - 1
		return util.Round(trend[lastIdx])
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
	return int64(5*period + 1 + offset)
}

// CalculateAll calculates Supertrend values for all data points
// Following Excel flow: TR -> ATR -> Basic Bands -> Supertrend with Trend logic
func (indicator *Supertrend) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	period := int(config[0].(float64))
	multiplier := float64(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.multiplier = multiplier
	indicator.offset = offset

	stockData := *ctx.Data
	dataLen := len(stockData)

	if dataLen < period+1 {
		return []IndicatorValues{}
	}

	highs := make([]float64, dataLen)
	lows := make([]float64, dataLen)
	closes := make([]float64, dataLen)

	for i := 0; i < dataLen; i++ {
		highs[i] = stockData[i].High
		lows[i] = stockData[i].Low
		closes[i] = stockData[i].Close
	}

	tr := make([]float64, dataLen)
	for i := 1; i < dataLen; i++ {
		hl := highs[i] - lows[i]
		hc := math.Abs(highs[i] - closes[i-1])
		lc := math.Abs(lows[i] - closes[i-1])
		tr[i] = math.Max(hl, math.Max(hc, lc))
	}

	atr := make([]float64, dataLen)

	if period < dataLen {
		var sumTR float64
		for i := 1; i <= period; i++ {
			sumTR += tr[i]
		}
		atr[period] = sumTR / float64(period)

		for i := period + 1; i < dataLen; i++ {
			atr[i] = (atr[i-1]*float64(period-1) + tr[i]) / float64(period)
		}
	}

	upperBand := make([]float64, dataLen)
	lowerBand := make([]float64, dataLen)

	for i := 0; i < dataLen; i++ {
		hl2 := (highs[i] + lows[i]) / 2.0
		upperBand[i] = hl2 + multiplier*atr[i]
		lowerBand[i] = hl2 - multiplier*atr[i]
	}

	supertrend := make([]float64, dataLen)
	trend := make([]int, dataLen)

	if period < dataLen {
		if closes[period] < upperBand[period] {
			supertrend[period] = upperBand[period]
			trend[period] = -1 // Downtrend
		} else {
			supertrend[period] = lowerBand[period]
			trend[period] = 1 // Uptrend
		}
	}

	for i := period + 1; i < dataLen; i++ {
		if closes[i-1] > upperBand[i-1] {
			upperBand[i] = math.Min(upperBand[i], upperBand[i-1])
		}
		if closes[i-1] < lowerBand[i-1] {
			lowerBand[i] = math.Max(lowerBand[i], lowerBand[i-1])
		}

		// Immediate band switch on trend flip
		if closes[i] > supertrend[i-1] && trend[i-1] == -1 {
			supertrend[i] = lowerBand[i] // Switch to lower band for uptrend
			trend[i] = 1
		} else if closes[i] < supertrend[i-1] && trend[i-1] == 1 {
			supertrend[i] = upperBand[i] // Switch to upper band for downtrend
			trend[i] = -1
		} else {
			// Continue current trend, adjust within same band
			if trend[i-1] == 1 {
				supertrend[i] = math.Max(lowerBand[i], supertrend[i-1])
			} else {
				supertrend[i] = math.Min(upperBand[i], supertrend[i-1])
			}
			trend[i] = trend[i-1]
		}
	}

	var results []IndicatorValues
	for i := period; i < dataLen; i++ {
		adjustedIndex := i + offset
		if adjustedIndex >= 0 && adjustedIndex < dataLen {
			results = append(results, IndicatorValues{
				Timestamp: stockData[adjustedIndex].Timestamp,
				Value:     supertrend[i],
			})
		}
	}

	var filteredResults []IndicatorValues
	for _, result := range results {
		if !result.Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults
}
