package shared

import "time"

type Stock struct {
	Timestamp time.Time `json:"timestamp" db:"timestamp"`

	Open  float64 `json:"open" db:"open"`
	High  float64 `json:"high" db:"high"`
	Close float64 `json:"close" db:"close"`
	Low   float64 `json:"low" db:"low"`

	Volume int64 `json:"volume" db:"volume"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
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

type StockService interface {
	Converter(sd StockData) *StockColumn
	Get(name string, dateRangeStr string, lastNRecords int32) *StockData
	ConvertToDateRange(dateRangeStr string) DateRange
}
