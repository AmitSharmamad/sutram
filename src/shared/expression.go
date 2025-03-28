package shared

type ExpressionService interface {
	Eval(
		expression string,
		data *StockData,
	) bool
}
