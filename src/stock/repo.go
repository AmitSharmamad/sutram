package stock

import (
	"fmt"
	"time"

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

func (repo *Repo) Get(
	instrumentId int32,
	startPeriod time.Time,
	endPeriod time.Time,
	timeFrame string,
) (*shared.StockData, error) {
	var data shared.StockData

	query := fmt.Sprintf(`
		SELECT
			CASE
				WHEN '%s' = '1m' THEN toStartOfInterval(timestamp, INTERVAL 1 MINUTE)
				WHEN '%s' = '3m' THEN toStartOfInterval(timestamp, INTERVAL 3 MINUTE)
				WHEN '%s' = '5m' THEN toStartOfInterval(timestamp, INTERVAL 5 MINUTE)
				WHEN '%s' = '15m' THEN toStartOfInterval(timestamp, INTERVAL 15 MINUTE)
				WHEN '%s' = '30m' THEN toStartOfInterval(timestamp, INTERVAL 30 MINUTE)
				WHEN '%s' = '1h' THEN toStartOfInterval(timestamp, INTERVAL 1 HOUR)
				WHEN '%s' = '3h' THEN toStartOfInterval(timestamp, INTERVAL 3 HOUR)
				WHEN '%s' = '1d' THEN toStartOfInterval(timestamp, INTERVAL 1 DAY)
				WHEN '%s' = '1M' THEN toStartOfInterval(timestamp, INTERVAL 1 MONTH)
			END AS interval_time,
			id,
			MIN(low) AS low,
			MAX(high) AS high,
			first_value(open) AS open,
			last_value(close) AS close,
			SUM(volume) AS volume
		FROM
			instrument
		WHERE
			id = %s
			AND timestamp BETWEEN toDateTime64('%s', 6) AND toDateTime64('%s', 6)
		GROUP BY
			id,
			interval_time
		ORDER BY
			interval_time ASC;
	`,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		instrumentId,
		startPeriod.Format("2006-01-02 15:04:05"),
		endPeriod.Format("2006-01-02 15:04:05"),
	)

	err := repo.db.Select(&data, query, instrumentId)
	return &data, err
}

func(repo *Repo) GetLookBack(id int32, timeFrame string, startTime time.Time, lookback int32) (*shared.StockData, error) {
	var data shared.StockData

	query := fmt.Sprintf(`
		SELECT
			CASE
				WHEN '%s' = '1m' THEN toStartOfInterval(timestamp, INTERVAL 1 MINUTE)
				WHEN '%s' = '3m' THEN toStartOfInterval(timestamp, INTERVAL 3 MINUTE)
				WHEN '%s' = '5m' THEN toStartOfInterval(timestamp, INTERVAL 5 MINUTE)
				WHEN '%s' = '15m' THEN toStartOfInterval(timestamp, INTERVAL 15 MINUTE)
				WHEN '%s' = '30m' THEN toStartOfInterval(timestamp, INTERVAL 30 MINUTE)
				WHEN '%s' = '1h' THEN toStartOfInterval(timestamp, INTERVAL 1 HOUR)
				WHEN '%s' = '3h' THEN toStartOfInterval(timestamp, INTERVAL 3 HOUR)
				WHEN '%s' = '1d' THEN toStartOfInterval(timestamp, INTERVAL 1 DAY)
				WHEN '%s' = '1M' THEN toStartOfInterval(timestamp, INTERVAL 1 MONTH)
			END AS interval_time,
			id,
			MIN(low) AS low,
			MAX(high) AS high,
			first_value(open) AS open,
			last_value(close) AS close,
			SUM(volume) AS volume
		FROM
			instrument
		WHERE
			id = %d
			AND timestamp <= toDateTime64('%s', 6)
		GROUP BY
			id,
			interval_time
		ORDER BY
			interval_time DESC
		LIMIT %d;
	`,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		timeFrame,
		id,
		startTime.Format("2006-01-02 15:04:05"),
		lookback,
	)

	err := repo.db.Select(&data, query)
	return &data, err
}

func (repo *Repo) Insert(instruments []*shared.OHLCV) error {
	_, err := repo.db.Insert(`
		INSERT INTO sutram.instrument (
			id, symbol, open, high, low, close, volume, timestamp, created_at, updated_at
		) VALUES (
			:id, :instrument, :open, :high, :low, :close, :volume, :timestamp, :created_at, :updated_at
		)`,
		instruments,
	)
	return err
}
