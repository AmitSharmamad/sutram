package parser

import (
	util "algo/internal/util/math"
	"math"
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

	// Get all stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	// Check if we have enough data points
	maxPeriod := emaPeriod
	if len(stockData) < maxPeriod+offset+1 {
		return []IndicatorValues{}
	}

	var result []IndicatorValues

	// Find the start index for startDate
	startIndex := -1
	for i, data := range stockData {
		if !data.Timestamp.Before(startDate) {
			startIndex = i
			break
		}
	}

	if startIndex == -1 {
		return []IndicatorValues{}
	}

	// Calculate for each data point starting from startDate
	for i := startIndex; i < len(stockData); i++ {
		if i < maxPeriod {
			continue // Need at least 'maxPeriod+1' data points (for TR calculation)
		}

		// Get data for calculation window ending at current point
		var highs, lows, closes []float64
		for j := i - maxPeriod; j <= i; j++ {
			highs = append(highs, stockData[j].High)
			lows = append(lows, stockData[j].Low)
			closes = append(closes, stockData[j].Close)
		}

		// Calculate price series (using median price for Keltner)
		var prices = make([]float64, len(closes))
		for k := 0; k < len(closes); k++ {
			prices[k] = (highs[k] + lows[k]) / 2.0
		}

		// Calculate moving average using utility function
		emas := util.CalculateMovingAverage(prices, maType, indicator.emaPeriod)

		// Calculate ATR (Average True Range)
		var tr = make([]float64, len(closes))

		// Calculate True Range
		for k := 1; k < len(closes); k++ {
			tr[k] = max(
				highs[k]-lows[k],
				max(
					math.Abs(highs[k]-closes[k-1]),
					math.Abs(lows[k]-closes[k-1]),
				),
			)
		}

		// Calculate ATR using moving average on TR
		trData := tr[1:] // Skip the first zero value
		atrs := util.CalculateMovingAverage(trData, maType, indicator.atrPeriod)

		// Get the last values for MA and ATR
		var middleLine float64
		if len(emas) > 0 {
			middleLine = emas[len(emas)-1]
		}

		var atr float64
		if len(atrs) > 0 {
			atr = atrs[len(atrs)-1]
		}

		// Return the middle line value (like the original Calculate method for MEDIAN)
		var resultValue float64
		switch bandType {
		case "TOP":
			resultValue = middleLine + indicator.multiplier*atr
		case "BOTTOM":
			resultValue = middleLine - indicator.multiplier*atr
		default: // "MEDIAN" or any other value
			resultValue = middleLine
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
