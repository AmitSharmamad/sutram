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

func (service *Service) GetEntryDataRange(strategy *shared.Strategy) int32 {
	return NewExpressionStatement(strategy.EntryExpression).GetDataRange()
}

func (service *Service) GetExitDataRange(strategy *shared.Strategy) int32 {
	return NewExpressionStatement(strategy.ExitExpression).GetDataRange()
}
