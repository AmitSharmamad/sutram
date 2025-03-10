package indicator

import (
	"algo/src/shared"
	"algo/src/util"
	"math"
)

type ADX struct {
	period int
}

func (indicator *ADX) Calculate(
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

	if len(prices) < period*2 {
		return 0.0
	}

	// Calculate TR (True Range) and DM (Directional Movement)
	var trSum, plusDMSum, minusDMSum float64
	for i := 1; i < len(prices); i++ {
		high := prices[i]
		low := prices[i-1]
		close := prices[i-1]

		tr1 := math.Abs(high - low)
		tr2 := math.Abs(high - close)
		tr3 := math.Abs(low - close)
		tr := math.Max(tr1, math.Max(tr2, tr3))
		trSum += tr

		// Directional Movement
		upMove := high - prices[i-1]
		downMove := prices[i-1] - low

		plusDM := 0.0
		minusDM := 0.0

		if upMove > downMove && upMove > 0 {
			plusDM = upMove
		}
		if downMove > upMove && downMove > 0 {
			minusDM = downMove
		}

		plusDMSum += plusDM
		minusDMSum += minusDM
	}

	// Calculate smoothed averages
	trN := trSum / float64(period)
	plusDMN := plusDMSum / float64(period)
	minusDMN := minusDMSum / float64(period)

	// Calculate DI+ and DI-
	plusDI := (plusDMN / trN) * 100
	minusDI := (minusDMN / trN) * 100

	// Calculate ADX
	dx := math.Abs(plusDI-minusDI) / (plusDI + minusDI) * 100
	adx := dx // In a complete implementation, this would be smoothed over the period

	return util.Round(adx)
}

func (adx *ADX) GetDataRange() int32 {
	return int32(adx.period * 2)
}

func (indicator *ADX) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
