package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"
)

type VolumeOscillator struct {
	shortPeriod int
	longPeriod  int
	pop         string // "PCT" or "POINT"
	offset      int
}

func (indicator *VolumeOscillator) SliceData(volumes Prices) Prices {
	// We need enough data to calculate Volume Oscillator properly
	start := max(0, len(volumes)-(indicator.offset+2*indicator.longPeriod))
	end := len(volumes) - indicator.offset

	if start < 0 || end <= start || end > len(volumes) {
		return Prices{}
	}

	return volumes[start:end]
}

// generateIndicatorKey generates a unique key for VolumeOscillator indicator based on its configuration
func (indicator *VolumeOscillator) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 4 {
		return "", fmt.Errorf("VolumeOscillator requires 4 parameters: shortPeriod, longPeriod, pop, offset")
	}

	// Extract shortPeriod parameter
	shortPeriodFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("VolumeOscillator shortPeriod parameter must be a number")
	}
	shortPeriod := int(shortPeriodFloat)

	// Extract longPeriod parameter
	longPeriodFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("VolumeOscillator longPeriod parameter must be a number")
	}
	longPeriod := int(longPeriodFloat)

	// Extract pop parameter
	pop, ok := config[2].(string)
	if !ok {
		return "", fmt.Errorf("VolumeOscillator pop parameter must be a string")
	}

	// Extract offset parameter
	offsetFloat, ok := config[3].(float64)
	if !ok {
		return "", fmt.Errorf("VolumeOscillator offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include non-default values
	baseKey := fmt.Sprintf("VO_%d_%d_%s", shortPeriod, longPeriod, strings.ToUpper(pop))

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

func (indicator *VolumeOscillator) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	shortPeriod := int(config[0].(float64))
	longPeriod := int(config[1].(float64))
	pop := config[2].(string)
	offset := int(config[3].(float64))

	indicator.shortPeriod = shortPeriod
	indicator.longPeriod = longPeriod
	indicator.pop = pop
	indicator.offset = offset

	volumes := GetPrices(*ctx.Data, "volume")

	// Check if we have enough data points
	if len(volumes) < longPeriod+offset {
		return 0.0
	}

	volumes = indicator.SliceData(volumes)

	// Need at least longPeriod data points to calculate
	if len(volumes) < longPeriod {
		return 0.0
	}

	// Calculate short-term SMA of volume
	var shortSMA = make([]float64, len(volumes))

	// For each point, calculate SMA for the short period
	for i := indicator.shortPeriod - 1; i < len(volumes); i++ {
		var sum float64
		for j := i - indicator.shortPeriod + 1; j <= i; j++ {
			sum += volumes[j]
		}
		shortSMA[i] = sum / float64(indicator.shortPeriod)
	}

	// Calculate long-term SMA of volume
	var longSMA = make([]float64, len(volumes))

	// For each point, calculate SMA for the long period
	for i := indicator.longPeriod - 1; i < len(volumes); i++ {
		var sum float64
		for j := i - indicator.longPeriod + 1; j <= i; j++ {
			sum += volumes[j]
		}
		longSMA[i] = sum / float64(indicator.longPeriod)
	}

	// Calculate Volume Oscillator
	var vo = make([]float64, len(volumes))
	for i := indicator.longPeriod - 1; i < len(volumes); i++ {
		if indicator.pop == "PCT" {
			// Percentage-based Volume Oscillator
			// Avoid division by zero
			if longSMA[i] != 0 {
				vo[i] = 100 * ((shortSMA[i] - longSMA[i]) / longSMA[i])
			} else {
				vo[i] = 0
			}
		} else {
			// POINT-based Volume Oscillator is simply the difference
			vo[i] = shortSMA[i] - longSMA[i]
		}
	}

	// Return the last calculated Volume Oscillator value
	return util.Round(vo[len(vo)-1])
}

func (indicator *VolumeOscillator) GetDataRange(config map[string]any) int32 {
	if p, ok := config["shortPeriod"].(float64); ok {
		indicator.shortPeriod = int(p)
	} else {
		indicator.shortPeriod = 12 // default short period per indicators.json
	}

	if p, ok := config["longPeriod"].(float64); ok {
		indicator.longPeriod = int(p)
	} else {
		indicator.longPeriod = 26 // default long period per indicators.json
	}

	if p, ok := config["pop"].(string); ok {
		indicator.pop = p
	} else {
		indicator.pop = "PCT" // default pop type per indicators.json
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.longPeriod + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for Volume Oscillator calculation
// config[0] = shortPeriod (float64), config[1] = longPeriod (float64), config[2] = pop (string), config[3] = offset (float64)
func (indicator *VolumeOscillator) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is shortPeriod (float64) - not needed since longPeriod is larger
	longPeriod := int(config[1].(float64)) // default: 26
	// config[2] is pop (string) - not needed for lookback calculation
	offset := int(config[3].(float64)) // default: 0

	// For Volume Oscillator calculation, we need:
	// - longPeriod data points for long-term SMA calculation (larger of the two periods)
	// - Additional buffer for stability (using 2*longPeriod like SliceData)
	// - Plus any offset
	return int64(2*longPeriod + offset)
}

// CalculateAll calculates Volume Oscillator values for all data points
func (indicator *VolumeOscillator) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	// Parse config parameters
	shortPeriod := int(config[0].(float64))
	longPeriod := int(config[1].(float64))
	pop := config[2].(string)
	offset := int(config[3].(float64))

	indicator.shortPeriod = shortPeriod
	indicator.longPeriod = longPeriod
	indicator.pop = pop
	indicator.offset = offset

	// Get stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	// Check if we have enough data
	if len(stockData) < longPeriod {
		return []IndicatorValues{}
	}

	// Extract volume array
	volumes := make([]float64, len(stockData))
	for i, data := range stockData {
		volumes[i] = data.Volume
	}

	// Step 1: Calculate short-term SMA of volume using utility function
	shortSMAValues := util.CalculateMovingAverage(volumes, "SIMPLE", shortPeriod)

	// Create full-length array with proper positioning
	shortSMA := make([]float64, len(stockData))
	startIdxShort := shortPeriod - 1
	for i, val := range shortSMAValues {
		if startIdxShort+i < len(shortSMA) {
			shortSMA[startIdxShort+i] = val
		}
	}

	// Step 2: Calculate long-term SMA of volume using utility function
	longSMAValues := util.CalculateMovingAverage(volumes, "SIMPLE", longPeriod)

	// Create full-length array with proper positioning
	longSMA := make([]float64, len(stockData))
	startIdxLong := longPeriod - 1
	for i, val := range longSMAValues {
		if startIdxLong+i < len(longSMA) {
			longSMA[startIdxLong+i] = val
		}
	}

	// Step 3: Calculate Volume Oscillator
	vo := make([]float64, len(stockData))

	for i := longPeriod - 1; i < len(stockData); i++ {
		if pop == "PCT" {
			// Percentage-based Volume Oscillator
			// Avoid division by zero
			if longSMA[i] != 0 {
				vo[i] = 100 * ((shortSMA[i] - longSMA[i]) / longSMA[i])
			} else {
				vo[i] = 0
			}
		} else {
			// POINT-based Volume Oscillator is simply the difference
			vo[i] = shortSMA[i] - longSMA[i]
		}
	}

	// Step 4: Build results starting from valid calculation points
	var result []IndicatorValues

	for i := longPeriod - 1; i < len(stockData); i++ {
		// Apply offset
		if i < longPeriod-1+offset {
			continue
		}

		// Check start date filter
		if stockData[i].Timestamp.Before(startDate) {
			continue
		}

		result = append(result, IndicatorValues{
			Timestamp: stockData[i].Timestamp,
			Value:     util.Round(vo[i]),
		})
	}

	return result
}
