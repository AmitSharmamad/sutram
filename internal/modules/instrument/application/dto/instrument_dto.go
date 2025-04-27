package dto

import (
	"time"
)

// InstrumentDTO represents instrument data for API responses
type InstrumentDTO struct {
	ID        uint64    `json:"id"`
	Symbol    string    `json:"symbol"`
	Open      float32   `json:"open"`
	High      float32   `json:"high"`
	Low       float32   `json:"low"`
	Close     float32   `json:"close"`
	Volume    float32   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

// InstrumentMasterDTO represents instrument master data for API responses
type InstrumentMasterDTO struct {
	ID        uint64    `json:"id"`
	Exchange  string    `json:"exchange"`
	Segment   string    `json:"segment"`
	Symbol    string    `json:"symbol"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	Expiry    time.Time `json:"expiry,omitempty"`
}

// InstrumentWithMetadataDTO combines instrument data with its metadata
type InstrumentWithMetadataDTO struct {
	Instrument       *InstrumentDTO       `json:"instrument"`
	InstrumentMaster *InstrumentMasterDTO `json:"instrument_master"`
}

// CreateInstrumentRequest represents request to create an instrument
type CreateInstrumentRequest struct {
	Symbol    string    `json:"symbol" binding:"required"`
	Open      float32   `json:"open" binding:"required"`
	High      float32   `json:"high" binding:"required"`
	Low       float32   `json:"low" binding:"required"`
	Close     float32   `json:"close" binding:"required"`
	Volume    float32   `json:"volume" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

// UpdateInstrumentRequest represents request to update an instrument
type UpdateInstrumentRequest struct {
	Open      float32   `json:"open" binding:"required"`
	High      float32   `json:"high" binding:"required"`
	Low       float32   `json:"low" binding:"required"`
	Close     float32   `json:"close" binding:"required"`
	Volume    float32   `json:"volume" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

// CreateInstrumentMasterRequest represents request to create an instrument master
type CreateInstrumentMasterRequest struct {
	Exchange  string    `json:"exchange" binding:"required"`
	Segment   string    `json:"segment" binding:"required"`
	Symbol    string    `json:"symbol" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	Expiry    time.Time `json:"expiry,omitempty"`
}

// UpdateInstrumentMasterRequest represents request to update an instrument master
type UpdateInstrumentMasterRequest struct {
	Exchange  string    `json:"exchange"`
	Segment   string    `json:"segment"`
	Symbol    string    `json:"symbol"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	Expiry    time.Time `json:"expiry,omitempty"`
}

// GetInstrumentsByTimeRangeRequest represents a request to fetch instruments by time range
type GetInstrumentsByTimeRangeRequest struct {
	Symbol    string    `json:"symbol" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

// InstrumentMasterFilterDTO represents filtering and pagination options for instrument masters
type InstrumentMasterFilterDTO struct {
	Exchange string `form:"exchange"`
	Segment  string `form:"segment"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
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
