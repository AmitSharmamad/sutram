package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"math"
	"strings"
	"time"
)

type KeltnerChannels struct {
	emaPeriod  int
	atrPeriod  int
	multiplier float64
	offset     int
}

func (indicator *KeltnerChannels) SliceData(prices Prices) Prices {
	// We need enough data to calculate Keltner Channels properly
	maxPeriod := max(indicator.emaPeriod, indicator.atrPeriod)
	start := max(0, len(prices)-(indicator.offset+2*maxPeriod))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for KeltnerChannels indicator based on its configuration
func (indicator *KeltnerChannels) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters (4 required + 1 optional)
	if len(config) < 4 || len(config) > 5 {
		return "", fmt.Errorf("KeltnerChannels requires 4-5 parameters: emaPeriod, bandType, maType, multiplier, [offset]")
	}

	// Extract emaPeriod parameter
	emaPeriodFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("KeltnerChannels emaPeriod parameter must be a number")
	}
	emaPeriod := int(emaPeriodFloat)

	// Extract bandType parameter
	bandType, ok := config[1].(string)
	if !ok {
		return "", fmt.Errorf("KeltnerChannels bandType parameter must be a string")
	}

	// Extract maType parameter
	maType, ok := config[2].(string)
	if !ok {
		return "", fmt.Errorf("KeltnerChannels maType parameter must be a string")
	}

	// Extract multiplier parameter
	multiplier, ok := config[3].(float64)
	if !ok {
		return "", fmt.Errorf("KeltnerChannels multiplier parameter must be a number")
	}

	// Extract optional offset parameter
	offset := 0
	if len(config) > 4 {
		offsetFloat, ok := config[4].(float64)
		if !ok {
			return "", fmt.Errorf("KeltnerChannels offset parameter must be a number")
		}
		offset = int(offsetFloat)
	}

	// ATR period is same as EMA period (as per CalculateAll implementation)
	atrPeriod := emaPeriod

	// Generate the key - include non-default values
	baseKey := fmt.Sprintf("KELTNER_%d_%s_%s_%.1f_%d",
		emaPeriod, strings.ToUpper(bandType), strings.ToUpper(maType), multiplier, atrPeriod)

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// Helper method to slice all required data
func (indicator *KeltnerChannels) sliceAllData(highs, lows, closes Prices) (Prices, Prices, Prices) {
	// We need enough data to calculate Keltner Channels properly
	maxPeriod := max(indicator.emaPeriod, indicator.atrPeriod)
	start := max(0, len(closes)-(indicator.offset+2*maxPeriod))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end]
}

func (indicator *KeltnerChannels) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	// Parse config parameters according to indicators.json
	emaPeriod := int(config[0].(float64))      // default: 50
	bandType := config[1].(string)             // default: "MEDIAN"
	maType := config[2].(string)               // default: "EXPONENTIAL"
	multiplier := float64(config[3].(float64)) // default: 5
	offset := 0

	// Handle optional parameter if present
	if len(config) > 4 {
		offset = int(config[4].(float64))
	}

	indicator.emaPeriod = emaPeriod
	indicator.atrPeriod = emaPeriod // Using same period for ATR calculation
	indicator.multiplier = multiplier
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")

	// Check if we have enough data points
	maxPeriod := emaPeriod
	if len(highs) < maxPeriod+offset+1 || len(lows) < maxPeriod+offset+1 || len(closes) < maxPeriod+offset+1 {
		return 0.0
	}

	highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least maxPeriod+1 data points to calculate
	if len(highs) < maxPeriod+1 || len(lows) < maxPeriod+1 || len(closes) < maxPeriod+1 {
		return 0.0
	}

	// Calculate price series (always using median price for Keltner)
	var prices = make([]float64, len(closes))
	for i := 0; i < len(closes); i++ {
		// For Keltner Channels, we typically use median price (high+low)/2
		prices[i] = (highs[i] + lows[i]) / 2.0
	}

	// Calculate moving average using utility function
	emas := util.CalculateMovingAverage(prices, maType, indicator.emaPeriod)

	// Calculate ATR (Average True Range)
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

	// Calculate ATR using moving average on TR
	trData := tr[1:] // Skip the first zero value
	atrs := util.CalculateMovingAverage(trData, maType, indicator.atrPeriod)

	// Adjust ATR length to match the original data
	fullAtrs := make([]float64, len(closes))
	// Fill in missing values at the beginning
	for i := 0; i < indicator.atrPeriod+1; i++ {
		fullAtrs[i] = 0.0
	}
	// Copy calculated values
	copy(fullAtrs[indicator.atrPeriod+1:], atrs)

	// Calculate upper and lower bands
	var upperBand = make([]float64, len(closes))
	var lowerBand = make([]float64, len(closes))

	// Determine the start index for valid calculations
	startIndex := indicator.emaPeriod
	if len(emas) < len(closes) {
		startIndex += len(closes) - len(emas)
	}

	for i := startIndex; i < len(closes); i++ {
		if i < len(emas) {
			// Calculate bands using MA and ATR
			ema := 0.0
			if i < len(emas) {
				ema = emas[i-startIndex]
			}

			atr := 0.0
			if i < len(fullAtrs) {
				atr = fullAtrs[i]
			}

			upperBand[i] = ema + indicator.multiplier*atr
			lowerBand[i] = ema - indicator.multiplier*atr
		}
	}

	// Determine which band to return based on bandType
	if len(emas) > 0 {
		var result float64
		lastEmaIndex := len(emas) - 1
		lastAtrIndex := min(len(atrs)-1, len(fullAtrs)-1-indicator.atrPeriod-1)

		switch bandType {
		case "TOP":
			if lastAtrIndex >= 0 {
				result = emas[lastEmaIndex] + indicator.multiplier*atrs[lastAtrIndex]
			} else {
				result = emas[lastEmaIndex]
			}
		case "BOTTOM":
			if lastAtrIndex >= 0 {
				result = emas[lastEmaIndex] - indicator.multiplier*atrs[lastAtrIndex]
			} else {
				result = emas[lastEmaIndex]
			}
		default: // "MEDIAN" or any other value
			result = emas[lastEmaIndex]
		}
		return util.Round(result)
	}
	return 0.0
}

func (indicator *KeltnerChannels) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	// Parse config parameters according to indicators.json
	emaPeriod := int(config[0].(float64))      // default: 50
	bandType := config[1].(string)             // default: "MEDIAN"
	maType := config[2].(string)               // default: "EXPONENTIAL"
	multiplier := float64(config[3].(float64)) // default: 5
	offset := 0

	// Handle optional parameter if present
	if len(config) > 4 {
		offset = int(config[4].(float64))
	}

	indicator.emaPeriod = emaPeriod
	indicator.atrPeriod = emaPeriod // Using same period for ATR calculation
	indicator.multiplier = multiplier
	indicator.offset = offset

	// Get stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	// Check if we have enough data
	if len(stockData) < emaPeriod+1 {
		return []IndicatorValues{}
	}

	// Extract price arrays
	highs := make([]float64, len(stockData))
	lows := make([]float64, len(stockData))
	closes := make([]float64, len(stockData))
	medianPrices := make([]float64, len(stockData))

	for i, data := range stockData {
		highs[i] = data.High
		lows[i] = data.Low
		closes[i] = data.Close
		medianPrices[i] = (data.High + data.Low) / 2.0
	}

	// Step 1: Calculate True Range (TR) for each period
	trueRange := make([]float64, len(stockData))
	trueRange[0] = highs[0] - lows[0] // First TR is just High - Low

	for i := 1; i < len(stockData); i++ {
		tr1 := highs[i] - lows[i]
		tr2 := math.Abs(highs[i] - closes[i-1])
		tr3 := math.Abs(lows[i] - closes[i-1])
		trueRange[i] = math.Max(tr1, math.Max(tr2, tr3))
	}
	// Step 2: Calculate Middle Line using utility function for moving average
	middleLineValues := util.CalculateMovingAverage(medianPrices, maType, emaPeriod)

	// Create full-length array with proper positioning
	middleLine := make([]float64, len(stockData))
	startIdx := emaPeriod - 1
	for i, val := range middleLineValues {
		if startIdx+i < len(middleLine) {
			middleLine[startIdx+i] = val
		}
	}

	// Step 3: Calculate ATR using utility function for moving average
	atrValues := util.CalculateMovingAverage(trueRange, maType, emaPeriod)

	// Create full-length array with proper positioning
	atr := make([]float64, len(stockData))
	for i, val := range atrValues {
		if startIdx+i < len(atr) {
			atr[startIdx+i] = val
		}
	}

	// Step 4: Calculate Keltner Channel bands
	upperBand := make([]float64, len(stockData))
	lowerBand := make([]float64, len(stockData))

	for i := emaPeriod - 1; i < len(stockData); i++ {
		upperBand[i] = middleLine[i] + multiplier*atr[i]
		lowerBand[i] = middleLine[i] - multiplier*atr[i]
	}

	// Step 5: Build results starting from valid calculation points
	var result []IndicatorValues

	for i := emaPeriod - 1; i < len(stockData); i++ {
		// Apply offset
		if i < emaPeriod-1+offset {
			continue
		}

		// Check start date filter
		if stockData[i].Timestamp.Before(startDate) {
			continue
		}

		// Determine result value based on band type
		var resultValue float64
		switch bandType {
		case "TOP":
			resultValue = upperBand[i]
		case "BOTTOM":
			resultValue = lowerBand[i]
		default: // "MEDIAN" or any other value
			resultValue = middleLine[i]
		}

		result = append(result, IndicatorValues{
			Timestamp: stockData[i].Timestamp,
			Value:     util.Round(resultValue),
		})
	}

	return result
}

func (indicator *KeltnerChannels) GetDataRange(config map[string]any) int32 {
	// Based on indicators.json field names
	if p, ok := config["period"].(float64); ok {
		indicator.emaPeriod = int(p)
	} else {
		indicator.emaPeriod = 50 // default period value
	}

	// Using the same period for ATR as for the EMA
	indicator.atrPeriod = indicator.emaPeriod

	if m, ok := config["shift"].(float64); ok {
		indicator.multiplier = m
	} else {
		indicator.multiplier = 5.0 // default shift value
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// Return the maximum range required
	maxPeriod := max(indicator.emaPeriod, indicator.atrPeriod)
	// We need maxPeriod+1 data points (due to the use of previous values) + offset
	return int32(maxPeriod + 1 + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for Keltner Channels calculation
// config[0] = emaPeriod (float64), config[1] = bandType (string), config[2] = maType (string)
// config[3] = multiplier (float64), config[4] = offset (float64) [optional]
func (indicator *KeltnerChannels) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	emaPeriod := int(config[0].(float64)) // default: 50
	// config[1] is bandType (string) - not needed for lookback calculation
	// config[2] is maType (string) - not needed for lookback calculation
	// config[3] is multiplier (float64) - not needed for lookback calculation

	offset := 0
	if len(config) > 4 {
		offset = int(config[4].(float64)) // default: 0
	}

	// For Keltner Channels calculation, we need:
	// - emaPeriod for moving average calculation
	// - Same period for ATR calculation
	// - 1 additional for true range calculation
	// - Additional buffer for stability (using 2*maxPeriod like SliceData)
	// - Plus any offset
	maxPeriod := emaPeriod // ATR uses same period
	return int64(2*maxPeriod + 1 + offset)
}
