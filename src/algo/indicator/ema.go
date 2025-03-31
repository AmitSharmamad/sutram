package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type EMA struct {
	period int
	offset int
}

func (indicator *EMA) sliceData(prices Prices) Prices {
	start := len(prices) - indicator.offset - indicator.period
	end := len(prices) - indicator.offset

	return prices[start:end]
}

func (indicator *EMA) Calculate(
	prices Prices,
	config []any,
	point string,
) float64 {
	period := indicator.period
	prices = indicator.sliceData(prices)

	if len(prices) < period {
		return 0.0
	}

	// Calculate multiplier
	multiplier := 2.0 / float64(period+1)

	// Calculate initial SMA
	var sma float64
	for i := 0; i < period; i++ {
		sma += prices[i]
	}
	sma /= float64(period)

	// Calculate EMA
	ema := sma
	for i := period; i < len(prices); i++ {
		ema = (prices[i]-ema)*multiplier + ema
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

func (indicator *EMA) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
// 	prices := GetPrices(data, point)

// 	for i := 0; i < len(data); i++ {
// 		windowPrices := prices[:i+1]
// 		result := indicator.Calculate(windowPrices, config, point)
// 		results = append(results, result)
// 	}

	return results
}
