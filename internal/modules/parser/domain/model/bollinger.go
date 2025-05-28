package parser

import (
	util "algo/internal/util/math"
	"math"
)

type BollingerBands struct {
	period     int
	multiplier float64
	offset     int
}

func (indicator *BollingerBands) SliceData(prices Prices) Prices {
	// We need enough data to calculate Bollinger Bands properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// BOLLINGER_BANDS(SOURCE (OPEN, HIGH, LOW, CLOSE), PERIOD, STANDARD_DEVIATION, OFFSET)
// BOLLINGER_BANDS(CLOSE, 20, 2, 0)
func (indicator *BollingerBands) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	point := config[0].(string)
	period := int(config[1].(float64))
	multiplier := float64(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.period = period
	indicator.multiplier = multiplier
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, point)

	// Check if we have enough data points
	if len(prices) < period+offset {
		return 0.0
	}

	prices = indicator.SliceData(prices)

	// Need at least period data points to calculate
	if len(prices) < period {
		return 0.0
	}

	// Calculate SMA
	var sum float64
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)

	// Calculate standard deviation
	var variance float64
	for i := len(prices) - period; i < len(prices); i++ {
		diff := prices[i] - sma
		variance += diff * diff
	}
	variance /= float64(period)

	// Calculate Bollinger Bands
	middleBand := sma

	// Return the middle band value (SMA)
	return util.Round(middleBand)
}

func (indicator *BollingerBands) CalculateAll(
	config []any,
	ctx *EvaluationContext,
) (float64, float64, float64) {
	point := config[0].(string)
	period := int(config[1].(float64))
	multiplier := float64(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.period = period
	indicator.multiplier = multiplier
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, point)

	// Check if we have enough data points
	if len(prices) < period+offset {
		return 0.0, 0.0, 0.0
	}

	prices = indicator.SliceData(prices)

	// Need at least period data points to calculate
	if len(prices) < period {
		return 0.0, 0.0, 0.0
	}

	// Calculate SMA
	var sum float64
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)

	// Calculate standard deviation
	var variance float64
	for i := len(prices) - period; i < len(prices); i++ {
		diff := prices[i] - sma
		variance += diff * diff
	}
	variance /= float64(period)
	stdDev := math.Sqrt(variance)

	// Calculate Bollinger Bands
	middleBand := sma
	upperBand := sma + indicator.multiplier*stdDev
	lowerBand := sma - indicator.multiplier*stdDev

	return util.Round(upperBand), util.Round(middleBand), util.Round(lowerBand)
}

func (indicator *BollingerBands) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 20 // default period
	}

	if m, ok := config["multiplier"].(float64); ok {
		indicator.multiplier = m
	} else {
		indicator.multiplier = 2.0 // default multiplier
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.period + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for Bollinger Bands calculation
// config[0] = point (string), config[1] = period (float64), config[2] = multiplier (float64), config[3] = offset (float64)
func (indicator *BollingerBands) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is point (string) - not needed for lookback calculation
	period := int(config[1].(float64)) // default: 20
	// config[2] is multiplier (float64) - not needed for lookback calculation
	offset := int(config[3].(float64)) // default: 0

	// For Bollinger Bands calculation, we need:
	// - Period data points for SMA and standard deviation calculation
	// - Additional buffer for stability (using 2*period like SliceData)
	// - Plus any offset
	return int64(2*period + offset)
}
