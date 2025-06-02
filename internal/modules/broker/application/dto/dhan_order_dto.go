package dto

import (
	"time"

	"github.com/google/uuid"
)

// DhanPlaceOrderRequest represents the request to place an order through Dhan
type DhanPlaceOrderRequest struct {
	DhanClientID       string  `json:"dhanClientId" binding:"required"`
	TransactionType    string  `json:"transactionType" binding:"required"` // BUY, SELL
	ExchangeSegment    string  `json:"exchangeSegment" binding:"required"` // NSE_EQ, NSE_FNO, BSE_EQ, etc.
	ProductType        string  `json:"productType" binding:"required"`     // CNC, MIS, NRML
	OrderType          string  `json:"orderType" binding:"required"`       // LIMIT, MARKET, SL, SL-M
	Validity           string  `json:"validity" binding:"required"`        // DAY, IOC
	SecurityID         string  `json:"securityId" binding:"required"`
	Quantity           int     `json:"quantity" binding:"required"`
	DisclosedQuantity  int     `json:"disclosedQuantity,omitempty"`
	Price              float64 `json:"price,omitempty"`
	TriggerPrice       float64 `json:"triggerPrice,omitempty"`
	AfterMarketOrder   bool    `json:"afterMarketOrder,omitempty"`
	AmoTime            string  `json:"amoTime,omitempty"`
	BolBlockDeal       bool    `json:"bolBlockDeal,omitempty"`
	CorrelationID      string  `json:"correlationId,omitempty"`
}

// DhanPlaceOrderResponse represents the response from Dhan when placing an order
type DhanPlaceOrderResponse struct {
	OrderID      string `json:"orderId"`
	DrmOrderID   string `json:"drmOrderId"`
	OrderStatus  string `json:"orderStatus"`
	TransactionType string `json:"transactionType"`
	ExchangeSegment string `json:"exchangeSegment"`
	ProductType  string `json:"productType"`
	OrderType    string `json:"orderType"`
	Validity     string `json:"validity"`
	TradingSymbol string `json:"tradingSymbol"`
	SecurityID   string `json:"securityId"`
	Quantity     int    `json:"quantity"`
	DisclosedQuantity int `json:"disclosedQuantity"`
	Price        float64 `json:"price"`
	TriggerPrice float64 `json:"triggerPrice"`
	AfterMarketOrder bool `json:"afterMarketOrder"`
	BolBlockDeal bool    `json:"bolBlockDeal"`
	CreateTime   string  `json:"createTime"`
	UpdateTime   string  `json:"updateTime"`
	ExchangeTime string  `json:"exchangeTime"`
	DrvExOrderID string  `json:"drvExOrderId"`
	AlgoID       string  `json:"algoId"`
}

// DhanModifyOrderRequest represents the request to modify an existing order
type DhanModifyOrderRequest struct {
	DhanClientID      string  `json:"dhanClientId" binding:"required"`
	OrderID           string  `json:"orderId" binding:"required"`
	OrderType         string  `json:"orderType" binding:"required"`
	LegName           string  `json:"legName,omitempty"`
	Quantity          int     `json:"quantity" binding:"required"`
	Price             float64 `json:"price,omitempty"`
	TriggerPrice      float64 `json:"triggerPrice,omitempty"`
	DisclosedQuantity int     `json:"disclosedQuantity,omitempty"`
	Validity          string  `json:"validity" binding:"required"`
}

// DhanModifyOrderResponse represents the response from Dhan when modifying an order
type DhanModifyOrderResponse struct {
	OrderID      string  `json:"orderId"`
	OrderStatus  string  `json:"orderStatus"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
	TriggerPrice float64 `json:"triggerPrice"`
	Validity     string  `json:"validity"`
	UpdateTime   string  `json:"updateTime"`
}

// DhanCancelOrderRequest represents the request to cancel an order
type DhanCancelOrderRequest struct {
	DhanClientID string `json:"dhanClientId" binding:"required"`
	OrderID      string `json:"orderId" binding:"required"`
}

// DhanCancelOrderResponse represents the response from Dhan when canceling an order
type DhanCancelOrderResponse struct {
	OrderID     string `json:"orderId"`
	OrderStatus string `json:"orderStatus"`
	UpdateTime  string `json:"updateTime"`
}

// DhanOrderStatusResponse represents a single order status from Dhan
type DhanOrderStatusResponse struct {
	OrderID           string  `json:"orderId"`
	DrmOrderID        string  `json:"drmOrderId"`
	OrderStatus       string  `json:"orderStatus"`
	TransactionType   string  `json:"transactionType"`
	ExchangeSegment   string  `json:"exchangeSegment"`
	ProductType       string  `json:"productType"`
	OrderType         string  `json:"orderType"`
	Validity          string  `json:"validity"`
	TradingSymbol     string  `json:"tradingSymbol"`
	SecurityID        string  `json:"securityId"`
	Quantity          int     `json:"quantity"`
	DisclosedQuantity int     `json:"disclosedQuantity"`
	Price             float64 `json:"price"`
	TriggerPrice      float64 `json:"triggerPrice"`
	AfterMarketOrder  bool    `json:"afterMarketOrder"`
	BolBlockDeal      bool    `json:"bolBlockDeal"`
	LegName           string  `json:"legName"`
	CreateTime        string  `json:"createTime"`
	UpdateTime        string  `json:"updateTime"`
	ExchangeTime      string  `json:"exchangeTime"`
	DrvExOrderID      string  `json:"drvExOrderId"`
	AlgoID            string  `json:"algoId"`
	FilledQty         int     `json:"filledQty"`
	RemainingQuantity int     `json:"remainingQuantity"`
	AveragePrice      float64 `json:"averagePrice"`
	OrderValue        float64 `json:"orderValue"`
	BrokerageAmount   float64 `json:"brokerageAmount"`
}

// DhanOrderListResponse represents the response for getting all orders
type DhanOrderListResponse struct {
	Orders []DhanOrderStatusResponse `json:"data"`
}

// DhanOrderByIDResponse represents the response for getting a specific order by ID
type DhanOrderByIDResponse struct {
	Order DhanOrderStatusResponse `json:"data"`
}

// DhanTradeResponse represents a trade execution from Dhan
type DhanTradeResponse struct {
	DhanClientID    string  `json:"dhanClientId"`
	OrderID         string  `json:"orderId"`
	ExchangeOrderID string  `json:"exchangeOrderId"`
	TradeID         string  `json:"tradeId"`
	TransactionType string  `json:"transactionType"`
	ExchangeSegment string  `json:"exchangeSegment"`
	ProductType     string  `json:"productType"`
	TradingSymbol   string  `json:"tradingSymbol"`
	SecurityID      string  `json:"securityId"`
	TradeQuantity   int     `json:"tradeQuantity"`
	TradePrice      float64 `json:"tradePrice"`
	TradeValue      float64 `json:"tradeValue"`
	TradeTime       string  `json:"tradeTime"`
	BrokerageAmount float64 `json:"brokerageAmount"`
}

// DhanTradeListResponse represents the response for getting all trades
type DhanTradeListResponse struct {
	Trades []DhanTradeResponse `json:"data"`
}

// Internal DTOs for our application

// PlaceOrderDTO represents the internal DTO for placing orders
type PlaceOrderDTO struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id" binding:"required"`
	PositionID      uuid.UUID `json:"position_id" binding:"required"`
	Exchange        string    `json:"exchange" binding:"required"`
	Segment         string    `json:"segment" binding:"required"`
	Instrument      string    `json:"instrument" binding:"required"`
	TransactionType string    `json:"transaction_type" binding:"required"`
	Quantity        int       `json:"quantity" binding:"required"`
	Price           float64   `json:"price" binding:"required"`
	ProductType     string    `json:"product_type" binding:"required"`
	OrderType       string    `json:"order_type" binding:"required"`
	AlgoID          string    `json:"algo_id,omitempty"`
}

// OrderResponseDTO represents the response when placing/modifying orders
type OrderResponseDTO struct {
	ID                   uuid.UUID `json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	PositionID           uuid.UUID `json:"position_id"`
	Exchange             string    `json:"exchange"`
	Segment              string    `json:"segment"`
	Instrument           string    `json:"instrument"`
	TransactionType      string    `json:"transaction_type"`
	Quantity             int       `json:"quantity"`
	Price                float64   `json:"price"`
	State                string    `json:"state"`
	ProductType          string    `json:"product_type"`
	OrderType            string    `json:"order_type"`
	BrokerOrderID        string    `json:"broker_order_id,omitempty"`
	BrokerOrderResponse  any       `json:"broker_order_response,omitempty"`
	AlgoID               string    `json:"algo_id,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	ExecutedAt           time.Time `json:"executed_at,omitempty"`
}