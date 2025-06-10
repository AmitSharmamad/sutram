package dto

import (
	"time"
)

// GrowwPlaceOrderRequest represents the request to place an order through Groww
type GrowwPlaceOrderRequest struct {
	Symbol          string  `json:"symbol" binding:"required"`
	Exchange        string  `json:"exchange" binding:"required"`        // NSE, BSE
	Segment         string  `json:"segment" binding:"required"`         // EQT, FUT, OPT
	Side            string  `json:"side" binding:"required"`            // BUY, SELL
	OrderType       string  `json:"orderType" binding:"required"`       // LIMIT, MARKET, SL, SLM
	ProductType     string  `json:"productType" binding:"required"`     // DELIVERY, INTRADAY, MARGIN
	Validity        string  `json:"validity" binding:"required"`        // DAY, IOC, GTD
	Quantity        int     `json:"quantity" binding:"required"`
	Price           float64 `json:"price,omitempty"`
	StopPrice       float64 `json:"stopPrice,omitempty"`
	DisclosedQty    int     `json:"disclosedQty,omitempty"`
	FillType        string  `json:"fillType,omitempty"`
	ValidityDate    string  `json:"validityDate,omitempty"`
}

// GrowwPlaceOrderResponse represents the response from Groww when placing an order
type GrowwPlaceOrderResponse struct {
	OrderID       string    `json:"orderId"`
	Symbol        string    `json:"symbol"`
	Exchange      string    `json:"exchange"`
	Segment       string    `json:"segment"`
	Side          string    `json:"side"`
	OrderType     string    `json:"orderType"`
	ProductType   string    `json:"productType"`
	Validity      string    `json:"validity"`
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"price"`
	StopPrice     float64   `json:"stopPrice"`
	DisclosedQty  int       `json:"disclosedQty"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
	OrderTime     time.Time `json:"orderTime"`
}

// GrowwModifyOrderRequest represents the request to modify an existing order
type GrowwModifyOrderRequest struct {
	OrderID      string  `json:"orderId" binding:"required"`
	OrderType    string  `json:"orderType" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required"`
	Price        float64 `json:"price,omitempty"`
	StopPrice    float64 `json:"stopPrice,omitempty"`
	DisclosedQty int     `json:"disclosedQty,omitempty"`
	Validity     string  `json:"validity" binding:"required"`
}

// GrowwModifyOrderResponse represents the response from Groww when modifying an order
type GrowwModifyOrderResponse struct {
	OrderID   string    `json:"orderId"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GrowwCancelOrderRequest represents the request to cancel an order
type GrowwCancelOrderRequest struct {
	OrderID string `json:"orderId" binding:"required"`
}

// GrowwCancelOrderResponse represents the response from Groww when canceling an order
type GrowwCancelOrderResponse struct {
	OrderID   string    `json:"orderId"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GrowwOrderStatusResponse represents a single order status from Groww
type GrowwOrderStatusResponse struct {
	OrderID          string    `json:"orderId"`
	Symbol           string    `json:"symbol"`
	Exchange         string    `json:"exchange"`
	Segment          string    `json:"segment"`
	Side             string    `json:"side"`
	OrderType        string    `json:"orderType"`
	ProductType      string    `json:"productType"`
	Validity         string    `json:"validity"`
	Quantity         int       `json:"quantity"`
	Price            float64   `json:"price"`
	StopPrice        float64   `json:"stopPrice"`
	DisclosedQty     int       `json:"disclosedQty"`
	Status           string    `json:"status"`
	FilledQty        int       `json:"filledQty"`
	RemainingQty     int       `json:"remainingQty"`
	AveragePrice     float64   `json:"averagePrice"`
	OrderValue       float64   `json:"orderValue"`
	OrderTime        time.Time `json:"orderTime"`
	UpdateTime       time.Time `json:"updateTime"`
	ExchangeOrderID  string    `json:"exchangeOrderId"`
	RejectionReason  string    `json:"rejectionReason,omitempty"`
}

// GrowwOrderListResponse represents the response for getting all orders
type GrowwOrderListResponse struct {
	Orders []GrowwOrderStatusResponse `json:"orders"`
	Status string                     `json:"status"`
	Message string                    `json:"message"`
}

// GrowwOrderByIDResponse represents the response for getting a specific order by ID
type GrowwOrderByIDResponse struct {
	Order   GrowwOrderStatusResponse `json:"order"`
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
}

// GrowwTradeResponse represents a trade execution from Groww
type GrowwTradeResponse struct {
	TradeID         string    `json:"tradeId"`
	OrderID         string    `json:"orderId"`
	ExchangeOrderID string    `json:"exchangeOrderId"`
	Symbol          string    `json:"symbol"`
	Exchange        string    `json:"exchange"`
	Segment         string    `json:"segment"`
	Side            string    `json:"side"`
	Quantity        int       `json:"quantity"`
	Price           float64   `json:"price"`
	TradeValue      float64   `json:"tradeValue"`
	TradeTime       time.Time `json:"tradeTime"`
	Brokerage       float64   `json:"brokerage"`
	Tax             float64   `json:"tax"`
	NetValue        float64   `json:"netValue"`
}

// GrowwTradeListResponse represents the response for getting all trades
type GrowwTradeListResponse struct {
	Trades  []GrowwTradeResponse `json:"trades"`
	Status  string               `json:"status"`
	Message string               `json:"message"`
}