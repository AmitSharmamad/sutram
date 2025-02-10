package shared

import "algo/src/util"

type RSI struct {
	period int
}

func (indicator *RSI) Calculate(
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

	// Calculate price changes
	changes := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		changes[i-1] = prices[i] - prices[i-1]
	}

	// Calculate initial averages
	var gains, losses float64
	for i := 0; i < period; i++ {
		if changes[i] > 0 {
			gains += changes[i]
		} else {
			losses -= changes[i]
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// Calculate subsequent values using Wilder's smoothing
	for i := period; i < len(changes); i++ {
		if changes[i] > 0 {
			avgGain = (avgGain*float64(period-1) + changes[i]) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + 0) / float64(period)
		} else {
			avgGain = (avgGain*float64(period-1) + 0) / float64(period)
			avgLoss = (avgLoss*float64(period-1) - changes[i]) / float64(period)
		}
	}

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	rsi := 100.0 - (100.0 / (1.0 + rs))

	return util.Round(rsi)
}

func (rsi *RSI) GetDataRange() int32 {
	return int32(rsi.period)
}
