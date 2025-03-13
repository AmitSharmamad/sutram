package stock

import (
	"algo/src/core/db"
	"algo/src/shared"
	"fmt"
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

func (repo *Repo) Get(
	instrumentId int32,
	dateRange shared.DateRange,
	granularity shared.StockGranularity,
	lastNRecords int32,
) (*shared.StockData, error) {
	var data shared.StockData
	var query string

	// Set interval value and unit based on granularity
	var intervalValue, intervalUnit string
	switch granularity {
	case shared.StockGranularity1m:
		intervalValue, intervalUnit = "1", "MINUTE"
	case shared.StockGranularity3m:
		intervalValue, intervalUnit = "3", "MINUTE"
	case shared.StockGranularity5m:
		intervalValue, intervalUnit = "5", "MINUTE"
	case shared.StockGranularity15m:
		intervalValue, intervalUnit = "15", "MINUTE"
    case shared.StockGranularity30m:
		intervalValue, intervalUnit = "30", "MINUTE"
    case shared.StockGranularity1h:
		intervalValue, intervalUnit = "1", "HOUR"
	case shared.StockGranularity3h:
		intervalValue, intervalUnit = "3", "HOUR"
	case shared.StockGranularity1d:
		intervalValue, intervalUnit = "1", "DAY"
	case shared.StockGranularity1M:
		intervalValue, intervalUnit = "1", "MONTH"
	default:
		return nil, fmt.Errorf("unsupported granularity: %s", granularity)
	}

	// Base aggregation query
	query = fmt.Sprintf(`
        SELECT
            toStartOfInterval(timestamp, INTERVAL %s %s) AS timestamp,
            argMin(open, timestamp) AS open,
            max(high) AS high,
            min(low) AS low,
            argMax(close, timestamp) AS close,
            sum(volume) AS volume,
            now() AS created_at
        FROM sutram.instrument
        WHERE id = ?
    `, intervalValue, intervalUnit)

	// Add date filtering based on the requested range
	switch dateRange {
	case shared.LastDay:
		query += " AND timestamp >= now() - INTERVAL 1 DAY"
	case shared.LastWeek:
		query += " AND timestamp >= now() - INTERVAL 7 DAY"
	case shared.LastMonth:
		query += " AND timestamp >= now() - INTERVAL 1 MONTH"
	case shared.Last6Months:
		query += " AND timestamp >= now() - INTERVAL 6 MONTH"
	case shared.LastYear:
		query += " AND timestamp >= now() - INTERVAL 1 YEAR"
	case shared.YearToDate:
		query += " AND timestamp >= toStartOfYear(now())"
	case shared.Last3Years:
		query += " AND timestamp >= now() - INTERVAL 3 YEAR"
	case shared.Last5Years:
		query += " AND timestamp >= now() - INTERVAL 5 YEAR"
	case shared.All:
		// No additional filtering
	}

	// Add grouping for aggregated queries
	query += " GROUP BY timestamp"

	// Add ordering and limit
	query += " ORDER BY timestamp ASC"
	// Execute the query without limit
	err := repo.db.Select(&data, query, instrumentId)
	return &data, err
}
