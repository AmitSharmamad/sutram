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
		Name: "instrument",
		db:   factory.Get("clickhouse"),
	}
}

func (repo *Repo) GetInstrumentId(instrumentName string) (int32, error) {
	var id int32
	err := repo.db.Get(&id, `
		SELECT instrument_id
		FROM sutram.instrument_master
		WHERE instrument = $1
		LIMIT 1
		`,
		instrumentName,
	)
	return id, err
}

func (repo *Repo) Get(instrumentId int32, lastNRecords int32) (*shared.StockData, error) {
	var data shared.StockData
	err := repo.db.Select(&data, `
		SELECT close
		FROM sutram.instrument
		WHERE id = $1
		LIMIT $2
		`,
		instrumentId,
		lastNRecords,
	)
	return &data, err
}
