package stock

import (
	"fmt"

	"algo/src/core/db"
	"algo/src/shared"
)

type MasterRepo struct {
	Name string
	db   db.Database
}

func NewMasterRepoProv(factory *db.DBFactory) *MasterRepo {
	return &MasterRepo{
		Name: "instrument_master",
		db:   factory.Get(string(db.CLICKHOUSE)),
	}
}

/**
ID             uint64    `json:"id" db:"id"`
Exchange       Exchange  `json:"exchange" db:"exchange"`
Segment        Segment   `json:"segment" db:"segment"`
Symbol         string    `json:"symbol" db:"symbol"`
Instrument     string    `json:"instrument" db:"instrument"`
InstrumentID   uint64    `json:"instrument_id" db:"instrument_id"`
InstrumentName string    `json:"instrument_name" db:"instrument_name"`
StartDate      time.Time `json:"start_date" db:"start_date"`
Expiry         time.Time `json:"expiry" db:"expiry"`

*/

func (r *MasterRepo) Insert(master *shared.InstrumentMaster) error {
	_, err := r.db.GetDb().NamedQuery(`
		INSERT INTO sutram.instrument_master (
			id, exchange, segment, instrument, instrument_id
		) VALUES (
			:id, :exchange, :segment, :instrument, :instrument_id
		)
		`,
		master,
	)
	fmt.Printf("%+v\n", err)
	return err
}

func (r *MasterRepo) GetInstrument(name string) (uint64, error) {
	rows, err := r.db.GetDb().NamedQuery("SELECT id FROM sutram.instrument WHERE symbol = :symbol;", struct{ symbol string }{symbol: name})
	fmt.Printf("coming from here %+v\n", rows)
	return 1, err
}
