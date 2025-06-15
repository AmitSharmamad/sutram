package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"
)

type MACD struct {
	fastPeriod   int
	slowPeriod   int
	signalPeriod int
	offset       int
}

func (indicator *MACD) SliceData(prices Prices) Prices {
	// We need enough data to calculate MACD properly
	// For MACD calculation we need at least 2*slowPeriod data points
	maxPeriod := max(indicator.slowPeriod, indicator.fastPeriod)
	start := max(0, len(prices)-(indicator.offset+2*maxPeriod+indicator.signalPeriod))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for MACD indicator based on its configuration
func (indicator *MACD) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 5 {
		return "", fmt.Errorf("MACD requires 5 parameters: fastPeriod, slowPeriod, signalPeriod, macdType, offset")
	}

	// Extract fastPeriod parameter
	fastPeriodFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("MACD fastPeriod parameter must be a number")
	}
	fastPeriod := int(fastPeriodFloat)

	// Extract slowPeriod parameter
	slowPeriodFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("MACD slowPeriod parameter must be a number")
	}
	slowPeriod := int(slowPeriodFloat)

	// Extract signalPeriod parameter
	signalPeriodFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("MACD signalPeriod parameter must be a number")
	}
	signalPeriod := int(signalPeriodFloat)

	// Extract macdType parameter
	macdType, ok := config[3].(string)
	if !ok {
		return "", fmt.Errorf("MACD macdType parameter must be a string")
	}

	// Extract offset parameter
	offsetFloat, ok := config[4].(float64)
	if !ok {
		return "", fmt.Errorf("MACD offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include type and offset only if non-default
	baseKey := fmt.Sprintf("MACD_%d_%d_%d", fastPeriod, slowPeriod, signalPeriod)

	if strings.ToUpper(macdType) != "STANDARD" {
		baseKey += "_" + strings.ToUpper(macdType)
	}

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

func (indicator *MACD) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	fastPeriod := int(config[0].(float64))   // default: 12
	slowPeriod := int(config[1].(float64))   // default: 26
	signalPeriod := int(config[2].(float64)) // default: 9
	macdType := config[3].(string)           // default: "STANDARD"
	offset := int(config[4].(float64))       // default: 0

	indicator.fastPeriod = fastPeriod
	indicator.slowPeriod = slowPeriod
	indicator.signalPeriod = signalPeriod
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, "close") // MACD typically uses close prices

	// Check if we have enough data points
	minDataPoints := indicator.slowPeriod + indicator.signalPeriod + indicator.offset
	if len(prices) < minDataPoints {
		return 0.0
	}

	prices = indicator.SliceData(prices)

	// Need at least slowPeriod data points to calculate
	if len(prices) < indicator.slowPeriod {
		return 0.0
	}

	// Calculate fast EMA
	fastMultiplier := 2.0 / float64(indicator.fastPeriod+1)
	var fastSMA float64
	for i := 0; i < indicator.fastPeriod; i++ {
		fastSMA += prices[i]
	}
	fastSMA /= float64(indicator.fastPeriod)

	fastEMA := fastSMA
	var fastEMAs = make([]float64, len(prices))
	fastEMAs[indicator.fastPeriod-1] = fastEMA

	// Calculate slow EMA
	slowMultiplier := 2.0 / float64(indicator.slowPeriod+1)
	var slowSMA float64
	for i := 0; i < indicator.slowPeriod; i++ {
		slowSMA += prices[i]
	}
	slowSMA /= float64(indicator.slowPeriod)

	slowEMA := slowSMA
	var slowEMAs = make([]float64, len(prices))
	slowEMAs[indicator.slowPeriod-1] = slowEMA

	// Calculate fast and slow EMAs
	for i := indicator.fastPeriod; i < len(prices); i++ {
		fastEMA = prices[i]*fastMultiplier + fastEMA*(1-fastMultiplier)
		fastEMAs[i] = fastEMA
	}

	for i := indicator.slowPeriod; i < len(prices); i++ {
		slowEMA = prices[i]*slowMultiplier + slowEMA*(1-slowMultiplier)
		slowEMAs[i] = slowEMA
	}

	// Calculate MACD line
	var macdLine = make([]float64, len(prices))
	for i := indicator.slowPeriod - 1; i < len(prices); i++ {
		macdLine[i] = fastEMAs[i] - slowEMAs[i]
	}

	// Calculate signal line using EMA of MACD line
	var signalLine = make([]float64, len(prices))
	signalMultiplier := 2.0 / float64(indicator.signalPeriod+1)

	// Calculate initial signal line value as SMA of MACD
	startIndex := indicator.slowPeriod + indicator.signalPeriod - 2
	var signalSMA float64
	for i := 0; i < indicator.signalPeriod; i++ {
		signalSMA += macdLine[indicator.slowPeriod-1+i]
	}
	signalSMA /= float64(indicator.signalPeriod)

	signalLine[startIndex] = signalSMA
	var signalEMA = signalSMA

	// Calculate subsequent signal line values
	for i := startIndex + 1; i < len(prices); i++ {
		signalEMA = macdLine[i]*signalMultiplier + signalEMA*(1-signalMultiplier)
		signalLine[i] = signalEMA
	}

	// Calculate histogram (MACD line - signal line)
	var histogram = make([]float64, len(prices))
	for i := startIndex; i < len(prices); i++ {
		histogram[i] = macdLine[i] - signalLine[i]
	}

	// Return the appropriate value based on macdType
	lastIndex := len(macdLine) - 1
	switch macdType {
	case "HISTOGRAM":
		if lastIndex >= startIndex {
			return util.Round(histogram[lastIndex])
		}
	case "SIGNAL":
		if lastIndex >= startIndex {
			return util.Round(signalLine[lastIndex])
		}
	default: // "STANDARD" or any other value
		return util.Round(macdLine[lastIndex])
	}

	// Fallback return in case of any issues
	return 0.0
}

func (indicator *MACD) GetDataRange(config map[string]any) int32 {
	if p, ok := config["fastPeriod"].(float64); ok {
		indicator.fastPeriod = int(p)
	} else {
		indicator.fastPeriod = 12 // default fast period
	}

	if p, ok := config["slowPeriod"].(float64); ok {
		indicator.slowPeriod = int(p)
	} else {
		indicator.slowPeriod = 26 // default slow period
	}

	if p, ok := config["signalPeriod"].(float64); ok {
		indicator.signalPeriod = int(p)
	} else {
		indicator.signalPeriod = 9 // default signal period
	}

	// No need to store macdType in struct as it's only used for return value selection

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// Return the maximum data range required
	return int32(indicator.slowPeriod + indicator.signalPeriod + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for MACD calculation
// config[0] = fastPeriod (float64), config[1] = slowPeriod (float64), config[2] = signalPeriod (float64)
// config[3] = macdType (string), config[4] = offset (float64)
func (indicator *MACD) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is fastPeriod (float64) - not needed for lookback calculation since slowPeriod is larger
	slowPeriod := int(config[1].(float64))   // default: 26
	signalPeriod := int(config[2].(float64)) // default: 9
	// config[3] is macdType (string) - not needed for lookback calculation
	offset := int(config[4].(float64)) // default: 0

	// For MACD calculation, we need:
	// - Enough data for slow EMA (slowPeriod)
	// - Additional data for signal line calculation (signalPeriod)
	// - Additional buffer for EMA stabilization (using 2*slowPeriod like SliceData)
	// - Plus any offset
	return int64(5*slowPeriod + signalPeriod + offset)
}

// CalculateAll calculates MACD values for all data points
// Following Excel flow: Fast EMA -> Slow EMA -> MACD Line -> Signal Line -> Histogram
func (indicator *MACD) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	fastPeriod := int(config[0].(float64))   // default: 12
	slowPeriod := int(config[1].(float64))   // default: 26
	signalPeriod := int(config[2].(float64)) // default: 9
	macdType := config[3].(string)           // default: "STANDARD"
	offset := int(config[4].(float64))       // default: 0

	indicator.fastPeriod = fastPeriod
	indicator.slowPeriod = slowPeriod
	indicator.signalPeriod = signalPeriod
	indicator.offset = offset

	stockData := *ctx.Data
	dataLen := len(stockData)

	// Need enough data for calculation (slow period + signal period)
	if dataLen < slowPeriod+signalPeriod {
		return []IndicatorValues{}
	}

	// Extract close prices
	closes := make([]float64, dataLen)
	for i := 0; i < dataLen; i++ {
		closes[i] = stockData[i].Close
	}

	// Step 1: Calculate 12-Period EMA (Fast EMA)
	fastEMA := make([]float64, dataLen)
	fastMultiplier := 2.0 / float64(fastPeriod+1)

	// Initial SMA for fast period (C13 in Excel)
	if fastPeriod <= dataLen {
		var sum float64
		for i := 0; i < fastPeriod; i++ {
			sum += closes[i]
		}
		fastEMA[fastPeriod-1] = sum / float64(fastPeriod)

		// Calculate EMA for subsequent periods (C14 onwards)
		for i := fastPeriod; i < dataLen; i++ {
			// EMA = Close * multiplier + Previous EMA * (1 - multiplier)
			fastEMA[i] = closes[i]*fastMultiplier + fastEMA[i-1]*(1-fastMultiplier)
		}
	}

	// Step 2: Calculate 26-Period EMA (Slow EMA)
	slowEMA := make([]float64, dataLen)
	slowMultiplier := 2.0 / float64(slowPeriod+1)

	// Initial SMA for slow period (D27 in Excel)
	if slowPeriod <= dataLen {
		var sum float64
		for i := 0; i < slowPeriod; i++ {
			sum += closes[i]
		}
		slowEMA[slowPeriod-1] = sum / float64(slowPeriod)

		// Calculate EMA for subsequent periods (D28 onwards)
		for i := slowPeriod; i < dataLen; i++ {
			// EMA = Close * multiplier + Previous EMA * (1 - multiplier)
			slowEMA[i] = closes[i]*slowMultiplier + slowEMA[i-1]*(1-slowMultiplier)
		}
	}

	// Step 3: Calculate MACD Line (E27 onwards)
	macdLine := make([]float64, dataLen)
	for i := slowPeriod - 1; i < dataLen; i++ {
		// MACD Line = Fast EMA - Slow EMA
		macdLine[i] = fastEMA[i] - slowEMA[i]
	}

	// Step 4: Calculate 9-Period EMA of MACD (Signal Line)
	signalLine := make([]float64, dataLen)
	signalMultiplier := 2.0 / float64(signalPeriod+1)

	// Calculate initial signal line starting index (F35 in Excel)
	signalStartIndex := slowPeriod + signalPeriod - 2
	if signalStartIndex < dataLen {
		// Initial SMA of MACD for signal period
		var sum float64
		for i := 0; i < signalPeriod; i++ {
			sum += macdLine[slowPeriod-1+i]
		}
		signalLine[signalStartIndex] = sum / float64(signalPeriod)

		// Calculate EMA for subsequent periods (F36 onwards)
		for i := signalStartIndex + 1; i < dataLen; i++ {
			// Signal EMA = MACD * multiplier + Previous Signal * (1 - multiplier)
			signalLine[i] = macdLine[i]*signalMultiplier + signalLine[i-1]*(1-signalMultiplier)
		}
	}

	// Step 5: Calculate MACD Histogram (G35 onwards)
	histogram := make([]float64, dataLen)
	for i := signalStartIndex; i < dataLen; i++ {
		// Histogram = MACD Line - Signal Line
		histogram[i] = macdLine[i] - signalLine[i]
	}

	// Build results starting from when all calculations are valid
	var results []IndicatorValues
	for i := signalStartIndex; i < dataLen; i++ {
		// Apply offset if specified
		adjustedIndex := i + offset
		if adjustedIndex >= 0 && adjustedIndex < dataLen {
			var value float64
			// Return the appropriate value based on macdType
			switch macdType {
			case "HISTOGRAM":
				value = util.Round(histogram[i])
			case "SIGNAL":
				value = util.Round(signalLine[i])
			default: // "STANDARD" or any other value
				value = util.Round(macdLine[i])
			}

			results = append(results, IndicatorValues{
				Timestamp: stockData[adjustedIndex].Timestamp,
				Value:     value,
			})
		}
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
