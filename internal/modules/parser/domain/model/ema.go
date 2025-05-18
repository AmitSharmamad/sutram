package parser

import (
	util "algo/internal/util/math"
)

type EMA struct {
	period int
	offset int
}

func (indicator *EMA) SliceData(prices Prices) Prices {
	// We need enough data to calculate EMA properly
	// For EMA calculation we need at least 2*period data points
	// This ensures we have enough data for initial SMA and subsequent EMA calculation
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

func (indicator *EMA) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	point := config[0].(string)
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
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

	// Calculate multiplier (alpha)
	multiplier := 2.0 / float64(period+1)

	// Calculate initial SMA for the first EMA value
	var sma float64
	for i := 0; i < period; i++ {
		sma += prices[i]
	}
	sma /= float64(period)

	// Calculate EMA recursively
	// Start with the initial SMA value
	ema := sma

	// Only process remaining data points after the initial period
	for i := period; i < len(prices); i++ {
		// EMA formula: EMA(today) = (Price(today) * multiplier) + (EMA(yesterday) * (1 - multiplier))
		// Alternative form: EMA = Price * multiplier + EMA * (1 - multiplier)
		// Simplified: EMA = EMA + multiplier * (Price - EMA)
		ema = ema + multiplier*(prices[i]-ema)
	}

	return util.Round(ema)
}

func (indicator *EMA) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 14 // default period
	}

	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.period + indicator.offset)
}
