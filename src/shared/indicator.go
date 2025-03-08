package shared

import (
	"context"
)

type IndicatorService interface {
	GetIndicatorData(ctx context.Context, name string, data *StockData) []float64
}
