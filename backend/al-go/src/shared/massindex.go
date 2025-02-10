package shared

import "algo/src/util"

type MassIndex struct {
	period    int
	emaPeriod int
}

func (indicator *MassIndex) Calculate(
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

	emaPeriod := indicator.emaPeriod
	if emaPeriod == 0 {
		if p, ok := config["emaPeriod"].(int); ok {
			emaPeriod = p
		} else {
			emaPeriod = 9 // default EMA period
		}
	}

	if len(prices) < period+emaPeriod {
		return 0.0
	}

	// Calculate price ranges
	ranges := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		ranges[i-1] = prices[i] - prices[i-1]
	}

	// Calculate first EMA of ranges
	ema1 := &EMA{period: emaPeriod}
	ema1Value := ema1.Calculate(Prices(ranges), nil, point)

	// Calculate second EMA of the first EMA
	ema2 := &EMA{period: emaPeriod}
	ema2Value := ema2.Calculate(Prices{ema1Value}, nil, point)

	if ema2Value == 0 {
		return 0.0
	}

	// Calculate Mass Index ratio
	ratio := ema1Value / ema2Value

	// Sum the ratios over the period
	massIndex := ratio

	return util.Round(massIndex)
}

func (mi *MassIndex) GetDataRange() int32 {
	return int32(mi.period + mi.emaPeriod)
}
