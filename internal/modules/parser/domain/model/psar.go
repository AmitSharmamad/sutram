package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"time"
)

type ParabolicSAR struct {
	acceleration    float64
	maxAcceleration float64
	offset          int
}

func (indicator *ParabolicSAR) SliceData(prices Prices) Prices {
	// We need enough data to calculate Parabolic SAR properly
	start := max(0, len(prices)-(indicator.offset+20)) // Use more data for trend identification
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for ParabolicSAR indicator based on its configuration
func (indicator *ParabolicSAR) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 3 {
		return "", fmt.Errorf("ParabolicSAR requires 3 parameters: minimumAF, maximumAF, offset")
	}

	// Extract minimumAF parameter
	minimumAF, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("ParabolicSAR minimumAF parameter must be a number")
	}

	// Extract maximumAF parameter
	maximumAF, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("ParabolicSAR maximumAF parameter must be a number")
	}

	// Extract offset parameter
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("ParabolicSAR offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	baseKey := fmt.Sprintf("PSAR_%.3f_%.3f", minimumAF, maximumAF)

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// Helper method to slice all required data
func (indicator *ParabolicSAR) sliceAllData(highs, lows Prices) (Prices, Prices) {
	// We need enough data to calculate PSAR properly
	start := max(0, len(highs)-(indicator.offset+10)) // Need some initial data points
	end := len(highs) - indicator.offset

	if start < 0 || end <= start || end > len(highs) {
		return Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end]
}

func (indicator *ParabolicSAR) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	minimumAF := float64(config[0].(float64))
	maximumAF := float64(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.acceleration = minimumAF
	indicator.maxAcceleration = maximumAF
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")

	// Check if we have enough data points
	if len(highs) < 2+offset || len(lows) < 2+offset {
		return 0.0
	}

	highs, lows = indicator.sliceAllData(highs, lows)

	// Need at least 2 data points to calculate
	if len(highs) < 2 || len(lows) < 2 {
		return 0.0
	}

	// Initialize variables
	var sar = make([]float64, len(highs))
	var ep = make([]float64, len(highs)) // Extreme Point
	var af = make([]float64, len(highs)) // Acceleration Factor
	var isUptrend = make([]bool, len(highs))

	// Determine initial trend
	// Assuming uptrend initially for simplicity
	isUptrend[0] = true
	isUptrend[1] = highs[1] > highs[0]

	// Set initial extreme point (EP)
	if isUptrend[1] {
		ep[1] = highs[1]
		sar[1] = lows[0] // SAR starts at the lowest low of the previous period
	} else {
		ep[1] = lows[1]
		sar[1] = highs[0] // SAR starts at the highest high of the previous period
	}

	af[1] = indicator.acceleration // Start with the initial acceleration factor

	// Calculate SAR for the rest of the data
	for i := 2; i < len(highs); i++ {
		// Previous SAR
		sar[i] = sar[i-1]

		// If previous trend was up
		if isUptrend[i-1] {
			// Check if current low is less than previous SAR (trend reversal)
			if lows[i] < sar[i] {
				// Trend changes to down
				isUptrend[i] = false

				// SAR becomes the previous extreme point
				sar[i] = ep[i-1]

				// Reset acceleration factor and set new extreme point
				af[i] = indicator.acceleration
				ep[i] = lows[i]
			} else {
				// Continue uptrend
				isUptrend[i] = true

				// Update SAR: SAR(today) = SAR(yesterday) + AF * (EP - SAR(yesterday))
				sar[i] = sar[i-1] + af[i-1]*(ep[i-1]-sar[i-1])

				// Ensure SAR is not higher than the lows of the two previous periods
				sar[i] = min(sar[i], min(lows[i-1], lows[i-2]))

				// Update extreme point if new high
				if highs[i] > ep[i-1] {
					ep[i] = highs[i]
					// Increase acceleration factor, up to maximum
					af[i] = min(af[i-1]+indicator.acceleration, indicator.maxAcceleration)
				} else {
					ep[i] = ep[i-1]
					af[i] = af[i-1]
				}
			}
		} else { // Previous trend was down
			// Check if current high is greater than previous SAR (trend reversal)
			if highs[i] > sar[i] {
				// Trend changes to up
				isUptrend[i] = true

				// SAR becomes the previous extreme point
				sar[i] = ep[i-1]

				// Reset acceleration factor and set new extreme point
				af[i] = indicator.acceleration
				ep[i] = highs[i]
			} else {
				// Continue downtrend
				isUptrend[i] = false

				// Update SAR: SAR(today) = SAR(yesterday) + AF * (EP - SAR(yesterday))
				sar[i] = sar[i-1] + af[i-1]*(ep[i-1]-sar[i-1])

				// Ensure SAR is not lower than the highs of the two previous periods
				sar[i] = max(sar[i], max(highs[i-1], highs[i-2]))

				// Update extreme point if new low
				if lows[i] < ep[i-1] {
					ep[i] = lows[i]
					// Increase acceleration factor, up to maximum
					af[i] = min(af[i-1]+indicator.acceleration, indicator.maxAcceleration)
				} else {
					ep[i] = ep[i-1]
					af[i] = af[i-1]
				}
			}
		}
	}

	// Return the last calculated SAR value
	return util.Round(sar[len(sar)-1])
}

func (indicator *ParabolicSAR) GetDataRange(config map[string]any) int32 {
	if a, ok := config["minimumAF"].(float64); ok {
		indicator.acceleration = a
	} else {
		indicator.acceleration = 0.02 // default initial acceleration
	}

	if m, ok := config["maximumAF"].(float64); ok {
		indicator.maxAcceleration = m
	} else {
		indicator.maxAcceleration = 0.2 // default maximum acceleration
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// We need at least 20 data points to establish a good trend + offset
	return int32(20 + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for Parabolic SAR calculation
// config[0] = acceleration (float64), config[1] = maxAcceleration (float64), config[2] = offset (float64)
func (indicator *ParabolicSAR) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is acceleration (float64) - not needed for lookback calculation
	// config[1] is maxAcceleration (float64) - not needed for lookback calculation
	offset := int(config[2].(float64)) // default: 0

	// For Parabolic SAR calculation, we need:
	// - At least 20 data points to establish a good trend pattern
	// - Additional buffer for stability (minimum 50 points recommended for PSAR)
	// - Plus any offset
	// PSAR needs significant historical data to properly establish trend direction
	return int64(50 + offset)
}

// CalculateAll calculates Parabolic SAR values for all data points
func (indicator *ParabolicSAR) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	minimumAF := float64(config[0].(float64))
	maximumAF := float64(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.acceleration = minimumAF
	indicator.maxAcceleration = maximumAF
	indicator.offset = offset

	stockData := *ctx.Data
	var results []IndicatorValues

	// Need enough data for calculation
	minDataPoints := 20 + offset // Using the same as GetDataRange
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

		// Calculate Parabolic SAR value for this point
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
