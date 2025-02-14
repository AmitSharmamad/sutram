package indicator

import (
	"algo/src/util"
	"math"
)

type Entropy struct {
	period int
	bins   int
}

func (indicator *Entropy) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 10 // default period
		}
	}

	bins := indicator.bins
	if bins == 0 {
		if b, ok := config["bins"].(int); ok {
			bins = b
		} else {
			bins = 10 // default number of bins
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Get the relevant price subset
	subset := prices[len(prices)-period:]

	// Find min and max for binning
	min := subset[0]
	max := subset[0]
	for _, price := range subset {
		if price < min {
			min = price
		}
		if price > max {
			max = price
		}
	}

	// Create histogram bins
	binCounts := make([]int, bins)
	binWidth := (max - min) / float64(bins)
	if binWidth == 0 {
		return 0.0
	}

	// Fill bins
	for _, price := range subset {
		binIndex := int((price - min) / binWidth)
		if binIndex == bins {
			binIndex--
		}
		binCounts[binIndex]++
	}

	// Calculate probabilities and entropy
	entropy := 0.0
	for _, count := range binCounts {
		if count > 0 {
			prob := float64(count) / float64(period)
			entropy -= prob * math.Log2(prob)
		}
	}

	// Normalize entropy to [0, 100]
	maxEntropy := math.Log2(float64(bins))
	if maxEntropy > 0 {
		entropy = (entropy / maxEntropy) * 100
	}

	return util.Round(entropy)
}

func (entropy *Entropy) GetDataRange() int32 {
	return int32(entropy.period)
}
