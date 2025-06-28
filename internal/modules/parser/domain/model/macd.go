package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
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

	// Use go-talib MACD calculation
	macdLine, signalLine, histogram := talib.Macd(prices, fastPeriod, slowPeriod, signalPeriod)

	// Get the value at the correct offset position
	targetIndex := len(macdLine) - 1 - offset
	if targetIndex < 0 || targetIndex >= len(macdLine) {
		return 0.0
	}

	// Return the appropriate value based on macdType
	switch macdType {
	case "HISTOGRAM":
		if !isNaN(histogram[targetIndex]) {
			return util.Round(histogram[targetIndex])
		}
	case "SIGNAL":
		if !isNaN(signalLine[targetIndex]) {
			return util.Round(signalLine[targetIndex])
		}
	default: // "STANDARD" or any other value
		if !isNaN(macdLine[targetIndex]) {
			return util.Round(macdLine[targetIndex])
		}
	}

	// Fallback return in case of any issues
	return 0.0
}

// isNaN checks if a float64 value is NaN
func isNaN(f float64) bool {
	return f != f
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

	// Use go-talib MACD calculation for all values
	macdLine, signalLine, histogram := talib.Macd(closes, fastPeriod, slowPeriod, signalPeriod)

	var filteredResults []IndicatorValues

	// Process MACD values and apply offset
	for i := 0; i < len(macdLine); i++ {
		// Skip NaN values (initial values where MACD can't be calculated)
		var value float64
		var isValid bool

		// Get the appropriate value based on macdType
		switch macdType {
		case "HISTOGRAM":
			if !isNaN(histogram[i]) {
				value = util.Round(histogram[i])
				isValid = true
			}
		case "SIGNAL":
			if !isNaN(signalLine[i]) {
				value = util.Round(signalLine[i])
				isValid = true
			}
		default: // "STANDARD" or any other value
			if !isNaN(macdLine[i]) {
				value = util.Round(macdLine[i])
				isValid = true
			}
		}

		if !isValid {
			continue
		}

		// Calculate the corresponding data index considering offset
		dataIndex := i + offset
		if dataIndex >= len(stockData) {
			break
		}

		// Only include values from startDate onwards
		if !stockData[dataIndex].Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     value,
			})
		}
	}

	return filteredResults
}
