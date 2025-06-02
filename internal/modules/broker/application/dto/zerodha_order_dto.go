package dto

import (
	"time"
)

// ZerodhaPlaceOrderRequest represents the request to place an order through Zerodha Kite
type ZerodhaPlaceOrderRequest struct {
	TradingSymbol   string  `json:"tradingsymbol" binding:"required"`
	Exchange        string  `json:"exchange" binding:"required"`       // NSE, BSE, NFO, etc.
	TransactionType string  `json:"transaction_type" binding:"required"` // BUY, SELL
	OrderType       string  `json:"order_type" binding:"required"`     // MARKET, LIMIT, SL, SL-M
	Quantity        int     `json:"quantity" binding:"required"`
	Price           float64 `json:"price,omitempty"`
	Product         string  `json:"product" binding:"required"`        // CNC, MIS, NRML
	Validity        string  `json:"validity" binding:"required"`       // DAY, IOC
	DisclosedQuantity int   `json:"disclosed_quantity,omitempty"`
	TriggerPrice    float64 `json:"trigger_price,omitempty"`
	SquareOff       float64 `json:"squareoff,omitempty"`
	StopLoss        float64 `json:"stoploss,omitempty"`
	TrailingStopLoss float64 `json:"trailing_stoploss,omitempty"`
	Variety         string  `json:"variety,omitempty"`                 // regular, bo, co, amo
	Tag             string  `json:"tag,omitempty"`
}

// ZerodhaPlaceOrderResponse represents the response from Zerodha when placing an order
type ZerodhaPlaceOrderResponse struct {
	OrderID string `json:"order_id"`
}

// ZerodhaModifyOrderRequest represents the request to modify an existing order
type ZerodhaModifyOrderRequest struct {
	OrderID           string  `json:"order_id" binding:"required"`
	Quantity          int     `json:"quantity,omitempty"`
	Price             float64 `json:"price,omitempty"`
	OrderType         string  `json:"order_type,omitempty"`
	TriggerPrice      float64 `json:"trigger_price,omitempty"`
	Validity          string  `json:"validity,omitempty"`
	DisclosedQuantity int     `json:"disclosed_quantity,omitempty"`
	ParentOrderID     string  `json:"parent_order_id,omitempty"`
}

// ZerodhaModifyOrderResponse represents the response from Zerodha when modifying an order
type ZerodhaModifyOrderResponse struct {
	OrderID string `json:"order_id"`
}

// ZerodhaCancelOrderRequest represents the request to cancel an order
type ZerodhaCancelOrderRequest struct {
	OrderID       string `json:"order_id" binding:"required"`
	Variety       string `json:"variety,omitempty"`
	ParentOrderID string `json:"parent_order_id,omitempty"`
}

// ZerodhaCancelOrderResponse represents the response from Zerodha when canceling an order
type ZerodhaCancelOrderResponse struct {
	OrderID string `json:"order_id"`
}

// ZerodhaOrderResponse represents a single order from Zerodha
type ZerodhaOrderResponse struct {
	OrderID             string    `json:"order_id"`
	ExchangeOrderID     string    `json:"exchange_order_id"`
	ParentOrderID       string    `json:"parent_order_id"`
	Status              string    `json:"status"`
	StatusMessage       string    `json:"status_message"`
	StatusMessageRaw    string    `json:"status_message_raw"`
	OrderTimestamp      time.Time `json:"order_timestamp"`
	ExchangeTimestamp   time.Time `json:"exchange_timestamp"`
	Variety             string    `json:"variety"`
	Exchange            string    `json:"exchange"`
	TradingSymbol       string    `json:"tradingsymbol"`
	InstrumentToken     int       `json:"instrument_token"`
	TransactionType     string    `json:"transaction_type"`
	OrderType           string    `json:"order_type"`
	Product             string    `json:"product"`
	Validity            string    `json:"validity"`
	Price               float64   `json:"price"`
	Quantity            int       `json:"quantity"`
	TriggerPrice        float64   `json:"trigger_price"`
	AveragePrice        float64   `json:"average_price"`
	PendingQuantity     int       `json:"pending_quantity"`
	FilledQuantity      int       `json:"filled_quantity"`
	DisclosedQuantity   int       `json:"disclosed_quantity"`
	Market              bool      `json:"market"`
	Tag                 string    `json:"tag"`
	GUID                string    `json:"guid"`
}

// ZerodhaOrderListResponse represents the response for getting all orders
type ZerodhaOrderListResponse struct {
	Status string                 `json:"status"`
	Data   []ZerodhaOrderResponse `json:"data"`
}

// ZerodhaOrderHistoryResponse represents the response for getting order history
type ZerodhaOrderHistoryResponse struct {
	Status string                 `json:"status"`
	Data   []ZerodhaOrderResponse `json:"data"`
}

// ZerodhaTradeResponse represents a trade execution from Zerodha
type ZerodhaTradeResponse struct {
	TradeID           string    `json:"trade_id"`
	OrderID           string    `json:"order_id"`
	ExchangeOrderID   string    `json:"exchange_order_id"`
	TradingSymbol     string    `json:"tradingsymbol"`
	Exchange          string    `json:"exchange"`
	InstrumentToken   int       `json:"instrument_token"`
	TransactionType   string    `json:"transaction_type"`
	Product           string    `json:"product"`
	AveragePrice      float64   `json:"average_price"`
	Quantity          int       `json:"quantity"`
	FillTimestamp     time.Time `json:"fill_timestamp"`
	ExchangeTimestamp time.Time `json:"exchange_timestamp"`
}

// ZerodhaTradeListResponse represents the response for getting all trades
type ZerodhaTradeListResponse struct {
	Status string                 `json:"status"`
	Data   []ZerodhaTradeResponse `json:"data"`
}

// ZerodhaPositionResponse represents a position from Zerodha
type ZerodhaPositionResponse struct {
	TradingSymbol     string  `json:"tradingsymbol"`
	Exchange          string  `json:"exchange"`
	InstrumentToken   int     `json:"instrument_token"`
	Product           string  `json:"product"`
	Quantity          int     `json:"quantity"`
	OvernightQuantity int     `json:"overnight_quantity"`
	Multiplier        float64 `json:"multiplier"`
	AveragePrice      float64 `json:"average_price"`
	ClosePrice        float64 `json:"close_price"`
	LastPrice         float64 `json:"last_price"`
	Value             float64 `json:"value"`
	PNL               float64 `json:"pnl"`
	M2M               float64 `json:"m2m"`
	Unrealised        float64 `json:"unrealised"`
	Realised          float64 `json:"realised"`
	BuyQuantity       int     `json:"buy_quantity"`
	BuyPrice          float64 `json:"buy_price"`
	BuyValue          float64 `json:"buy_value"`
	BuyM2M            float64 `json:"buy_m2m"`
	SellQuantity      int     `json:"sell_quantity"`
	SellPrice         float64 `json:"sell_price"`
	SellValue         float64 `json:"sell_value"`
	SellM2M           float64 `json:"sell_m2m"`
	DayBuyQuantity    int     `json:"day_buy_quantity"`
	DayBuyPrice       float64 `json:"day_buy_price"`
	DayBuyValue       float64 `json:"day_buy_value"`
	DaySellQuantity   int     `json:"day_sell_quantity"`
	DaySellPrice      float64 `json:"day_sell_price"`
	DaySellValue      float64 `json:"day_sell_value"`
}

// ZerodhaPositionListResponse represents the response for getting positions
type ZerodhaPositionListResponse struct {
	Status string `json:"status"`
	Data   struct {
		Net []ZerodhaPositionResponse `json:"net"`
		Day []ZerodhaPositionResponse `json:"day"`
	} `json:"data"`
}