package strategy

import (
	"algo/src/algo"
	"algo/src/data"
)

type Strategy struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Position    Position    `json:"position"`
	Entry       Entry       `json:"entry"`
	Exit        Exit        `json:"exit"`
	TestParams  TradeParams `json:"test_params"`
}

type Position struct {
	Type        string   `json:"type"` // Type can be either "BUY" or "SELL".
	Instruments []string `json:"instruments"`
	Chart       Chart    `json:"chart"`
}

type Chart struct {
	Type      string `json:"type"`      // Type can be one of: "CANDLESTICK", "HEIKEN-ASHI", "RENKO".
	Timeframe string `json:"timeframe"` // Timeframe can be one of: "1m", "3m", "5m", "10m", "15m", "30m", "1h", "1d".
	Unknown   string `json:"unknown"`   // Unknown can be one of: "MIS", "CNC", "NRML".
	Qty       int    `json:"qty"`
}

type Entry struct {
	Expression algo.Expression `json:"expression"`
	Prompt     string          `json:"prompt"`
}

type Exit struct {
	Expression       algo.Expression `json:"expression"`
	Prompt           string          `json:"prompt"`
	TargetProfit     float64         `json:"target_profit"`
	StopLoss         float64         `json:"stop_loss"`
	TrailingStopLoss float64         `json:"trailing_stop_loss"`
	Type             string          `json:"type"`
}

type TradeParams struct {
	Capital        int64  `json:"capital"`
	MaxTrades      int    `json:"max_trades"`
	EntryStartTime int64  `json:"entry_start_time"`
	EntryEndTime   int64  `json:"entry_end_time"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
}

func Get(id string) *Strategy {
	return &Strategy{}
}

func (strategy *Strategy) GetEntryDataRange() int {
	return 0
}

func (strategy *Strategy) TryEnter(data data.StockData) bool {
	/// eval entry expression
	return true
}

func (strategy *Strategy) GetExitDataRange() int {
	return 0
}

func (strategy *Strategy) TryExit(data data.StockData) bool {
	/// eval exit expression
	return true
}
