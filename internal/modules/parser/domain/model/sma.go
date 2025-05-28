package parser

import (
	util "algo/internal/util/math"
)

type SMA struct {
	period int
	offset int
}

func (indicator *SMA) SliceData(prices Prices) Prices {
	start := len(prices) - indicator.offset - indicator.period
	end := len(prices) - indicator.offset
	return prices[start:end]
}

// SMA(field, period, offset)
func (indicator *SMA) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	field := config[0].(string)
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, field)

	if len(prices) < period+offset {
		return 0.0
	}

	// Calculate SMA
	var sum float64
	start := len(prices) - period - offset
	end := len(prices) - offset

	for i := start; i < end; i++ {
		sum += prices[i]
	}

	return util.Round(sum / float64(period))
}

func (indicator *SMA) GetDataRange(config map[string]any) int32 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(float64); ok {
			period = int(p)
		} else {
			period = 20 // default period from indicators.json
		}
	}

	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(float64); ok {
			offset = int(o)
		} else {
			offset = 0
		}
	}

	return int32(period + offset)
}

// GetLookBack returns the number of historical data points needed for SMA calculation
// config[0] = field (string), config[1] = period (float64), config[2] = offset (float64)
func (indicator *SMA) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is field (string) - not needed for lookback calculation
	// config[1] is period (float64)
	// config[2] is offset (float64)

	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	// For SMA calculation, we need exactly period + offset data points
	// The SMA needs 'period' consecutive data points plus any offset
	return int64(period + offset)
}
