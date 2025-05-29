package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"time"
)

type UltimateOscillator struct {
	length1 int
	length2 int
	length3 int
	offset  int
}

func (indicator *UltimateOscillator) SliceData(prices Prices) Prices {
	// We need enough data to calculate Ultimate Oscillator properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.length3))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for UltimateOscillator indicator based on its configuration
func (indicator *UltimateOscillator) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 4 {
		return "", fmt.Errorf("UltimateOscillator requires 4 parameters: length1, length2, length3, offset")
	}

	// Extract length1 parameter
	length1Float, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("UltimateOscillator length1 parameter must be a number")
	}
	length1 := int(length1Float)

	// Extract length2 parameter
	length2Float, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("UltimateOscillator length2 parameter must be a number")
	}
	length2 := int(length2Float)

	// Extract length3 parameter
	length3Float, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("UltimateOscillator length3 parameter must be a number")
	}
	length3 := int(length3Float)

	// Extract offset parameter
	offsetFloat, ok := config[3].(float64)
	if !ok {
		return "", fmt.Errorf("UltimateOscillator offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	baseKey := fmt.Sprintf("ULTIMATE_%d_%d_%d", length1, length2, length3)

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// Helper method to slice all required data
func (indicator *UltimateOscillator) sliceAllData(highs, lows, closes Prices) (Prices, Prices, Prices) {
	// We need enough data to calculate Ultimate Oscillator properly
	start := max(0, len(closes)-(indicator.offset+2*indicator.length3))
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
	length1 := int(config[0].(float64))
	length2 := int(config[1].(float64))
	length3 := int(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.length1 = length1
	indicator.length2 = length2
	indicator.length3 = length3
	indicator.offset = offset

	// Standard weights for Ultimate Oscillator calculation
	// Using 4:2:1 ratio as per traditional implementation
	weightShort := 4.0
	weightMid := 2.0
	weightLong := 1.0

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")

	// Check if we have enough data points
	if len(highs) < length3+offset+1 || len(lows) < length3+offset+1 || len(closes) < length3+offset+1 {
		return 0.0
	}

	highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least length3+1 data points to calculate
	if len(highs) < length3+1 || len(lows) < length3+1 || len(closes) < length3+1 {
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

	// Calculate averages for each period
	var sumBP1, sumTR1 float64
	var sumBP2, sumTR2 float64
	var sumBP3, sumTR3 float64

	// Calculate UO
	var uo = make([]float64, len(closes))

	for i := indicator.length3; i < len(closes); i++ {
		// Reset sums for each period
		sumBP1 = 0
		sumTR1 = 0
		sumBP2 = 0
		sumTR2 = 0
		sumBP3 = 0
		sumTR3 = 0

		// Calculate short period sums
		for j := i - indicator.length1 + 1; j <= i; j++ {
			sumBP1 += bp[j]
			sumTR1 += tr[j]
		}

		// Calculate mid period sums
		for j := i - indicator.length2 + 1; j <= i; j++ {
			sumBP2 += bp[j]
			sumTR2 += tr[j]
		}

		// Calculate long period sums
		for j := i - indicator.length3 + 1; j <= i; j++ {
			sumBP3 += bp[j]
			sumTR3 += tr[j]
		}

		// Calculate averages
		var avg1, avg2, avg3 float64

		// Avoid division by zero
		if sumTR1 > 0 {
			avg1 = sumBP1 / sumTR1
		}

		if sumTR2 > 0 {
			avg2 = sumBP2 / sumTR2
		}

		if sumTR3 > 0 {
			avg3 = sumBP3 / sumTR3
		}

		// Calculate Ultimate Oscillator
		totalWeight := weightShort + weightMid + weightLong
		uo[i] = 100 * ((weightShort*avg1 + weightMid*avg2 + weightLong*avg3) / totalWeight)
	}

	// Return the last calculated Ultimate Oscillator value
	return util.Round(uo[len(uo)-1])
}

func (indicator *UltimateOscillator) GetDataRange(config map[string]any) int32 {
	if p, ok := config["length1"].(float64); ok {
		indicator.length1 = int(p)
	} else {
		indicator.length1 = 7 // default short period
	}

	if p, ok := config["length2"].(float64); ok {
		indicator.length2 = int(p)
	} else {
		indicator.length2 = 14 // default mid period
	}

	if p, ok := config["length3"].(float64); ok {
		indicator.length3 = int(p)
	} else {
		indicator.length3 = 28 // default long period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// We need length3+1 data points (due to the use of previous values) + offset
	return int32(indicator.length3 + 1 + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for Ultimate Oscillator calculation
// config[0] = length1 (float64), config[1] = length2 (float64), config[2] = length3 (float64), config[3] = offset (float64)
func (indicator *UltimateOscillator) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is length1 (float64) - not needed since length3 is largest
	// config[1] is length2 (float64) - not needed since length3 is largest
	length3 := int(config[2].(float64)) // default: 28 (longest period)
	offset := int(config[3].(float64))  // default: 0

	// For Ultimate Oscillator calculation, we need:
	// - length3 data points for longest period calculation
	// - 1 additional for buying pressure and true range calculation
	// - Additional buffer for stability (using 2*length3 like SliceData)
	// - Plus any offset
	return int64(2*length3 + 1 + offset)
}

// CalculateAll calculates Ultimate Oscillator values for all data points
func (indicator *UltimateOscillator) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	length1 := int(config[0].(float64))
	length2 := int(config[1].(float64))
	length3 := int(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.length1 = length1
	indicator.length2 = length2
	indicator.length3 = length3
	indicator.offset = offset

	stockData := *ctx.Data
	var results []IndicatorValues

	// Need enough data for calculation
	minDataPoints := length3 + offset + 1
	if len(stockData) < minDataPoints {
		return []IndicatorValues{}
	}

	// Calculate for each valid data point
	for i := minDataPoints - 1; i < len(stockData); i++ {
		// Create a temporary context with data up to this point
		tempData := stockData[:i+1]
		tempCtx := &EvaluationContext{
			Data:     &tempData,
			Strategy: ctx.Strategy,
			Results:  make(map[int64]map[string]interface{}),
		}

		// Calculate Ultimate Oscillator value for this point
		value := indicator.Calculate(config, tempCtx)

		results = append(results, IndicatorValues{
			Timestamp: stockData[i].Timestamp,
			Value:     value,
		})
	}

	// Filter results by startDate
	var filteredResults []IndicatorValues
	for _, result := range results {
		if !result.Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults
}
