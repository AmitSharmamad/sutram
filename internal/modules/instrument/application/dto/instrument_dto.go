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
	// Core Identification Columns
	ExchangeSecurityID int32  `json:"exchange_security_id"`
	Exchange           string `json:"exchange"`
	Segment            string `json:"segment"`
	Symbol             string `json:"symbol"`
	Instrument         string `json:"instrument"`

	// Trading Symbols and Names
	InstrumentName string `json:"instrument_name"`

	// ISIN and Security Identifiers
	ISIN                 string  `json:"isin"`
	UnderlyingSecurityID *int32  `json:"underlying_security_id,omitempty"`
	UnderlyingSymbol     *string `json:"underlying_symbol,omitempty"`

	// Date and Expiry Information
	Expiry     time.Time `json:"expiry"`
	ExpiryFlag *string   `json:"expiry_flag,omitempty"`

	// Trading Parameters
	LotUnits       int32    `json:"lot_units"`
	TickSize       float64  `json:"tick_size"`
	StrikePrice    *float64 `json:"strike_price,omitempty"`
	OptionType     *string  `json:"option_type,omitempty"`
	InstrumentType string   `json:"instrument_type"`

	// Order Type Availability Flags
	BracketFlag string `json:"bracket_flag"`
	CoverFlag   string `json:"cover_flag"`

	// Surveillance and Trading Permissions
	AsmGsmFlag       string `json:"asm_gsm_flag"`
	BuySellIndicator string `json:"buy_sell_indicator"`

	// Cover Order Margin Parameters
	BuyCoMinMarginPer    *float64 `json:"buy_co_min_margin_per,omitempty"`
	SellCoMinMarginPer   *float64 `json:"sell_co_min_margin_per,omitempty"`
	BuyCoSlRangeMaxPerc  *float64 `json:"buy_co_sl_range_max_perc,omitempty"`
	SellCoSlRangeMinPerc *float64 `json:"sell_co_sl_range_min_perc,omitempty"`

	// Bracket Order Parameters
	BuyBoMinMarginPer        *float64 `json:"buy_bo_min_margin_per,omitempty"`
	BuyBoSlRangeMaxPerc      *float64 `json:"buy_bo_sl_range_max_perc,omitempty"`
	SellBoSlMinRange         *float64 `json:"sell_bo_sl_min_range,omitempty"`
	BuyBoProfitRangeMaxPerc  *float64 `json:"buy_bo_profit_range_max_perc,omitempty"`
	SellBoProfitRangeMaxPerc *float64 `json:"sell_bo_profit_range_max_perc,omitempty"`

	// Margin Trading Facility
	MtfLeverage *float64 `json:"mtf_leverage,omitempty"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	// Core Identification Columns
	ExchangeSecurityID int32  `json:"exchange_security_id" binding:"required"`
	Exchange           string `json:"exchange" binding:"required"`
	Segment            string `json:"segment" binding:"required"`
	Symbol             string `json:"symbol" binding:"required"`
	Instrument         string `json:"instrument" binding:"required"`

	// Trading Symbols and Names
	InstrumentName string `json:"instrument_name"`

	// ISIN and Security Identifiers
	ISIN                 string  `json:"isin"`
	UnderlyingSecurityID *int32  `json:"underlying_security_id,omitempty"`
	UnderlyingSymbol     *string `json:"underlying_symbol,omitempty"`

	// Date and Expiry Information
	Expiry     time.Time `json:"expiry" binding:"required"`
	ExpiryFlag *string   `json:"expiry_flag,omitempty"`

	// Trading Parameters
	LotUnits       int32    `json:"lot_units" binding:"required"`
	TickSize       float64  `json:"tick_size" binding:"required"`
	StrikePrice    *float64 `json:"strike_price,omitempty"`
	OptionType     *string  `json:"option_type,omitempty"`
	InstrumentType string   `json:"instrument_type" binding:"required"`

	// Order Type Availability Flags
	BracketFlag string `json:"bracket_flag"`
	CoverFlag   string `json:"cover_flag"`

	// Surveillance and Trading Permissions
	AsmGsmFlag       string `json:"asm_gsm_flag"`
	BuySellIndicator string `json:"buy_sell_indicator"`

	// Cover Order Margin Parameters
	BuyCoMinMarginPer    *float64 `json:"buy_co_min_margin_per,omitempty"`
	SellCoMinMarginPer   *float64 `json:"sell_co_min_margin_per,omitempty"`
	BuyCoSlRangeMaxPerc  *float64 `json:"buy_co_sl_range_max_perc,omitempty"`
	SellCoSlRangeMinPerc *float64 `json:"sell_co_sl_range_min_perc,omitempty"`

	// Bracket Order Parameters
	BuyBoMinMarginPer        *float64 `json:"buy_bo_min_margin_per,omitempty"`
	BuyBoSlRangeMaxPerc      *float64 `json:"buy_bo_sl_range_max_perc,omitempty"`
	SellBoSlMinRange         *float64 `json:"sell_bo_sl_min_range,omitempty"`
	BuyBoProfitRangeMaxPerc  *float64 `json:"buy_bo_profit_range_max_perc,omitempty"`
	SellBoProfitRangeMaxPerc *float64 `json:"sell_bo_profit_range_max_perc,omitempty"`

	// Margin Trading Facility
	MtfLeverage *float64 `json:"mtf_leverage,omitempty"`
}

// UpdateInstrumentMasterRequest represents request to update an instrument master
type UpdateInstrumentMasterRequest struct {
	// Core Identification Columns
	Exchange   string `json:"exchange"`
	Segment    string `json:"segment"`
	Symbol     string `json:"symbol"`
	Instrument string `json:"instrument"`

	// Trading Symbols and Names
	InstrumentName string `json:"instrument_name"`

	// ISIN and Security Identifiers
	ISIN                 string  `json:"isin"`
	UnderlyingSecurityID *int32  `json:"underlying_security_id,omitempty"`
	UnderlyingSymbol     *string `json:"underlying_symbol,omitempty"`

	// Date and Expiry Information
	Expiry     time.Time `json:"expiry"`
	ExpiryFlag *string   `json:"expiry_flag,omitempty"`

	// Trading Parameters
	LotUnits       int32    `json:"lot_units"`
	TickSize       float64  `json:"tick_size"`
	StrikePrice    *float64 `json:"strike_price,omitempty"`
	OptionType     *string  `json:"option_type,omitempty"`
	InstrumentType string   `json:"instrument_type"`

	// Order Type Availability Flags
	BracketFlag string `json:"bracket_flag"`
	CoverFlag   string `json:"cover_flag"`

	// Surveillance and Trading Permissions
	AsmGsmFlag       string `json:"asm_gsm_flag"`
	BuySellIndicator string `json:"buy_sell_indicator"`

	// Cover Order Margin Parameters
	BuyCoMinMarginPer    *float64 `json:"buy_co_min_margin_per,omitempty"`
	SellCoMinMarginPer   *float64 `json:"sell_co_min_margin_per,omitempty"`
	BuyCoSlRangeMaxPerc  *float64 `json:"buy_co_sl_range_max_perc,omitempty"`
	SellCoSlRangeMinPerc *float64 `json:"sell_co_sl_range_min_perc,omitempty"`

	// Bracket Order Parameters
	BuyBoMinMarginPer        *float64 `json:"buy_bo_min_margin_per,omitempty"`
	BuyBoSlRangeMaxPerc      *float64 `json:"buy_bo_sl_range_max_perc,omitempty"`
	SellBoSlMinRange         *float64 `json:"sell_bo_sl_min_range,omitempty"`
	BuyBoProfitRangeMaxPerc  *float64 `json:"buy_bo_profit_range_max_perc,omitempty"`
	SellBoProfitRangeMaxPerc *float64 `json:"sell_bo_profit_range_max_perc,omitempty"`

	// Margin Trading Facility
	MtfLeverage *float64 `json:"mtf_leverage,omitempty"`
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
	Segment        string                   `form:"segment"`
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
	Exchange       string `json:"exchange,omitempty"`
	Segment        string `json:"segment,omitempty"`
	InstrumentType string `json:"instrument_type,omitempty"`
	Query          string `json:"query,omitempty"`
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
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
		dbFilter.Segment = "Equities"
		dbFilter.InstrumentType = "Stocks (Cash Market)"

	// EQ + NSE + Index -> NA (no specific instrument type filtering for index in EQ)
	case f.Market == model.MarketEQ && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryIndex:
		dbFilter.Segment = "Equities"
		// No specific instrument type filter for EQ Index

	// F&O + NSE + Stock + Spot -> Indices + Index Spot
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractSpot:
		dbFilter.Segment = "Indices"
		dbFilter.InstrumentType = "Index Spot"

	// F&O + NSE + Stock + Futures -> Derivatives + Stock Futures
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractFutures:
		dbFilter.Segment = "Derivatives"
		dbFilter.InstrumentType = "Stock Futures"

	// F&O + NSE + Index + Futures -> Derivatives + Index Futures
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryIndex && f.Contract == model.ContractFutures:
		dbFilter.Segment = "Derivatives"
		dbFilter.InstrumentType = "Index Futures"

	// F&O + NSE + Stock + Options -> Derivatives + Stock Options
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractOptions:
		dbFilter.Segment = "Derivatives"
		dbFilter.InstrumentType = "Stock Options"

	// F&O + NSE + Index + Options -> Derivatives + Index Options
	case f.Market == model.MarketFO && exchange == "NSE" && f.InstrumentType == model.InstrumentCategoryIndex && f.Contract == model.ContractOptions:
		dbFilter.Segment = "Derivatives"
		dbFilter.InstrumentType = "Index Options"

	// F&O + MCX + Stock + Futures -> Commodities + Commodity Futures
	case f.Market == model.MarketFO && exchange == "MCX" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractFutures:
		dbFilter.Segment = "Commodities"
		dbFilter.InstrumentType = "Commodity Futures"

	// F&O + MCX + Stock + Options -> Commodities + Commodity Options
	case f.Market == model.MarketFO && exchange == "MCX" && f.InstrumentType == model.InstrumentCategoryStock && f.Contract == model.ContractOptions:
		dbFilter.Segment = "Commodities"
		dbFilter.InstrumentType = "Commodity Options"

	default:
		// Handle partial filters or use individual field mappings
		if f.Segment != "" {
			dbFilter.Segment = f.Segment
		}
		// For other cases, we might need to add more specific logic or leave fields empty
	}

	return dbFilter
}
