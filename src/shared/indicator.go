package shared

import (
	"context"
)

type Indicator struct {
	Config map[string]any
}

type IndicatorService interface {
	GetIndicatorData(ctx context.Context, name string, period int, data *StockData) []float64
}
