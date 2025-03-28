package expression

import "algo/src/shared"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (service *Service) Eval(
	expression string,
	data *shared.StockData,
) bool {
	return NewExpressionStatement(expression).Eval(*data)
}
