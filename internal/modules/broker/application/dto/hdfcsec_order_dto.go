package dto

// HDFCSecPlaceOrderRequest represents the request to place an order through HDFC Securities
type HDFCSecPlaceOrderRequest struct {
	ClientCode       string  `json:"clientCode" binding:"required"`
	TransactionType  string  `json:"transactionType" binding:"required"` // BUY, SELL
	Exchange         string  `json:"exchange" binding:"required"`         // NSE, BSE
	Segment          string  `json:"segment" binding:"required"`          // EQ, FO, CD, CO
	ProductType      string  `json:"productType" binding:"required"`      // CNC, MIS, NRML
	OrderType        string  `json:"orderType" binding:"required"`        // LIMIT, MARKET, SL, SLM
	Validity         string  `json:"validity" binding:"required"`         // DAY, IOC, GTC
	Symbol           string  `json:"symbol" binding:"required"`
	Quantity         int     `json:"quantity" binding:"required"`
	Price            float64 `json:"price,omitempty"`
	TriggerPrice     float64 `json:"triggerPrice,omitempty"`
	DisclosedQty     int     `json:"disclosedQty,omitempty"`
	StopLoss         float64 `json:"stopLoss,omitempty"`
	Target           float64 `json:"target,omitempty"`
	TrailingStopLoss float64 `json:"trailingStopLoss,omitempty"`
	AMO              bool    `json:"amo,omitempty"`
	Source           string  `json:"source,omitempty"`
}

// HDFCSecPlaceOrderResponse represents the response from HDFC Securities when placing an order
type HDFCSecPlaceOrderResponse struct {
	OrderID         string  `json:"orderId"`
	OrderStatus     string  `json:"orderStatus"`
	TransactionType string  `json:"transactionType"`
	Exchange        string  `json:"exchange"`
	Segment         string  `json:"segment"`
	ProductType     string  `json:"productType"`
	OrderType       string  `json:"orderType"`
	Validity        string  `json:"validity"`
	Symbol          string  `json:"symbol"`
	Quantity        int     `json:"quantity"`
	Price           float64 `json:"price"`
	TriggerPrice    float64 `json:"triggerPrice"`
	DisclosedQty    int     `json:"disclosedQty"`
	OrderTime       string  `json:"orderTime"`
	UpdateTime      string  `json:"updateTime"`
	ExchangeOrderID string  `json:"exchangeOrderId"`
	Message         string  `json:"message"`
}

// HDFCSecModifyOrderRequest represents the request to modify an existing order
type HDFCSecModifyOrderRequest struct {
	ClientCode   string  `json:"clientCode" binding:"required"`
	OrderID      string  `json:"orderId" binding:"required"`
	OrderType    string  `json:"orderType" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required"`
	Price        float64 `json:"price,omitempty"`
	TriggerPrice float64 `json:"triggerPrice,omitempty"`
	DisclosedQty int     `json:"disclosedQty,omitempty"`
	Validity     string  `json:"validity" binding:"required"`
}

// HDFCSecModifyOrderResponse represents the response from HDFC Securities when modifying an order
type HDFCSecModifyOrderResponse struct {
	OrderID      string  `json:"orderId"`
	OrderStatus  string  `json:"orderStatus"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
	TriggerPrice float64 `json:"triggerPrice"`
	Validity     string  `json:"validity"`
	UpdateTime   string  `json:"updateTime"`
	Message      string  `json:"message"`
}

// HDFCSecCancelOrderRequest represents the request to cancel an order
type HDFCSecCancelOrderRequest struct {
	ClientCode string `json:"clientCode" binding:"required"`
	OrderID    string `json:"orderId" binding:"required"`
}

// HDFCSecCancelOrderResponse represents the response from HDFC Securities when canceling an order
type HDFCSecCancelOrderResponse struct {
	OrderID     string `json:"orderId"`
	OrderStatus string `json:"orderStatus"`
	UpdateTime  string `json:"updateTime"`
	Message     string `json:"message"`
}

// HDFCSecOrderStatusResponse represents a single order status from HDFC Securities
type HDFCSecOrderStatusResponse struct {
	OrderID           string  `json:"orderId"`
	OrderStatus       string  `json:"orderStatus"`
	TransactionType   string  `json:"transactionType"`
	Exchange          string  `json:"exchange"`
	Segment           string  `json:"segment"`
	ProductType       string  `json:"productType"`
	OrderType         string  `json:"orderType"`
	Validity          string  `json:"validity"`
	Symbol            string  `json:"symbol"`
	Quantity          int     `json:"quantity"`
	PendingQuantity   int     `json:"pendingQuantity"`
	FilledQuantity    int     `json:"filledQuantity"`
	CancelledQuantity int     `json:"cancelledQuantity"`
	Price             float64 `json:"price"`
	TriggerPrice      float64 `json:"triggerPrice"`
	DisclosedQty      int     `json:"disclosedQty"`
	AveragePrice      float64 `json:"averagePrice"`
	OrderValue        float64 `json:"orderValue"`
	OrderTime         string  `json:"orderTime"`
	UpdateTime        string  `json:"updateTime"`
	ExchangeTime      string  `json:"exchangeTime"`
	ExchangeOrderID   string  `json:"exchangeOrderId"`
	RejectionReason   string  `json:"rejectionReason"`
	Source            string  `json:"source"`
}

// HDFCSecOrderListResponse represents the response for getting all orders
type HDFCSecOrderListResponse struct {
	Orders []HDFCSecOrderStatusResponse `json:"data"`
}

// HDFCSecOrderByIDResponse represents the response for getting a specific order by ID
type HDFCSecOrderByIDResponse struct {
	Order HDFCSecOrderStatusResponse `json:"data"`
}

// HDFCSecTradeResponse represents a trade execution from HDFC Securities
type HDFCSecTradeResponse struct {
	ClientCode      string  `json:"clientCode"`
	OrderID         string  `json:"orderId"`
	ExchangeOrderID string  `json:"exchangeOrderId"`
	TradeID         string  `json:"tradeId"`
	TransactionType string  `json:"transactionType"`
	Exchange        string  `json:"exchange"`
	Segment         string  `json:"segment"`
	ProductType     string  `json:"productType"`
	Symbol          string  `json:"symbol"`
	TradeQuantity   int     `json:"tradeQuantity"`
	TradePrice      float64 `json:"tradePrice"`
	TradeValue      float64 `json:"tradeValue"`
	TradeTime       string  `json:"tradeTime"`
	BrokerageAmount float64 `json:"brokerageAmount"`
	ExchangeFee     float64 `json:"exchangeFee"`
	STT             float64 `json:"stt"`
	NetAmount       float64 `json:"netAmount"`
}

// HDFCSecTradeListResponse represents the response for getting all trades
type HDFCSecTradeListResponse struct {
	Trades []HDFCSecTradeResponse `json:"data"`
}

// HDFCSecPositionResponse represents a position from HDFC Securities
type HDFCSecPositionResponse struct {
	ClientCode       string  `json:"clientCode"`
	Exchange         string  `json:"exchange"`
	Segment          string  `json:"segment"`
	ProductType      string  `json:"productType"`
	Symbol           string  `json:"symbol"`
	NetQuantity      int     `json:"netQuantity"`
	BuyQuantity      int     `json:"buyQuantity"`
	SellQuantity     int     `json:"sellQuantity"`
	AverageBuyPrice  float64 `json:"averageBuyPrice"`
	AverageSellPrice float64 `json:"averageSellPrice"`
	NetAmount        float64 `json:"netAmount"`
	RealizedPL       float64 `json:"realizedPL"`
	UnrealizedPL     float64 `json:"unrealizedPL"`
	MTM              float64 `json:"mtm"`
	LTP              float64 `json:"ltp"`
	LastTradeTime    string  `json:"lastTradeTime"`
}

// HDFCSecPositionListResponse represents the response for getting all positions
type HDFCSecPositionListResponse struct {
	Positions []HDFCSecPositionResponse `json:"data"`
}