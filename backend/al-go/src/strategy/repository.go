package strategy

import "algo/src/data"

func (strategy *Strategy) validate() {}

func (strategy *Strategy) Create() {
	strategy.validate()
	/// save the strategy to the db
}

func (strategy *Strategy) GetData() data.StockData {
	return nil
}

func (strategy *Strategy) Update() {
	strategy.validate()
	/// update the strategy in the db
}

func (strategy *Strategy) Delete(id string) bool {
	return true
}
