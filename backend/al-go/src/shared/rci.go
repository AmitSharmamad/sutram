package shared

import (
	"algo/src/util"
	"sort"
)

type RCI struct {
	period int
}

func (indicator *RCI) Calculate(
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

	// Get the subset of prices we're interested in
	subset := prices[len(prices)-period:]

	// Create price-rank pairs
	type PriceRank struct {
		price float64
		index int
	}
	pairs := make([]PriceRank, period)
	for i := 0; i < period; i++ {
		pairs[i] = PriceRank{subset[i], i}
	}

	// Sort by price to determine price ranks
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].price < pairs[j].price
	})

	// Assign price ranks
	priceRanks := make([]int, period)
	for rank, pair := range pairs {
		priceRanks[pair.index] = rank + 1
	}

	// Calculate time ranks (most recent = highest rank)
	timeRanks := make([]int, period)
	for i := 0; i < period; i++ {
		timeRanks[i] = period - i
	}

	// Calculate sum of squared rank differences
	var sumD2 float64
	for i := 0; i < period; i++ {
		d := float64(priceRanks[i] - timeRanks[i])
		sumD2 += d * d
	}

	// RCI formula: RCI = (1 - (6 * sum(d^2)) / (n * (n^2 - 1))) * 100
	n := float64(period)
	rci := (1.0 - (6.0*sumD2)/(n*(n*n-1.0))) * 100.0

	return util.Round(rci)
}

func (rci *RCI) GetDataRange() int32 {
	return int32(rci.period)
}
