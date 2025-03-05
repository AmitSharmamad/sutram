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

// DateRange represents common time periods for filtering
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

func (repo *Repo) Get(instrumentId int32, dateRange shared.DateRange, lastNRecords int32) (*shared.StockData, error) {
    var data shared.StockData

    // Base query
    query := `
        SELECT timestamp, open, high, close, low, volume, created_at
        FROM sutram.instrument
        WHERE id = $1
    `

    // Add date filtering based on the requested range
    switch dateRange {
        case shared.LastDay:
            query += " AND timestamp >= now() - INTERVAL '1 day'"
        case shared.LastWeek:
            query += " AND timestamp >= now() - INTERVAL '7 days'"
        case shared.LastMonth:
            query += " AND timestamp >= now() - INTERVAL '1 month'"
        case shared.Last6Months:
            query += " AND timestamp >= now() - INTERVAL '6 months'"
        case shared.LastYear:
            query += " AND timestamp >= now() - INTERVAL '1 year'"
        case shared.YearToDate:
            query += " AND timestamp >= date_trunc('year', now())"
        case shared.Last3Years:
            query += " AND timestamp >= now() - INTERVAL '3 years'"
        case shared.Last5Years:
            query += " AND timestamp >= now() - INTERVAL '5 years'"
        case shared.All:
            println(query)
    }

    // Add ordering and limit
    query += " ORDER BY timestamp DESC"

    // Execute the query
    err := repo.db.Select(&data, query, instrumentId)

    return &data, err
}
