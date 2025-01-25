package strategy

import (
	"algo/src/stock"
	"algo/src/util"
	"time"
)

// Strategy represents the entire structure.
type Strategy struct {
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
	Position    Position    `json:"position" db:"position"`
	Entry       Expression  `json:"entry" db:"entry"`
	Exit        Exit        `json:"exit" db:"exit"`
	TradeParams TradeParams `json:"trade_params" db:"trade_params"`
}

// Position represents the "position" object.
type Position struct {
	Type        string   `json:"type" db:"type"`               // Expected values: "BUY" or "SELL"
	Instruments []string `json:"instruments" db:"instruments"` // List of instrument names
	Chart       Chart    `json:"chart" db:"chart"`
}

// Chart represents the "chart" object within "position".
type Chart struct {
	Type      string `json:"type" db:"type"`           // e.g., "CANDLESTICK", "HEIKEN-ASHI", "RENKO"
	Timeframe string `json:"timeframe" db:"timeframe"` // e.g., "1m", "3m", "5m", "10m", "15m", "30m", "1h", "1d"
	Unknown   string `json:"unknown" db:"unknown"`     // e.g., "MIS", "CNC", "NRML"
	Qty       int    `json:"qty" db:"qty"`             // Quantity value
}

// Expression represents the "entry" object.
type Expression struct {
	Expression string `json:"expression" db:"expression"`
}

// Exit represents the "exit" object.
type Exit struct {
	Expression       string  `json:"expression" db:"expression"`
	TargetProfit     float64 `json:"target_profit" db:"target_profit"`
	StopLoss         float64 `json:"stop_loss" db:"stop_loss"`
	TrailingStopLoss float64 `json:"trailing_stop_loss" db:"trailing_stop_loss"`
	Type             string  `json:"type" db:"type"` // e.g., "PERCENTAGE" or "ABSOLUTE"
}

// TradeParams represents the "trade_params" object.
type TradeParams struct {
	Capital        float64 `json:"capital" db:"capital"`
	MaxTrades      int     `json:"max_trades" db:"max_trades"`
	EntryStartTime int64   `json:"entry_start_time" db:"entry_start_time"` // Could be Unix timestamp
	EntryEndTime   int64   `json:"entry_end_time" db:"entry_end_time"`     // Could be Unix timestamp
	// For date fields, you may choose to use time.Time if your JSON payload is in RFC3339 format.
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
}

func (strategy *Strategy) Validate() {
	path := "../data/llm.json"

	exp := util.GetFileContent(path)

	util.Validate(exp)
}

func Get(id string) *Strategy {
	return &Strategy{}
}

func GetData() stock.StockData {
	return nil
}

func (strategy *Strategy) GetEntryDataRange() int {
	return 0
}

func (strategy *Strategy) TryEnter(data stock.StockData) bool {
	/// eval entry expression
	return true
}

func (strategy *Strategy) GetExitDataRange() int {
	return 0
}

func (strategy *Strategy) TryExit(data stock.StockData) bool {
	/// eval exit expression
	return true
}
