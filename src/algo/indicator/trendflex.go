package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type TrendFlex struct {
	period int
}

func (indicator *TrendFlex) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 20 // default period
		}
	}

	if len(prices) < period*2 {
		return 0.0
	}

	// Calculate linear regression coefficients
	var sumX, sumY, sumXY, sumX2 float64
	n := float64(period)

	for i := 0; i < period; i++ {
		x := float64(i)
		y := prices[len(prices)-period+i]

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	// Calculate standard error
	var stdError float64
	for i := 0; i < period; i++ {
		x := float64(i)
		y := prices[len(prices)-period+i]
		yFit := slope*x + intercept
		stdError += math.Pow(y-yFit, 2)
	}
	stdError = math.Sqrt(stdError / float64(period-2))

	if stdError == 0 {
		return 0.0
	}

	// Calculate TrendFlex
	trendFlex := slope / stdError

	return util.Round(trendFlex * 100) // Scaled for readability
}

func (tf *TrendFlex) GetDataRange() int32 {
	return int32(tf.period * 2)
}

func (indicator *TrendFlex) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
