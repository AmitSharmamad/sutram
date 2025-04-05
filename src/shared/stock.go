package shared

import (
	"time"
)

type Exchange string

const (
	NSE Exchange = "NSE"
	BSE Exchange = "BSE"
	MCX Exchange = "MCX"
)

type Segment string

const (
	EQ  = "EQ"
	FNO = "FNO"
)

type InstrumentMaster struct {
	ID             uint64    `json:"id" db:"id"`
	Exchange       Exchange  `json:"exchange" db:"exchange"`
	Segment        Segment   `json:"segment" db:"segment"`
	Symbol         string    `json:"symbol" db:"symbol"`
	Instrument     string    `json:"instrument" db:"instrument"`
	InstrumentID   uint64    `json:"instrument_id" db:"instrument_id"`
	InstrumentName string    `json:"instrument_name" db:"instrument_name"`
	StartDate      time.Time `json:"start_date" db:"start_date"`
	Expiry         time.Time `json:"expiry" db:"expiry"`
}

type Instrument struct {
	ID        uint64    `json:"id" db:"id"`
	Open      float64   `json:"open" db:"open"`
	High      float64   `json:"high" db:"high"`
	Close     float64   `json:"close" db:"close"`
	Low       float64   `json:"low" db:"low"`
	Volume    int64     `json:"volume" db:"volume"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Stock struct {
	Timestamp time.Time `json:"timestamp" db:"timestamp"`

	Open  float64 `json:"open" db:"open"`
	High  float64 `json:"high" db:"high"`
	Close float64 `json:"close" db:"close"`
	Low   float64 `json:"low" db:"low"`

	Volume int64 `json:"volume" db:"volume"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// OHLCV represents a single candlestick data point
type OHLCV struct {
	ID         int64     `json:"id" db:"id"`
	Symbol     string    `json:"symbol" db:"symbol"`
	Instrument string    `json:"instrument" db:"instrument"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	Open       float64   `json:"open" db:"open"`
	High       float64   `json:"high" db:"high"`
	Low        float64   `json:"low" db:"low"`
	Close      float64   `json:"close" db:"close"`
	Volume     int64     `json:"volume" db:"volume"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type StockData []Stock

type StockColumn struct {
	Timestamp []time.Time `json:"timestamps" db:"timestamp"`

	Open  []float64 `json:"opens" db:"open"`
	High  []float64 `json:"highs" db:"high"`
	Close []float64 `json:"closes" db:"close"`
	Low   []float64 `json:"lows" db:"low"`

	Volume []int64 `json:"volumes" db:"volume"`

	CreatedAt []time.Time `json:"created_ats" db:"created_at"`
}

type StockIndicator struct {
	StockColumn
	Indicators map[string][]float64 `json:"indicators" db:"indicators"`
}

type DateRange string

const (
	LastDay     DateRange = "last_day"
	LastWeek    DateRange = "last_week"
	LastMonth   DateRange = "last_month"
	Last6Months DateRange = "last_6_months"
	LastYear    DateRange = "last_year"
	YearToDate  DateRange = "year_to_date"
	Last3Years  DateRange = "last_3_years"
	Last5Years  DateRange = "last_5_years"
	All         DateRange = ""
)

type StockGranularity string

const (
	StockGranularity1m  StockGranularity = "1m"
	StockGranularity3m  StockGranularity = "3m"
	StockGranularity5m  StockGranularity = "5m"
	StockGranularity15m StockGranularity = "15m"
	StockGranularity30m StockGranularity = "30m"
	StockGranularity1h  StockGranularity = "1h"
	StockGranularity3h  StockGranularity = "3h"
	StockGranularity1d  StockGranularity = "1d"
	StockGranularity1M  StockGranularity = "1M"
)

type StockService interface {
	Converter(sd StockData) *StockColumn
	Get(
		id int32,
		startPeriod CustomTime,
		endPeriod CustomTime,
		timeFrame string,
		lookBack int32,
	) *StockData
	ConvertToDateRange(dateRangeStr string) DateRange
	ConvertToGranularity(granularityStr string) StockGranularity
	Insert(master InstrumentMaster, instruments []*OHLCV) error
}

type InstrumentService interface {
	InsertMaster(master InstrumentMaster)
	InsertData(master InstrumentMaster, instruments []Instrument) error
}
