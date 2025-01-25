package stock

import (
	"algo/src/core/db"
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

func (this *Repo) Insert() {}

func (this *Repo) Select() (any, error) {
	id := "TESLA"

	var data StockData

	this.db.Select(
		&data,
		"SELECT * FROM "+this.TableName(id)+" LIMIT 1;",
	)

	println(data)

	return data, nil
}
func (this *Repo) Update() {}
func (this *Repo) Delete() {}

func (this *Repo) TableName(id string) string {
	return "sutram.stock_TESLA"
}

func (this *Repo) Model() []interface{} {
	var stocks []interface{}

	stocks = append(stocks, Stock{})

	return stocks
}
