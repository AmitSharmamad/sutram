package parser

import (
	util "algo/internal/util/math"
)

type VolumeOscillator struct {
	shortPeriod int
	longPeriod  int
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

func (indicator *VolumeOscillator) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	shortPeriod := int(config[0].(float64))
	longPeriod := int(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.shortPeriod = shortPeriod
	indicator.longPeriod = longPeriod
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
		// Avoid division by zero
		if longSMA[i] != 0 {
			vo[i] = 100 * ((shortSMA[i] - longSMA[i]) / longSMA[i])
		} else {
			vo[i] = 0
		}
	}

	// Return the last calculated Volume Oscillator value
	return util.Round(vo[len(vo)-1])
}

func (indicator *VolumeOscillator) GetDataRange(config map[string]any) int32 {
	if p, ok := config["shortPeriod"].(float64); ok {
		indicator.shortPeriod = int(p)
	} else {
		indicator.shortPeriod = 5 // default short period
	}

	if p, ok := config["longPeriod"].(float64); ok {
		indicator.longPeriod = int(p)
	} else {
		indicator.longPeriod = 10 // default long period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.longPeriod + indicator.offset)
}
