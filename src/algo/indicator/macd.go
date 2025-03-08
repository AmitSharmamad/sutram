package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type MACD struct {
	fastPeriod   int
	slowPeriod   int
	signalPeriod int
}

func (indicator *MACD) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	// Set default periods if not configured
	fastPeriod := indicator.fastPeriod
	if fastPeriod == 0 {
		if p, ok := config["fastPeriod"].(int); ok {
			fastPeriod = p
		} else {
			fastPeriod = 12 // default fast period
		}
	}

	slowPeriod := indicator.slowPeriod
	if slowPeriod == 0 {
		if p, ok := config["slowPeriod"].(int); ok {
			slowPeriod = p
		} else {
			slowPeriod = 26 // default slow period
		}
	}

	signalPeriod := indicator.signalPeriod
	if signalPeriod == 0 {
		if p, ok := config["signalPeriod"].(int); ok {
			signalPeriod = p
		} else {
			signalPeriod = 9 // default signal period
		}
	}

	if len(prices) < slowPeriod {
		return 0.0
	}

	// Calculate fast EMA
	fastEMA := &EMA{period: fastPeriod}
	fastValue := fastEMA.Calculate(prices, nil, point)

	// Calculate slow EMA
	slowEMA := &EMA{period: slowPeriod}
	slowValue := slowEMA.Calculate(prices, nil, point)

	// MACD Line = Fast EMA - Slow EMA
	macdLine := fastValue - slowValue

	// Calculate Signal Line (EMA of MACD Line)
	// For simplicity, we're returning the MACD line value
	// In a real implementation, you'd calculate the signal line and histogram

	return util.Round(macdLine)
}

func (macd *MACD) GetDataRange() int32 {
	// Return the longest period needed
	if macd.slowPeriod > 0 && macd.signalPeriod > 0 {
		return int32(macd.slowPeriod + macd.signalPeriod)
	}
	return 26 + 9 // default slowPeriod + signalPeriod
}

func (indicator *MACD) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	prices := GetPrices(data, point)

	for i := 0; i < len(data); i++ {
		windowPrices := prices[:i+1]
		result := indicator.Calculate(windowPrices, config, point)
		results = append(results, result)
	}

	return results
}
