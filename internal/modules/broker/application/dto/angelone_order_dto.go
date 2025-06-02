package dto

// Note: PlaceOrderDTO and OrderResponseDTO are defined in dhan_order_dto.go for consistency across brokers

// AngelOnePlaceOrderRequest represents the request to place an order through AngelOne
type AngelOnePlaceOrderRequest struct {
	Variety           string  `json:"variety" binding:"required"`           // regular, bo, co, iceberg, auction
	TradingSymbol     string  `json:"tradingsymbol" binding:"required"`
	SymbolToken       string  `json:"symboltoken" binding:"required"`
	TransactionType   string  `json:"transactiontype" binding:"required"`   // BUY, SELL
	Exchange          string  `json:"exchange" binding:"required"`          // NSE, BSE, NFO, etc.
	OrderType         string  `json:"ordertype" binding:"required"`         // MARKET, LIMIT, STOPLOSS_LIMIT, STOPLOSS_MARKET
	ProductType       string  `json:"producttype" binding:"required"`       // DELIVERY, CARRYFORWARD, MARGIN, INTRADAY, BO
	Duration          string  `json:"duration" binding:"required"`          // DAY, IOC
	Price             string  `json:"price,omitempty"`
	SquareOff         string  `json:"squareoff,omitempty"`
	StopLoss          string  `json:"stoploss,omitempty"`
	Quantity          string  `json:"quantity" binding:"required"`
	DisclosedQuantity string  `json:"disclosedquantity,omitempty"`
	TriggerPrice      string  `json:"triggerprice,omitempty"`
	Tag               string  `json:"tag,omitempty"`
}

// AngelOnePlaceOrderResponse represents the response from AngelOne when placing an order
type AngelOnePlaceOrderResponse struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode,omitempty"`
	Data      struct {
		Script   string `json:"script"`
		OrderID  string `json:"orderid"`
		UniqueID string `json:"uniqueid"`
	} `json:"data"`
}

// AngelOneModifyOrderRequest represents the request to modify an order
type AngelOneModifyOrderRequest struct {
	Variety           string `json:"variety" binding:"required"`
	OrderID           string `json:"orderid" binding:"required"`
	OrderType         string `json:"ordertype" binding:"required"`
	ProductType       string `json:"producttype" binding:"required"`
	Duration          string `json:"duration" binding:"required"`
	Price             string `json:"price,omitempty"`
	Quantity          string `json:"quantity" binding:"required"`
	DisclosedQuantity string `json:"disclosedquantity,omitempty"`
	TriggerPrice      string `json:"triggerprice,omitempty"`
}

// AngelOneModifyOrderResponse represents the response from AngelOne when modifying an order
type AngelOneModifyOrderResponse struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode,omitempty"`
	Data      struct {
		OrderID  string `json:"orderid"`
		UniqueID string `json:"uniqueid"`
	} `json:"data"`
}

// AngelOneCancelOrderRequest represents the request to cancel an order
type AngelOneCancelOrderRequest struct {
	Variety string `json:"variety" binding:"required"`
	OrderID string `json:"orderid" binding:"required"`
}

// AngelOneCancelOrderResponse represents the response from AngelOne when canceling an order
type AngelOneCancelOrderResponse struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode,omitempty"`
	Data      struct {
		OrderID  string `json:"orderid"`
		UniqueID string `json:"uniqueid"`
	} `json:"data"`
}

// AngelOneOrderResponse represents a single order in the order book
type AngelOneOrderResponse struct {
	Variety            string `json:"variety"`
	OrderType          string `json:"ordertype"`
	ProductType        string `json:"producttype"`
	Duration           string `json:"duration"`
	Price              string `json:"price"`
	TriggerPrice       string `json:"triggerprice"`
	Quantity           string `json:"quantity"`
	DisclosedQuantity  string `json:"disclosedquantity"`
	PendingQuantity    string `json:"pendingquantity"`
	FilledQuantity     string `json:"filledquantity"`
	UnfilledQuantity   string `json:"unfilledquantity"`
	OrderID            string `json:"orderid"`
	Text               string `json:"text"`
	Status             string `json:"status"`
	OrderTag           string `json:"ordertag"`
	Instrument         string `json:"instrument"`
	Exchange           string `json:"exchange"`
	TradingSymbol      string `json:"tradingsymbol"`
	SymbolToken        string `json:"symboltoken"`
	TransactionType    string `json:"transactiontype"`
	FillID             string `json:"fillid"`
	FillTime           string `json:"filltime"`
	FillPrice          string `json:"fillprice"`
	FillSize           string `json:"fillsize"`
	CancelledQuantity  string `json:"cancelledquantity"`
	CancelledTime      string `json:"cancelledtime"`
	SquareOff          string `json:"squareoff"`
	StopLoss           string `json:"stoploss"`
	TrailingStopLoss   string `json:"trailingstoploss"`
	LTP                string `json:"ltp"`
	AveragePrice       string `json:"averageprice"`
	LotSize            string `json:"lotsize"`
	CancelSize         string `json:"cancelsize"`
	UpdateTime         string `json:"updatetime"`
	ExchangeUpdateTime string `json:"exchuptime"`
	ExchangeOrderID    string `json:"exchorderid"`
	ExchangeOrderUpdateTime string `json:"exchorderupdatetime"`
	FillSharesCount    string `json:"fillshares"`
	UniqueID           string `json:"uniqueid"`
	OrderPlacedBy      string `json:"orderplacedby"`
}

// AngelOneOrderListResponse represents the response for getting all orders
type AngelOneOrderListResponse struct {
	Status    bool                    `json:"status"`
	Message   string                  `json:"message"`
	ErrorCode string                  `json:"errorcode,omitempty"`
	Data      []AngelOneOrderResponse `json:"data"`
}

// AngelOneOrderByIDResponse represents the response for getting a specific order by ID
type AngelOneOrderByIDResponse struct {
	Status    bool                   `json:"status"`
	Message   string                 `json:"message"`
	ErrorCode string                 `json:"errorcode,omitempty"`
	Data      AngelOneOrderResponse  `json:"data"`
}

// AngelOneTradeResponse represents a single trade
type AngelOneTradeResponse struct {
	Exchange           string `json:"exchange"`
	ProductType        string `json:"producttype"`
	TradingSymbol      string `json:"tradingsymbol"`
	Instrument         string `json:"instrument"`
	SymbolGroup        string `json:"symbolgroup"`
	StrikePrice        string `json:"strikeprice"`
	OptionType         string `json:"optiontype"`
	ExpiryDate         string `json:"expirydate"`
	MarketLot          string `json:"marketlot"`
	Precision          string `json:"precision"`
	Multiplier         string `json:"multiplier"`
	TradeID            string `json:"tradeid"`
	ExchangeOrderID    string `json:"exchorderid"`
	OrderID            string `json:"orderid"`
	TransactionType    string `json:"transactiontype"`
	TradedPrice        string `json:"tradedprice"`
	TradedQuantity     string `json:"tradedquantity"`
	TradedValue        string `json:"tradedvalue"`
	TradeDate          string `json:"tradedate"`
	TradeTime          string `json:"tradetime"`
	OrderTag           string `json:"ordertag"`
	ClientCode         string `json:"clientcode"`
	PanNo              string `json:"panno"`
	SymbolToken        string `json:"symboltoken"`
	OrderPlacedBy      string `json:"orderplacedby"`
	BoardLotQuantity   string `json:"boardlotquantity"`
	ExchangeSegment    string `json:"exchangesegment"`
	ExchangeInstrumentID string `json:"exchangeinstrumentid"`
	UniqueID           string `json:"uniqueid"`
	BrokerageAmount    string `json:"brokerageamt"`
	BrokeragePercent   string `json:"brokeragepct"`
}

// AngelOneTradeListResponse represents the response for getting all trades
type AngelOneTradeListResponse struct {
	Status    bool                    `json:"status"`
	Message   string                  `json:"message"`
	ErrorCode string                  `json:"errorcode,omitempty"`
	Data      []AngelOneTradeResponse `json:"data"`
}

// Note: PlaceOrderDTO and OrderResponseDTO are defined in dhan_order_dto.go for consistency across brokers