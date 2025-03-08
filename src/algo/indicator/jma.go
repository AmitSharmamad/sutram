package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type JMA struct {
	period int
	phase  float64
	power  float64
}

func (indicator *JMA) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 7 // default period
		}
	}

	phase := indicator.phase
	if phase == 0 {
		if p, ok := config["phase"].(float64); ok {
			phase = p
		} else {
			phase = 0 // default phase
		}
	}

	power := indicator.power
	if power == 0 {
		if p, ok := config["power"].(float64); ok {
			power = p
		} else {
			power = 2 // default power
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate beta and alpha parameters
	beta := 0.45 * float64(period-1) / (0.45*float64(period-1) + 2)
	alpha := math.Pow(beta, power)

	// Phase-adjusted alpha
	phaseRatio := phase/100.0 + 1.5
	alpha = alpha * phaseRatio

	// Initialize e0 and e1 series
	e0 := make([]float64, len(prices))
	e1 := make([]float64, len(prices))
	e2 := make([]float64, len(prices))

	// Calculate first value
	e0[0] = prices[0]
	e1[0] = prices[0]
	e2[0] = prices[0]

	// Calculate subsequent values
	for i := 1; i < len(prices); i++ {
		// First smoothing
		e0[i] = alpha*prices[i] + (1-alpha)*e0[i-1]

		// Second smoothing
		e1[i] = alpha*e0[i] + (1-alpha)*e1[i-1]

		// Third smoothing
		e2[i] = alpha*e1[i] + (1-alpha)*e2[i-1]
	}

	// Return final JMA value
	return util.Round(e2[len(e2)-1])
}

func (jma *JMA) GetDataRange() int32 {
	return int32(jma.period)
}

func (indicator *JMA) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
