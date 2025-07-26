package dto

import (
	"algo/internal/modules/instrument/domain/model"
	"time"
)

// InstrumentDTO represents instrument data for API responses
type InstrumentDTO struct {
	ExchangeSecurityID int32     `json:"exchange_security_id"`
	Open               float64   `json:"open"`
	High               float64   `json:"high"`
	Low                float64   `json:"low"`
	Close              float64   `json:"close"`
	Volume             float64   `json:"volume"`
	Timestamp          time.Time `json:"timestamp"`
}

// InstrumentMasterDTO represents instrument master data for API responses
type InstrumentMasterDTO struct {
	ExchangeSecurityID int32                `json:"exchange_security_id"`
	Exchange           string               `json:"exchange"`
	Segment            model.Segment        `json:"segment"`
	Symbol             string               `json:"symbol"`
	Name               string               `json:"name"`
	StartDate          time.Time            `json:"start_date"`
	Expiry             time.Time            `json:"expiry,omitempty"`
	ExmExchID          int32                `json:"exm_exch_id"`
	SmstSecurityID     int32                `json:"smst_security_id"`
	InstrumentName     string               `json:"instrument_name"`
	TradingSymbol      string               `json:"trading_symbol"`
	LotUnits           int32                `json:"lot_units"`
	StrikePrice        float64              `json:"strike_price"`
	OptionType         string               `json:"option_type"`
	TickSize           float64              `json:"tick_size"`
	ExpiryFlag         string               `json:"expiry_flag"`
	ExchInstrumentType model.InstrumentType `json:"exch_instrument_type"`
	SymbolName         string               `json:"symbol_name"`
}

// InstrumentWithMetadataDTO combines instrument data with its metadata
type InstrumentWithMetadataDTO struct {
	Instrument       *InstrumentDTO       `json:"instrument"`
	InstrumentMaster *InstrumentMasterDTO `json:"instrument_master"`
}

// CreateInstrumentRequest represents request to create an instrument
type CreateInstrumentRequest struct {
	ExchangeSecurityID int32     `json:"exchange_security_id" binding:"required"`
	Open               float64   `json:"open" binding:"required"`
	High               float64   `json:"high" binding:"required"`
	Low                float64   `json:"low" binding:"required"`
	Close              float64   `json:"close" binding:"required"`
	Volume             float64   `json:"volume" binding:"required"`
	Timestamp          time.Time `json:"timestamp" binding:"required"`
}

// UpdateInstrumentRequest represents request to update an instrument
type UpdateInstrumentRequest struct {
	Open      float64   `json:"open" binding:"required"`
	High      float64   `json:"high" binding:"required"`
	Low       float64   `json:"low" binding:"required"`
	Close     float64   `json:"close" binding:"required"`
	Volume    float64   `json:"volume" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

// CreateInstrumentMasterRequest represents request to create an instrument master
type CreateInstrumentMasterRequest struct {
	Exchange           string               `json:"exchange" binding:"required"`
	Segment            model.Segment        `json:"segment" binding:"required"`
	Symbol             string               `json:"symbol" binding:"required"`
	Name               string               `json:"name" binding:"required"`
	StartDate          time.Time            `json:"start_date" binding:"required"`
	Expiry             time.Time            `json:"expiry,omitempty"`
	ExmExchID          int32                `json:"exm_exch_id"`
	SmstSecurityID     int32                `json:"smst_security_id"`
	InstrumentName     string               `json:"instrument_name"`
	TradingSymbol      string               `json:"trading_symbol"`
	LotUnits           int32                `json:"lot_units"`
	StrikePrice        float64              `json:"strike_price"`
	OptionType         string               `json:"option_type"`
	TickSize           float64              `json:"tick_size"`
	ExpiryFlag         string               `json:"expiry_flag"`
	ExchInstrumentType model.InstrumentType `json:"exch_instrument_type"`
	SymbolName         string               `json:"symbol_name"`
}

// UpdateInstrumentMasterRequest represents request to update an instrument master
type UpdateInstrumentMasterRequest struct {
	Exchange           string               `json:"exchange"`
	Segment            model.Segment        `json:"segment"`
	Symbol             string               `json:"symbol"`
	Name               string               `json:"name"`
	StartDate          time.Time            `json:"start_date"`
	Expiry             time.Time            `json:"expiry,omitempty"`
	ExmExchID          int32                `json:"exm_exch_id"`
	SmstSecurityID     int32                `json:"smst_security_id"`
	InstrumentName     string               `json:"instrument_name"`
	TradingSymbol      string               `json:"trading_symbol"`
	LotUnits           int32                `json:"lot_units"`
	StrikePrice        float64              `json:"strike_price"`
	OptionType         string               `json:"option_type"`
	TickSize           float64              `json:"tick_size"`
	ExpiryFlag         string               `json:"expiry_flag"`
	ExchInstrumentType model.InstrumentType `json:"exch_instrument_type"`
	SymbolName         string               `json:"symbol_name"`
}

// GetInstrumentsByTimeRangeRequest represents a request to fetch instruments by time range
type GetInstrumentsByTimeRangeRequest struct {
	ExchangeSecurityID int32     `json:"exchange_security_id" binding:"required"`
	StartTime          time.Time `json:"start_time" binding:"required"`
	EndTime            time.Time `json:"end_time" binding:"required"`
	TimeFrame          string    `json:"time_frame" binding:"omitempty,oneof=1s 1m 3m 5m 15m 30m 1h 3h 1d 1M"`
}

// InstrumentMasterFilterDTO represents filtering and pagination options for instrument masters
type InstrumentMasterFilterDTO struct {
	Exchange string        `form:"exchange"`
	Segment  model.Segment `form:"segment"`
	Query    string        `form:"query"`
	Page     int           `form:"page,default=1"`
	PageSize int           `form:"pageSize,default=10"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalCount int         `json:"totalCount"`
	TotalPages int         `json:"totalPages"`
}

// ApiResponse represents a standard API response structure
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// GetInstrumentsByExchangeSecurityIDRequest represents request to get paginated instruments
type GetInstrumentsByExchangeSecurityIDRequest struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}
