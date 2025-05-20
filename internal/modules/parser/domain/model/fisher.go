package parser

import (
	util "algo/internal/util/math"
	"math"
)

type FisherTransform struct {
	period int
	offset int
}

func (indicator *FisherTransform) SliceData(prices Prices) Prices {
	// We need enough data to calculate Fisher Transform properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// Helper method to slice all required data
func (indicator *FisherTransform) sliceAllData(highs, lows Prices) (Prices, Prices) {
	// We need enough data to calculate Fisher Transform properly
	start := max(0, len(highs)-(indicator.offset+2*indicator.period))
	end := len(highs) - indicator.offset

	if start < 0 || end <= start || end > len(highs) {
		return Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end]
}

func (indicator *FisherTransform) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	period := int(config[0].(float64))
	offset := int(config[1].(float64))

	indicator.period = period
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")

	// Check if we have enough data points
	if len(highs) < period+offset || len(lows) < period+offset {
		return 0.0
	}

	highs, lows = indicator.sliceAllData(highs, lows)

	// Need at least period data points to calculate
	if len(highs) < period || len(lows) < period {
		return 0.0
	}

	// Calculate median prices
	var medianPrices = make([]float64, len(highs))
	for i := 0; i < len(highs); i++ {
		medianPrices[i] = (highs[i] + lows[i]) / 2.0
	}

	// Find highest high and lowest low in the period for normalization
	var fisher = make([]float64, len(highs))
	var value = make([]float64, len(highs))

	// Initialize
	value[0] = 0.0
	fisher[0] = 0.0

	// Calculate normalized price and apply Fisher Transform
	for i := 1; i < len(highs); i++ {
		// Find highest high and lowest low in the lookback period
		var highest, lowest float64
		highest = medianPrices[i]
		lowest = medianPrices[i]

		lookbackStart := max(0, i-period+1)
		for j := lookbackStart; j <= i; j++ {
			if medianPrices[j] > highest {
				highest = medianPrices[j]
			}
			if medianPrices[j] < lowest {
				lowest = medianPrices[j]
			}
		}

		// Avoid division by zero
		if highest == lowest {
			value[i] = 0.0
		} else {
			// Normalize price to -1 to +1 range
			normalizedPrice := 2.0*((medianPrices[i]-lowest)/(highest-lowest)) - 1.0

			// Apply smoothing
			value[i] = 0.33*normalizedPrice + 0.67*value[i-1]

			// Bound value between -0.99 and 0.99 to avoid infinity in log transform
			if value[i] > 0.99 {
				value[i] = 0.99
			}
			if value[i] < -0.99 {
				value[i] = -0.99
			}
		}

		// Calculate Fisher Transform
		// Fisher Transform = 0.5 * ln((1 + value) / (1 - value))
		fisher[i] = 0.5 * math.Log((1.0+value[i])/(1.0-value[i]))
	}

	// Return the last calculated Fisher Transform value
	return util.Round(fisher[len(fisher)-1])
}

func (indicator *FisherTransform) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 10 // default period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.period + indicator.offset)
}
