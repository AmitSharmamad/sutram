package indicator

import (
	"algo/src/util"
	"math"
)

type MesaSineWave struct {
	period int
}

func (indicator *MesaSineWave) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(int); ok {
			period = p
		} else {
			period = 25 // default period
		}
	}

	if len(prices) < period {
		return 0.0
	}

	// Calculate smoothed price
	smooth := make([]float64, period)
	smooth[0] = prices[len(prices)-period]
	for i := 1; i < period; i++ {
		smooth[i] = (4*prices[len(prices)-period+i] + 3*smooth[i-1]) / 7
	}

	// Calculate phase and quadrature
	var i1, q1 float64
	for i := 0; i < period; i++ {
		angle := 2.0 * math.Pi * float64(i) / float64(period)
		i1 += smooth[i] * math.Cos(angle)
		q1 += smooth[i] * math.Sin(angle)
	}
	i1 = i1 * 2.0 / float64(period)
	q1 = q1 * 2.0 / float64(period)

	// Calculate phase
	phase := math.Atan2(q1, i1)

	// Convert phase to sine wave
	sineWave := math.Sin(phase)

	// Scale to percentage
	sineWave = sineWave * 100

	return util.Round(sineWave)
}

func (msw *MesaSineWave) GetDataRange() int32 {
	return int32(msw.period)
}
