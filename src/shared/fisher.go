package shared

import (
	"algo/src/util"
	"math"
)

type Fisher struct {
	period int
}

func (indicator *Fisher) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 9 // default period
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Find highest high and lowest low in the period
	highest := prices[len(prices)-period]
	lowest := prices[len(prices)-period]
	for i := len(prices) - period; i < len(prices); i++ {
		if prices[i] > highest {
			highest = prices[i]
		}
		if prices[i] < lowest {
			lowest = prices[i]
		}
	}

	// Calculate value position
	value := 0.0
	if highest != lowest {
		lastPrice := prices[len(prices)-1]
		value = 0.66 * ((lastPrice-lowest)/(highest-lowest) - 0.5 + value)
	}

	// Ensure value stays within [-1.0, 1.0]
	if value > 0.999 {
		value = 0.999
	}
	if value < -0.999 {
		value = -0.999
	}

	// Calculate Fisher Transform
	fisher := 0.5 * math.Log((1.0+value)/(1.0-value))

	return util.Round(fisher)
}

func (f *Fisher) GetDataRange() int32 {
	return int32(f.period)
}
