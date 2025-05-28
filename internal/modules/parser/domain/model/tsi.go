package parser

import (
	util "algo/internal/util/math"
	"math"
)

type TSI struct {
	longPeriod   int
	shortPeriod  int
	signalPeriod int
	offset       int
	lineType     string
}

func (indicator *TSI) SliceData(prices Prices) Prices {
	// We need enough data to calculate TSI properly
	// For TSI calculation we need at least 2*(longPeriod + shortPeriod) data points
	totalPeriod := indicator.longPeriod + indicator.shortPeriod + indicator.signalPeriod
	start := max(0, len(prices)-(indicator.offset+2*totalPeriod))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

func (indicator *TSI) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	field := config[0].(string)
	lineType := config[1].(string)
	longPeriod := int(config[2].(float64))
	shortPeriod := int(config[3].(float64))
	signalPeriod := int(config[4].(float64))
	offset := int(config[5].(float64))

	indicator.longPeriod = longPeriod
	indicator.shortPeriod = shortPeriod
	indicator.signalPeriod = signalPeriod
	indicator.offset = offset
	indicator.lineType = lineType

	prices := GetPrices(*ctx.Data, field)

	// Check if we have enough data points
	totalPeriod := longPeriod + shortPeriod + signalPeriod
	if len(prices) < totalPeriod+offset+1 {
		return 0.0
	}

	prices = indicator.SliceData(prices)

	// Need at least totalPeriod+1 data points to calculate
	if len(prices) < totalPeriod+1 {
		return 0.0
	}

	// Calculate price changes (momentum)
	var momentum = make([]float64, len(prices))
	for i := 1; i < len(prices); i++ {
		momentum[i] = prices[i] - prices[i-1]
	}

	// Calculate the first EMA of the momentum (long period)
	var ema1 = make([]float64, len(prices))
	multiplier1 := 2.0 / float64(indicator.longPeriod+1)

	// Calculate initial value as simple average
	var sumLong float64
	for i := 1; i <= indicator.longPeriod; i++ {
		sumLong += momentum[i]
	}
	ema1[indicator.longPeriod] = sumLong / float64(indicator.longPeriod)

	// Calculate subsequent values
	for i := indicator.longPeriod + 1; i < len(prices); i++ {
		ema1[i] = (momentum[i]-ema1[i-1])*multiplier1 + ema1[i-1]
	}

	// Calculate the second EMA of the first EMA (short period)
	var ema2 = make([]float64, len(prices))
	multiplier2 := 2.0 / float64(indicator.shortPeriod+1)

	// Calculate initial value as simple average
	var sumShort float64
	startIdx := indicator.longPeriod + indicator.shortPeriod
	for i := indicator.longPeriod + 1; i <= startIdx; i++ {
		sumShort += ema1[i]
	}
	ema2[startIdx] = sumShort / float64(indicator.shortPeriod)

	// Calculate subsequent values
	for i := startIdx + 1; i < len(prices); i++ {
		ema2[i] = (ema1[i]-ema2[i-1])*multiplier2 + ema2[i-1]
	}

	// Calculate the first EMA of the absolute momentum (long period)
	var absEma1 = make([]float64, len(prices))

	// Calculate initial value as simple average
	var absSum float64
	for i := 1; i <= indicator.longPeriod; i++ {
		absSum += math.Abs(momentum[i])
	}
	absEma1[indicator.longPeriod] = absSum / float64(indicator.longPeriod)

	// Calculate subsequent values
	for i := indicator.longPeriod + 1; i < len(prices); i++ {
		absEma1[i] = (math.Abs(momentum[i])-absEma1[i-1])*multiplier1 + absEma1[i-1]
	}

	// Calculate the second EMA of the first absolute EMA (short period)
	var absEma2 = make([]float64, len(prices))

	// Calculate initial value as simple average
	var absSumShort float64
	for i := indicator.longPeriod + 1; i <= startIdx; i++ {
		absSumShort += absEma1[i]
	}
	absEma2[startIdx] = absSumShort / float64(indicator.shortPeriod)

	// Calculate subsequent values
	for i := startIdx + 1; i < len(prices); i++ {
		absEma2[i] = (absEma1[i]-absEma2[i-1])*multiplier2 + absEma2[i-1]
	}

	// Calculate TSI
	var tsi = make([]float64, len(prices))
	for i := startIdx; i < len(prices); i++ {
		// Avoid division by zero
		if absEma2[i] != 0 {
			tsi[i] = 100 * (ema2[i] / absEma2[i])
		} else {
			tsi[i] = 0
		}
	}

	// Calculate signal line (EMA of TSI)
	var signal = make([]float64, len(prices))
	multiplierSignal := 2.0 / float64(indicator.signalPeriod+1)

	// Calculate initial value as simple average
	var sumSignal float64
	signalStartIdx := startIdx + indicator.signalPeriod
	for i := startIdx; i < signalStartIdx; i++ {
		sumSignal += tsi[i]
	}
	signal[signalStartIdx-1] = sumSignal / float64(indicator.signalPeriod)

	// Calculate subsequent values
	for i := signalStartIdx; i < len(prices); i++ {
		signal[i] = (tsi[i]-signal[i-1])*multiplierSignal + signal[i-1]
	}

	// Return the appropriate value based on lineType
	if indicator.lineType == "SIGNAL" {
		return util.Round(signal[len(signal)-1])
	} else {
		// Default is "TSI"
		return util.Round(tsi[len(tsi)-1])
	}
}

func (indicator *TSI) GetDataRange(config map[string]any) int32 {
	if p, ok := config["longPeriod"].(float64); ok {
		indicator.longPeriod = int(p)
	} else {
		indicator.longPeriod = 25 // default long period
	}

	if p, ok := config["shortPeriod"].(float64); ok {
		indicator.shortPeriod = int(p)
	} else {
		indicator.shortPeriod = 13 // default short period
	}

	if p, ok := config["signalPeriod"].(float64); ok {
		indicator.signalPeriod = int(p)
	} else {
		indicator.signalPeriod = 13 // default signal period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// We need (longPeriod + shortPeriod + signalPeriod) + 1 data points + offset
	totalPeriod := indicator.longPeriod + indicator.shortPeriod + indicator.signalPeriod
	return int32(totalPeriod + 1 + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for TSI calculation
// config[0] = field (string), config[1] = lineType (string), config[2] = longPeriod (float64)
// config[3] = shortPeriod (float64), config[4] = signalPeriod (float64), config[5] = offset (float64)
func (indicator *TSI) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is field (string) - not needed for lookback calculation
	// config[1] is lineType (string) - not needed for lookback calculation
	longPeriod := int(config[2].(float64))   // default: 25
	shortPeriod := int(config[3].(float64))  // default: 13
	signalPeriod := int(config[4].(float64)) // default: 13
	offset := int(config[5].(float64))       // default: 0

	// For TSI calculation, we need:
	// - Double smoothing requires longPeriod + shortPeriod + signalPeriod
	// - Additional buffer for EMA stabilization (using 2*totalPeriod like SliceData)
	// - Plus 1 for price change calculation
	// - Plus any offset
	totalPeriod := longPeriod + shortPeriod + signalPeriod
	return int64(2*totalPeriod + 1 + offset)
}
