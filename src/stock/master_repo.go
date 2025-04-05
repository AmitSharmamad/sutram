package stock

import (
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

func (r *MasterRepo) Insert(master *shared.InstrumentMaster) error {
	_, err := r.db.GetDB().Exec(`
		INSERT INTO sutram.instrument_master (
			id, exchange, segment, instrument
		) VALUES (
			?, ?, ?, ?
		);
		`,
		master.ID, master.Exchange, master.Segment, master.Instrument,
	)
	return err
}

func (r *MasterRepo) GetInstrument(name string) (*shared.InstrumentMaster, error) {
	master := make([]shared.InstrumentMaster, 0)
	err := r.db.Select(&master, "SELECT * FROM sutram.instrument_master WHERE instrument = $1 LIMIT 1", name)
	if len(master) == 0 {
		return nil, nil
	}
	return &master[0], err
}

func (r *MasterRepo) GetAll(exchange, segment string) ([]shared.InstrumentMaster, error) {
	master := make([]shared.InstrumentMaster, 0)
	query := "SELECT * FROM sutram.instrument_master WHERE 1=1"
	args := make([]interface{}, 0)

	if exchange != "" {
		query += " AND exchange = ?"
		args = append(args, exchange)
	}
	if segment != "" {
		query += " AND segment = ?"
		args = append(args, segment)
	}

	err := r.db.Select(&master, query, args...)
	if err != nil {
		return nil, err
	}
	return master, nil
}

func (r *MasterRepo) GetLatestID() (uint64, error) {
	ids := make([]uint64, 0)
	err := r.db.Select(&ids, "SELECT id FROM sutram.instrument_master ORDER BY id DESC LIMIT 1;")
	if len(ids) > 0 {
		return ids[0], err
	}
	return 0, err
}
