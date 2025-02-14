package shared

import "algo/src/util"

type ROC struct {
	period int
}

func (indicator *ROC) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 12 // default period
		}
	}

	if len(prices) < period+1 {
		return 0.0
	}

	// ROC = ((Current Price - Price n periods ago) / Price n periods ago) * 100
	currentPrice := prices[len(prices)-1]
	oldPrice := prices[len(prices)-1-period]

	if oldPrice == 0 {
		return 0.0
	}

	roc := ((currentPrice - oldPrice) / oldPrice) * 100

	return util.Round(roc)
}

func (roc *ROC) GetDataRange() int32 {
	return int32(roc.period)
}
