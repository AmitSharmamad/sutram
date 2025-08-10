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
	ExchangeSecurityID   int32                `json:"exchange_security_id"`
	Exchange             string               `json:"exchange"`
	Segment              model.Segment        `json:"segment"`
	Symbol               string               `json:"symbol"`
	Name                 string               `json:"name"`
	StartDate            time.Time            `json:"start_date"`
	Expiry               *time.Time           `json:"expiry,omitempty"`
	ExmExchID            int32                `json:"exm_exch_id"`
	SmstSecurityID       int32                `json:"smst_security_id"`
	InstrumentName       string               `json:"instrument_name"`
	TradingSymbol        string               `json:"trading_symbol"`
	LotUnits             int32                `json:"lot_units"`
	StrikePrice          float64              `json:"strike_price"`
	OptionType           string               `json:"option_type"`
	TickSize             float64              `json:"tick_size"`
	ExpiryFlag           string               `json:"expiry_flag"`
	ExchInstrumentType   model.InstrumentType `json:"exch_instrument_type"`
	SymbolName           string               `json:"symbol_name"`
	UnderlyingSecurityID int32                `json:"underlying_security_id"`
	UnderlyingSymbol     string               `json:"underlying_symbol"`
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
	Exchange             string               `json:"exchange" binding:"required"`
	Segment              model.Segment        `json:"segment" binding:"required"`
	Symbol               string               `json:"symbol" binding:"required"`
	Name                 string               `json:"name" binding:"required"`
	StartDate            time.Time            `json:"start_date" binding:"required"`
	Expiry               *time.Time           `json:"expiry,omitempty"`
	ExmExchID            int32                `json:"exm_exch_id"`
	SmstSecurityID       int32                `json:"smst_security_id"`
	InstrumentName       string               `json:"instrument_name"`
	TradingSymbol        string               `json:"trading_symbol"`
	LotUnits             int32                `json:"lot_units"`
	StrikePrice          float64              `json:"strike_price"`
	OptionType           string               `json:"option_type"`
	TickSize             float64              `json:"tick_size"`
	ExpiryFlag           string               `json:"expiry_flag"`
	ExchInstrumentType   model.InstrumentType `json:"exch_instrument_type"`
	SymbolName           string               `json:"symbol_name"`
	UnderlyingSecurityID int32                `json:"underlying_security_id"`
	UnderlyingSymbol     string               `json:"underlying_symbol"`
}

// UpdateInstrumentMasterRequest represents request to update an instrument master
type UpdateInstrumentMasterRequest struct {
	Exchange             string               `json:"exchange"`
	Segment              model.Segment        `json:"segment"`
	Symbol               string               `json:"symbol"`
	Name                 string               `json:"name"`
	StartDate            time.Time            `json:"start_date"`
	Expiry               *time.Time           `json:"expiry,omitempty"`
	ExmExchID            int32                `json:"exm_exch_id"`
	SmstSecurityID       int32                `json:"smst_security_id"`
	InstrumentName       string               `json:"instrument_name"`
	TradingSymbol        string               `json:"trading_symbol"`
	LotUnits             int32                `json:"lot_units"`
	StrikePrice          float64              `json:"strike_price"`
	OptionType           string               `json:"option_type"`
	TickSize             float64              `json:"tick_size"`
	ExpiryFlag           string               `json:"expiry_flag"`
	ExchInstrumentType   model.InstrumentType `json:"exch_instrument_type"`
	SymbolName           string               `json:"symbol_name"`
	UnderlyingSecurityID int32                `json:"underlying_security_id"`
	UnderlyingSymbol     string               `json:"underlying_symbol"`
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
	Exchange       string                   `form:"exchange"`
	Segment        model.Segment            `form:"segment"`
	Market         model.Market             `form:"market"`
	InstrumentType model.InstrumentCategory `form:"instrumentType"`
	Contract       model.Contract           `form:"contract"`
	Query          string                   `form:"query"`
	Page           int                      `form:"page,default=1"`
	PageSize       int                      `form:"pageSize,default=10"`
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

// DatabaseFilter represents the resolved database filter criteria
type DatabaseFilter struct {
	Exchange       string               `json:"exchange,omitempty"`
	Segment        model.Segment        `json:"segment,omitempty"`
	InstrumentType model.InstrumentType `json:"instrument_type,omitempty"`
	Query          string               `json:"query,omitempty"`
	Page           int                  `json:"page"`
	PageSize       int                  `json:"page_size"`
}

// ConvertToDatabase converts the filter DTO to database filter criteria
func (f *InstrumentMasterFilterDTO) ConvertToDatabase() *DatabaseFilter {
	dbFilter := &DatabaseFilter{
		Exchange: f.Exchange, // Use the exchange from the filter or default to NSE if not provided
		Query:    f.Query,
		Page:     f.Page,
		PageSize: f.PageSize,
	}

	// If no exchange is specified, default to NSE for the mapping logic
	exchange := f.Exchange
	if exchange == "" {
		exchange = "NSE"
	}
	dbFilter.Exchange = exchange

	// Apply conversion logic based on the mapping table
	switch {
	// EQ + NSE + Stock + Spot -> Equities + Stocks (Cash Market)
	case f.Market == model.MarketEQ && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractSpot:
		dbFilter.Segment = model.SegmentEquities
		dbFilter.InstrumentType = model.InstrumentTypeStockCash

	// EQ + NSE + Index -> NA (no specific instrument type filtering for index in EQ)
	case f.Market == model.MarketEQ && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryIndex:
		dbFilter.Segment = model.SegmentEquities
		// No specific instrument type filter for EQ Index

	// F&O + NSE + Stock + Spot -> Indices + Index Spot
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractSpot:
		dbFilter.Segment = model.SegmentIndices
		dbFilter.InstrumentType = model.InstrumentTypeIndexSpot

	// F&O + NSE + Stock + Futures -> Derivatives + Stock Futures
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractFutures:
		dbFilter.Segment = model.SegmentDerivatives
		dbFilter.InstrumentType = model.InstrumentTypeStockFutures

	// F&O + NSE + Index + Futures -> Derivatives + Index Futures
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryIndex && f.Contract == model.ContractFutures:
		dbFilter.Segment = model.SegmentDerivatives
		dbFilter.InstrumentType = model.InstrumentTypeIndexFutures

	// F&O + NSE + Stock + Options -> Derivatives + Stock Options
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractOptions:
		dbFilter.Segment = model.SegmentDerivatives
		dbFilter.InstrumentType = model.InstrumentTypeStockOptions

	// F&O + NSE + Index + Options -> Derivatives + Index Options
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryIndex && f.Contract == model.ContractOptions:
		dbFilter.Segment = model.SegmentDerivatives
		dbFilter.InstrumentType = model.InstrumentTypeIndexOptions

	// F&O + MCX + Stock + Futures -> Commodities + Commodity Futures
	case f.Market == model.MarketFO && exchange == "MCX" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractFutures:
		dbFilter.Segment = model.SegmentCommodities
		dbFilter.InstrumentType = model.InstrumentTypeCommodityFutures

	// F&O + MCX + Stock + Options -> Commodities + Commodity Options
	case f.Market == model.MarketFO && exchange == "MCX" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractOptions:
		dbFilter.Segment = model.SegmentCommodities
		dbFilter.InstrumentType = model.InstrumentTypeCommodityOptions

	default:
		// Handle partial filters or use individual field mappings
		if f.Segment != "" {
			dbFilter.Segment = f.Segment
		}
		// For other cases, we might need to add more specific logic or leave fields empty
	}

	return dbFilter
}
