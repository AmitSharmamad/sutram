package dto

import (
	"time"

	"github.com/google/uuid"
)

// SharekhanPlaceOrderRequest represents the request to place an order through Sharekhan
type SharekhanPlaceOrderRequest struct {
	TradingSymbol    string  `json:"tradingSymbol" binding:"required"`
	Exchange         string  `json:"exchange" binding:"required"`         // NSE, BSE, NFO, etc.
	TransactionType  string  `json:"transactionType" binding:"required"`  // BUY, SELL
	OrderType        string  `json:"orderType" binding:"required"`        // MARKET, LIMIT, SL, SL-M
	ProductCode      string  `json:"productCode" binding:"required"`      // CNC, MIS, NRML
	Quantity         int     `json:"quantity" binding:"required"`
	Price            float64 `json:"price,omitempty"`
	TriggerPrice     float64 `json:"triggerPrice,omitempty"`
	Validity         string  `json:"validity" binding:"required"`         // DAY, IOC, GTC
	DisclosedQty     int     `json:"disclosedQty,omitempty"`
	Remarks          string  `json:"remarks,omitempty"`
}

// SharekhanPlaceOrderResponse represents the response from Sharekhan when placing an order
type SharekhanPlaceOrderResponse struct {
	OrderNumber      string  `json:"orderNumber"`
	OrderType        string  `json:"orderType"`
	TradingSymbol    string  `json:"tradingSymbol"`
	Exchange         string  `json:"exchange"`
	TransactionType  string  `json:"transactionType"`
	ProductCode      string  `json:"productCode"`
	Quantity         int     `json:"quantity"`
	Price            float64 `json:"price"`
	TriggerPrice     float64 `json:"triggerPrice"`
	Validity         string  `json:"validity"`
	OrderStatus      string  `json:"orderStatus"`
	OrderTime        string  `json:"orderTime"`
	Message          string  `json:"message"`
	RejectionReason  string  `json:"rejectionReason,omitempty"`
}

// SharekhanModifyOrderRequest represents the request to modify an existing order
type SharekhanModifyOrderRequest struct {
	OrderNumber      string  `json:"orderNumber" binding:"required"`
	OrderType        string  `json:"orderType" binding:"required"`
	Quantity         int     `json:"quantity" binding:"required"`
	Price            float64 `json:"price,omitempty"`
	TriggerPrice     float64 `json:"triggerPrice,omitempty"`
	Validity         string  `json:"validity" binding:"required"`
	DisclosedQty     int     `json:"disclosedQty,omitempty"`
}

// SharekhanModifyOrderResponse represents the response from Sharekhan when modifying an order
type SharekhanModifyOrderResponse struct {
	OrderNumber      string  `json:"orderNumber"`
	OrderType        string  `json:"orderType"`
	TradingSymbol    string  `json:"tradingSymbol"`
	Exchange         string  `json:"exchange"`
	TransactionType  string  `json:"transactionType"`
	ProductCode      string  `json:"productCode"`
	Quantity         int     `json:"quantity"`
	Price            float64 `json:"price"`
	TriggerPrice     float64 `json:"triggerPrice"`
	Validity         string  `json:"validity"`
	OrderStatus      string  `json:"orderStatus"`
	ModifiedTime     string  `json:"modifiedTime"`
	Message          string  `json:"message"`
}

// SharekhanCancelOrderRequest represents the request to cancel an order
type SharekhanCancelOrderRequest struct {
	OrderNumber string `json:"orderNumber" binding:"required"`
}

// SharekhanCancelOrderResponse represents the response from Sharekhan when cancelling an order
type SharekhanCancelOrderResponse struct {
	OrderNumber      string `json:"orderNumber"`
	OrderStatus      string `json:"orderStatus"`
	CancelledTime    string `json:"cancelledTime"`
	Message          string `json:"message"`
}

// SharekhanOrderResponse represents a single order in the response
type SharekhanOrderResponse struct {
	OrderNumber       string    `json:"orderNumber"`
	ExchangeOrderID   string    `json:"exchangeOrderId"`
	TradingSymbol     string    `json:"tradingSymbol"`
	Exchange          string    `json:"exchange"`
	TransactionType   string    `json:"transactionType"`
	OrderType         string    `json:"orderType"`
	ProductCode       string    `json:"productCode"`
	Quantity          int       `json:"quantity"`
	FilledQuantity    int       `json:"filledQuantity"`
	RemainingQuantity int       `json:"remainingQuantity"`
	Price             float64   `json:"price"`
	TriggerPrice      float64   `json:"triggerPrice"`
	AveragePrice      float64   `json:"averagePrice"`
	OrderStatus       string    `json:"orderStatus"`
	Validity          string    `json:"validity"`
	OrderTime         time.Time `json:"orderTime"`
	UpdateTime        time.Time `json:"updateTime"`
	RejectionReason   string    `json:"rejectionReason,omitempty"`
	DisclosedQty      int       `json:"disclosedQty"`
	Remarks           string    `json:"remarks,omitempty"`
}

// SharekhanOrderListResponse represents the response for getting all orders
type SharekhanOrderListResponse struct {
	Orders []SharekhanOrderResponse `json:"orders"`
	Status string                   `json:"status"`
	Message string                  `json:"message"`
}

// SharekhanOrderByIDResponse represents the response for getting a specific order
type SharekhanOrderByIDResponse struct {
	Order   SharekhanOrderResponse `json:"order"`
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
}

// SharekhanTradeResponse represents a single trade in the response
type SharekhanTradeResponse struct {
	TradeNumber      string    `json:"tradeNumber"`
	OrderNumber      string    `json:"orderNumber"`
	ExchangeOrderID  string    `json:"exchangeOrderId"`
	TradingSymbol    string    `json:"tradingSymbol"`
	Exchange         string    `json:"exchange"`
	TransactionType  string    `json:"transactionType"`
	ProductCode      string    `json:"productCode"`
	Quantity         int       `json:"quantity"`
	Price            float64   `json:"price"`
	TradeValue       float64   `json:"tradeValue"`
	TradeTime        time.Time `json:"tradeTime"`
	Brokerage        float64   `json:"brokerage"`
	NetValue         float64   `json:"netValue"`
}

// SharekhanTradeListResponse represents the response for getting all trades
type SharekhanTradeListResponse struct {
	Trades  []SharekhanTradeResponse `json:"trades"`
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
}

// SharekhanOrderBookResponse represents the response for order book
type SharekhanOrderBookResponse struct {
	OrderBook []SharekhanOrderResponse `json:"orderBook"`
	Status    string                   `json:"status"`
	Message   string                   `json:"message"`
}

// SharekhanTradeBookResponse represents the response for trade book
type SharekhanTradeBookResponse struct {
	TradeBook []SharekhanTradeResponse `json:"tradeBook"`
	Status    string                   `json:"status"`
	Message   string                   `json:"message"`
}

// Internal DTOs for our application - these match the existing pattern

// SharekhanPlaceOrderDTO represents the internal DTO for placing orders through Sharekhan
type SharekhanPlaceOrderDTO struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id" binding:"required"`
	PositionID       uuid.UUID `json:"position_id" binding:"required"`
	TradingSymbol    string    `json:"trading_symbol" binding:"required"`
	Exchange         string    `json:"exchange" binding:"required"`
	TransactionType  string    `json:"transaction_type" binding:"required"`
	OrderType        string    `json:"order_type" binding:"required"`
	ProductCode      string    `json:"product_code" binding:"required"`
	Quantity         int       `json:"quantity" binding:"required"`
	Price            float64   `json:"price"`
	TriggerPrice     float64   `json:"trigger_price"`
	Validity         string    `json:"validity" binding:"required"`
	DisclosedQty     int       `json:"disclosed_qty"`
	Remarks          string    `json:"remarks"`
	AlgoID           string    `json:"algo_id,omitempty"`
}

// SharekhanOrderResponseDTO represents the response when placing/modifying orders through Sharekhan
type SharekhanOrderResponseDTO struct {
	ID                      uuid.UUID `json:"id"`
	UserID                  uuid.UUID `json:"user_id"`
	PositionID              uuid.UUID `json:"position_id"`
	TradingSymbol           string    `json:"trading_symbol"`
	Exchange                string    `json:"exchange"`
	TransactionType         string    `json:"transaction_type"`
	OrderType               string    `json:"order_type"`
	ProductCode             string    `json:"product_code"`
	Quantity                int       `json:"quantity"`
	Price                   float64   `json:"price"`
	TriggerPrice            float64   `json:"trigger_price"`
	Validity                string    `json:"validity"`
	State                   string    `json:"state"`
	BrokerOrderID           string    `json:"broker_order_id,omitempty"`
	BrokerOrderResponse     any       `json:"broker_order_response,omitempty"`
	AlgoID                  string    `json:"algo_id,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}