package indicator

import (
	"algo/src/util"
	"math"
)

type DMI struct {
	basePeriod int
	volRange   int
}

func (indicator *DMI) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	basePeriod := indicator.basePeriod
	if basePeriod == 0 {
		if p, ok := config["basePeriod"].(int); ok {
			basePeriod = p
		} else {
			basePeriod = 14 // default base period
		}
	}

	volRange := indicator.volRange
	if volRange == 0 {
		if v, ok := config["volRange"].(int); ok {
			volRange = v
		} else {
			volRange = 5 // default volatility range
		}
	}

	if len(prices) < basePeriod+volRange {
		return 0.0
	}

	// Calculate volatility using standard deviation
	var sum, sumSq float64
	for i := len(prices) - volRange; i < len(prices); i++ {
		sum += prices[i]
		sumSq += prices[i] * prices[i]
	}
	mean := sum / float64(volRange)
	variance := (sumSq - sum*sum/float64(volRange)) / float64(volRange)
	stdDev := math.Sqrt(variance)

	// Adjust period based on volatility
	// Lower period for higher volatility, higher period for lower volatility
	adjustedPeriod := int(float64(basePeriod) * (1 + stdDev/mean))
	if adjustedPeriod < 5 {
		adjustedPeriod = 5 // minimum period
	}
	if adjustedPeriod > basePeriod*2 {
		adjustedPeriod = basePeriod * 2 // maximum period
	}

	// Calculate gains and losses over adjusted period
	var gainSum, lossSum float64
	for i := len(prices) - adjustedPeriod; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gainSum += change
		} else {
			lossSum -= change
		}
	}

	// Calculate DMI similar to RSI but with adjusted period
	avgGain := gainSum / float64(adjustedPeriod)
	avgLoss := lossSum / float64(adjustedPeriod)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	dmi := 100.0 - (100.0 / (1.0 + rs))

	return util.Round(dmi)
}

func (dmi *DMI) GetDataRange() int32 {
	return int32(dmi.basePeriod + dmi.volRange)
}
