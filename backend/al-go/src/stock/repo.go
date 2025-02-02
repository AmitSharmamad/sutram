package stock

import (
	"algo/src/core/db"
	"algo/src/shared"
)

type Repo struct {
	Name string
	db   db.Database
}

func NewRepoProv(
	factory *db.DBFactory,
) *Repo {
	return &Repo{
		Name: "strategy",
		db:   factory.Get("clickhouse"),
	}
}

func (repo *Repo) Get(lastNRecords int) *shared.StockData {
	var data shared.StockData
	repo.db.Select(&data, `
		SELECT close 
		FROM sutram.stock_TESLA
		LIMIT $1
		`,
		lastNRecords,
	)
	return &data
}
