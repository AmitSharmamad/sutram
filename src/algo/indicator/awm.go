package indicator

import (
	"algo/src/util"
	"math"
)

type AWM struct {
	period      int
	sensitivity float64
}

func (indicator *AWM) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 14 // default period
		}
	}

	sensitivity := indicator.sensitivity
	if sensitivity == 0 {
		if s, ok := config["sensitivity"].(float64); ok {
			sensitivity = s
		} else {
			sensitivity = 2.0 // default sensitivity
		}
	}

	if len(prices) < period*2 {
		return 0.0
	}

	// Calculate momentum components
	var upWave, downWave float64
	for i := len(prices) - period; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			upWave += math.Pow(change, sensitivity)
		} else {
			downWave += math.Pow(-change, sensitivity)
		}
	}

	// Normalize waves
	if upWave+downWave == 0 {
		return 0.0
	}
	waveRatio := (upWave - downWave) / (upWave + downWave)

	// Calculate adaptive factor based on volatility
	volatility := 0.0
	for i := len(prices) - period; i < len(prices)-1; i++ {
		volatility += math.Abs(prices[i+1] - prices[i])
	}
	avgVolatility := volatility / float64(period)
	adaptiveFactor := math.Tanh(avgVolatility * sensitivity)

	// Calculate wave momentum
	waveMomentum := waveRatio * adaptiveFactor * 100

	// Smooth the result using EMA
	ema := &EMA{period: period / 2}
	smoothedAWM := ema.Calculate(Prices{waveMomentum}, nil, point)

	return util.Round(smoothedAWM)
}

func (awm *AWM) GetDataRange() int32 {
	return int32(awm.period * 2)
}

type AnchoredWMA struct {
	period int
	anchor float64
}

func (indicator *AnchoredWMA) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 14 // default period
		}
	}

	anchor := indicator.anchor
	if anchor == 0 {
		if a, ok := config["anchor"].(float64); ok {
			anchor = a
		} else {
			anchor = 2.0 // default anchor factor
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate trend strength for dynamic anchoring
	var trendStrength float64
	for i := len(prices) - period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		trendStrength += math.Abs(change)
	}
	trendStrength /= float64(period)

	// Calculate dynamic weights with anchor adjustment
	weights := make([]float64, period)
	var weightSum float64
	for i := 0; i < period; i++ {
		// Adjust weight based on trend strength and anchor
		weight := float64(i+1) * (1 + (trendStrength * anchor))
		weights[i] = weight
		weightSum += weight
	}

	// Normalize weights
	for i := range weights {
		weights[i] /= weightSum
	}

	// Calculate anchored WMA
	var awma float64
	for i := 0; i < period; i++ {
		price := prices[len(prices)-period+i]
		awma += price * weights[i]
	}

	return util.Round(awma)
}

func (awma *AnchoredWMA) GetDataRange() int32 {
	return int32(awma.period)
}
