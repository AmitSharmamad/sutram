package stock

import "time"

type Stock struct {
	Timestamp time.Time `json:"timestamp" db:"timestamp"`

	Open  float64 `json:"open" db:"open"`
	High  float64 `json:"high" db:"high"`
	Close float64 `json:"close" db:"close"`
	Low   float64 `json:"low" db:"low"`

	Volume int64 `json:"volume" db:"volume"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type StockData []Stock
