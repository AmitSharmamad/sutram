package shared

type ExpressionService interface {
	Eval(
		expression string,
		data *StockData,
	) bool
	GetDataRange(strategy *Strategy) int32

	GetEntryDataRange(strategy *Strategy) int32
	GetExitDataRange(strategy *Strategy) int32
}
