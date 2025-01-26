package strategy

import (
	"algo/src/shared"
	"algo/src/stock"
	"algo/src/util"
)

type Strategy shared.Strategy

func (strategy *Strategy) Validate() {
	path := "../data/llm.json"

	exp := util.GetFileContent(path)

	util.Validate(exp)
}

func Get(id string) *Strategy {
	return &Strategy{}
}

func GetData() stock.StockData {
	return nil
}

func (strategy *Strategy) GetEntryDataRange() int {
	return 0
}

func (strategy *Strategy) TryEnter(data stock.StockData) bool {
	/// eval entry expression
	return true
}

func (strategy *Strategy) GetExitDataRange() int {
	return 0
}

func (strategy *Strategy) TryExit(data stock.StockData) bool {
	/// eval exit expression
	return true
}
