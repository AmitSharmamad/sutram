package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type CMO struct {
	period int
}

func (indicator *CMO) Calculate(
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

	if len(prices) < period+1 {
		return 0.0
	}

	var sumGain, sumLoss float64

	// Calculate gains and losses
	for i := len(prices) - period; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			sumGain += change
		} else {
			sumLoss -= change // Make loss positive
		}
	}

	if sumGain+sumLoss == 0 {
		return 0
	}

	// CMO = 100 * ((sumGain - sumLoss) / (sumGain + sumLoss))
	cmo := 100 * ((sumGain - sumLoss) / (sumGain + sumLoss))

	return util.Round(cmo)
}

func (cmo *CMO) GetDataRange() int32 {
	return int32(cmo.period)
}

func (indicator *CMO) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
