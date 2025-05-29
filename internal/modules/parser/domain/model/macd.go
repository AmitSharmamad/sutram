package parser

import (
	util "algo/internal/util/math"
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
	return int64(2*slowPeriod + signalPeriod + offset)
}

// CalculateAll calculates MACD values for all data points
func (indicator *MACD) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	fastPeriod := int(config[0].(float64))   // default: 12
	slowPeriod := int(config[1].(float64))   // default: 26
	signalPeriod := int(config[2].(float64)) // default: 9
	// macdType := config[3].(string)           // default: "STANDARD" - handled by Calculate method
	offset := int(config[4].(float64)) // default: 0

	indicator.fastPeriod = fastPeriod
	indicator.slowPeriod = slowPeriod
	indicator.signalPeriod = signalPeriod
	indicator.offset = offset

	stockData := *ctx.Data
	var results []IndicatorValues

	// Need enough data for calculation
	minDataPoints := slowPeriod + signalPeriod + offset
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

		// Calculate MACD value for this point
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
