package shared

import (
	"algo/src/util"
	"math"
)

type ALMA struct {
	period int
	sigma  float64
	offset float64
}

func (indicator *ALMA) Calculate(
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

	sigma := indicator.sigma
	if sigma == 0 {
		if s, ok := config["sigma"].(float64); ok {
			sigma = s
		} else {
			sigma = 6.0 // default sigma
		}
	}

	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(float64); ok {
			offset = o
		} else {
			offset = 0.85 // default offset
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate window size and array position
	m := offset * float64(period-1)
	s := float64(period) / sigma

	// Calculate weights
	weights := make([]float64, period)
	var weightSum float64
	for i := 0; i < period; i++ {
		pos := float64(i) - m
		weights[i] = math.Exp(-pos * pos / (2 * s * s))
		weightSum += weights[i]
	}

	// Normalize weights
	if weightSum != 0 {
		for i := range weights {
			weights[i] /= weightSum
		}
	}

	// Calculate ALMA
	var alma float64
	for i := 0; i < period; i++ {
		price := prices[len(prices)-period+i]
		alma += price * weights[i]
	}

	return util.Round(alma)
}

func (alma *ALMA) GetDataRange() int32 {
	return int32(alma.period)
}
