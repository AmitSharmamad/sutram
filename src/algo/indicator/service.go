package indicator

import (
	"algo/src/shared"
	"context"
)

type Service struct{}

func NewServiceProv() *Service {
	return &Service{}
}

func (s *Service) GetIndicatorData(
	ctx context.Context,
	name string,
	period int,
	data *shared.StockData,
) []float64 {

	calculator := GetIndFactory(name)
	println(calculator)
	var results []float64
	config := make(map[string]any, 0)
	config["period"] = period

	results = (calculator).CalculateAll(*data, config, "close")
	return results
}
