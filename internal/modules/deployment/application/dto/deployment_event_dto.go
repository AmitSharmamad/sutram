package dto

import (
	"time"

	"github.com/google/uuid"
)

// DeploymentEventType represents the type of deployment event
type DeploymentEventType string

const (
	// Strategy evaluation events
	EventTypeEvaluationStarted   DeploymentEventType = "evaluation_started"
	EventTypeEntrySignalDetected DeploymentEventType = "entry_signal_detected"
	EventTypeExitSignalDetected  DeploymentEventType = "exit_signal_detected"
	EventTypeNoSignal            DeploymentEventType = "no_signal"

	// Position management events
	EventTypePositionOpened DeploymentEventType = "position_opened"
	EventTypePositionClosed DeploymentEventType = "position_closed"
	EventTypePositionUpdate DeploymentEventType = "position_update"

	// Risk management events
	EventTypeStopLossTriggered     DeploymentEventType = "stop_loss_triggered"
	EventTypeTargetProfitTriggered DeploymentEventType = "target_profit_triggered"
	EventTypeTrailingStopUpdated   DeploymentEventType = "trailing_stop_updated"
	EventTypeEODSquareOff          DeploymentEventType = "eod_square_off"

	// Deployment status events
	EventTypeDeploymentError    DeploymentEventType = "deployment_error"
	EventTypeEntryWindowClosed  DeploymentEventType = "entry_window_closed"
	EventTypeMarketDataReceived DeploymentEventType = "market_data_received"
)

// DeploymentEvent represents a real-time deployment event
type DeploymentEvent struct {
	EventType    DeploymentEventType `json:"event_type"`
	DeploymentID uuid.UUID           `json:"deployment_id"`
	StrategyID   uuid.UUID           `json:"strategy_id"`
	InstrumentID int                 `json:"instrument_id"`
	Timestamp    time.Time           `json:"timestamp"`
	Data         interface{}         `json:"data"`
	Metadata     EventMetadata       `json:"metadata"`
}

// EventMetadata contains additional context about the event
type EventMetadata struct {
	UserID           uuid.UUID `json:"user_id"`
	StrategyName     string    `json:"strategy_name"`
	InstrumentSymbol string    `json:"instrument_symbol,omitempty"`
	MarketPrice      *float64  `json:"market_price,omitempty"`
}

// EvaluationStartedData contains data for evaluation started events
type EvaluationStartedData struct {
	MarketPrice float64   `json:"market_price"`
	Timestamp   time.Time `json:"timestamp"`
	Timeframe   string    `json:"timeframe"`
}

// SignalDetectedData contains data for signal detection events
type SignalDetectedData struct {
	MarketPrice float64   `json:"market_price"`
	SignalTime  time.Time `json:"signal_time"`
	Expression  string    `json:"expression"`
	SignalValue bool      `json:"signal_value"`
}

// PositionOpenedData contains data for position opened events
type PositionOpenedData struct {
	PositionID        string    `json:"position_id"`
	EntryPrice        float64   `json:"entry_price"`
	Quantity          float64   `json:"quantity"`
	Direction         string    `json:"direction"`
	EntryTime         time.Time `json:"entry_time"`
	StopLossPrice     *float64  `json:"stop_loss_price,omitempty"`
	TargetProfitPrice *float64  `json:"target_profit_price,omitempty"`
}

// PositionClosedData contains data for position closed events
type PositionClosedData struct {
	PositionID string    `json:"position_id"`
	ExitPrice  float64   `json:"exit_price"`
	ExitTime   time.Time `json:"exit_time"`
	ExitReason string    `json:"exit_reason"`
	PNL        float64   `json:"pnl"`
	Duration   string    `json:"duration"`
	EntryPrice float64   `json:"entry_price"`
	Quantity   float64   `json:"quantity"`
}

// PositionUpdateData contains data for position update events
type PositionUpdateData struct {
	PositionID        string   `json:"position_id"`
	CurrentPrice      float64  `json:"current_price"`
	UnrealizedPNL     float64  `json:"unrealized_pnl"`
	StopLossPrice     *float64 `json:"stop_loss_price,omitempty"`
	TargetProfitPrice *float64 `json:"target_profit_price,omitempty"`
}

// RiskEventData contains data for risk management events
type RiskEventData struct {
	PositionID            string   `json:"position_id"`
	TriggerPrice          float64  `json:"trigger_price"`
	CurrentPrice          float64  `json:"current_price"`
	TriggerReason         string   `json:"trigger_reason"`
	PreviousStopLossPrice *float64 `json:"previous_stop_loss_price,omitempty"`
	NewStopLossPrice      *float64 `json:"new_stop_loss_price,omitempty"`
}

// ErrorEventData contains data for error events
type ErrorEventData struct {
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code,omitempty"`
	Context      string `json:"context"`
}

// MarketDataReceivedData contains data for market data events
type MarketDataReceivedData struct {
	Price     float64   `json:"price"`
	Volume    int64     `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
	OHLC      OHLC      `json:"ohlc"`
}

// OHLC represents OHLC data
type OHLC struct {
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

// TimeWindowEventData contains data for time window events
type TimeWindowEventData struct {
	WindowType   string     `json:"window_type"`   // "entry" or "exit"
	WindowStatus string     `json:"window_status"` // "opened", "closed"
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	CurrentTime  time.Time  `json:"current_time"`
}
