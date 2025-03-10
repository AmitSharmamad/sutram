package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type BopIndicator struct {
	scalar float64
}

func (indicator *BopIndicator) Calculate(
	open Prices,
	high Prices,
	low Prices,
	close Prices,
	config map[string]any,
	point string,
) float64 {
	scalar := indicator.scalar
	if scalar == 0 {
		if s, ok := config["scalar"].(float64); ok {
			scalar = s
		} else {
			scalar = 1.0 // default scalar
		}
	}

	if len(open) == 0 || len(high) == 0 || len(low) == 0 || len(close) == 0 {
		return 0.0
	}

	// Calculate BOP
	bop := scalar * (close[len(close)-1] - open[len(open)-1]) / (high[len(high)-1] - low[len(low)-1])

	return util.Round(bop)
}

func (indicator *BopIndicator) GetDataRange() int32 {
	return 1
}

func (indicator *BopIndicator) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
