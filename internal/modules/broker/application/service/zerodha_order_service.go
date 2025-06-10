package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"algo/internal/modules/broker/application/dto"
	"algo/internal/modules/broker/domain"
)

// ZerodhaBroker handles Zerodha Kite order API operations
type ZerodhaBroker struct {
	baseURL     string
	apiKey      string
	accessToken string
	httpClient  *http.Client
}

// Ensure ZerodhaBroker implements the Broker interface
var _ domain.Broker = (*ZerodhaBroker)(nil)

// NewZerodhaBroker creates a new Zerodha broker
func NewZerodhaBroker(baseURL, apiKey, accessToken string) *ZerodhaBroker {
	return &ZerodhaBroker{
		baseURL:     baseURL,
		apiKey:      apiKey,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetOrders retrieves all orders for the client
func (s *ZerodhaBroker) GetOrders(ctx context.Context) (*dto.ZerodhaOrderListResponse, error) {
	url := fmt.Sprintf("%s/orders", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", s.apiKey, s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.ZerodhaOrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrderHistory retrieves order history for a specific order
func (s *ZerodhaBroker) GetOrderHistory(ctx context.Context, orderID string) (*dto.ZerodhaOrderHistoryResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, orderID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", s.apiKey, s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.ZerodhaOrderHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTrades retrieves all trades for the client
func (s *ZerodhaBroker) GetTrades(ctx context.Context) (*dto.ZerodhaTradeListResponse, error) {
	url := fmt.Sprintf("%s/trades", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", s.apiKey, s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.ZerodhaTradeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetPositions retrieves all positions for the client
func (s *ZerodhaBroker) GetPositions(ctx context.Context) (*dto.ZerodhaPositionListResponse, error) {
	url := fmt.Sprintf("%s/portfolio/positions", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", s.apiKey, s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.ZerodhaPositionListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// Helper functions to convert request structs to form data

func (s *ZerodhaBroker) convertToFormData(request *dto.ZerodhaPlaceOrderRequest) string {
	params := make(map[string]string)
	
	params["tradingsymbol"] = request.TradingSymbol
	params["exchange"] = request.Exchange
	params["transaction_type"] = request.TransactionType
	params["order_type"] = request.OrderType
	params["quantity"] = fmt.Sprintf("%d", request.Quantity)
	params["product"] = request.Product
	params["validity"] = request.Validity
	
	if request.Price > 0 {
		params["price"] = fmt.Sprintf("%.2f", request.Price)
	}
	if request.DisclosedQuantity > 0 {
		params["disclosed_quantity"] = fmt.Sprintf("%d", request.DisclosedQuantity)
	}
	if request.TriggerPrice > 0 {
		params["trigger_price"] = fmt.Sprintf("%.2f", request.TriggerPrice)
	}
	if request.SquareOff > 0 {
		params["squareoff"] = fmt.Sprintf("%.2f", request.SquareOff)
	}
	if request.StopLoss > 0 {
		params["stoploss"] = fmt.Sprintf("%.2f", request.StopLoss)
	}
	if request.TrailingStopLoss > 0 {
		params["trailing_stoploss"] = fmt.Sprintf("%.2f", request.TrailingStopLoss)
	}
	if request.Variety != "" {
		params["variety"] = request.Variety
	}
	if request.Tag != "" {
		params["tag"] = request.Tag
	}
	
	// Convert to URL encoded string
	var formData string
	first := true
	for key, value := range params {
		if !first {
			formData += "&"
		}
		formData += fmt.Sprintf("%s=%s", key, value)
		first = false
	}
	
	return formData
}

func (s *ZerodhaBroker) convertModifyToFormData(request *dto.ZerodhaModifyOrderRequest) string {
	params := make(map[string]string)
	
	if request.Quantity > 0 {
		params["quantity"] = fmt.Sprintf("%d", request.Quantity)
	}
	if request.Price > 0 {
		params["price"] = fmt.Sprintf("%.2f", request.Price)
	}
	if request.OrderType != "" {
		params["order_type"] = request.OrderType
	}
	if request.TriggerPrice > 0 {
		params["trigger_price"] = fmt.Sprintf("%.2f", request.TriggerPrice)
	}
	if request.Validity != "" {
		params["validity"] = request.Validity
	}
	if request.DisclosedQuantity > 0 {
		params["disclosed_quantity"] = fmt.Sprintf("%d", request.DisclosedQuantity)
	}
	if request.ParentOrderID != "" {
		params["parent_order_id"] = request.ParentOrderID
	}
	
	// Convert to URL encoded string
	var formData string
	first := true
	for key, value := range params {
		if !first {
			formData += "&"
		}
		formData += fmt.Sprintf("%s=%s", key, value)
		first = false
	}
	
	return formData
}
// Broker interface implementation

// PlaceOrder implements the Broker interface for placing orders
func (s *ZerodhaBroker) PlaceOrder(ctx context.Context, request *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Zerodha-specific request
	zerodhaRequest := s.convertToPlaceOrderRequest(request)
	
	// Call Zerodha API directly
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, "regular")
	
	// Convert to form data as Zerodha expects form-encoded data
	formData := s.convertToFormData(zerodhaRequest)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte(formData)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", s.apiKey, s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var zerodhaResponse dto.ZerodhaPlaceOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&zerodhaResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	// Convert Zerodha response to common response
	return s.convertToOrderResponse(request, &zerodhaResponse), nil
}

// UpdateOrder implements the Broker interface for modifying orders
func (s *ZerodhaBroker) UpdateOrder(ctx context.Context, request *dto.ModifyOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Zerodha-specific request
	zerodhaRequest := s.convertToModifyOrderRequest(request)
	
	// Call Zerodha API directly
	url := fmt.Sprintf("%s/orders/%s/%s", s.baseURL, "regular", request.OrderID)
	
	// Convert to form data as Zerodha expects form-encoded data
	formData := s.convertModifyToFormData(zerodhaRequest)
	
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer([]byte(formData)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", s.apiKey, s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var zerodhaResponse dto.ZerodhaModifyOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&zerodhaResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	// Convert Zerodha response to common response
	return s.convertModifyToOrderResponse(request, &zerodhaResponse), nil
}

// CancelOrder implements the Broker interface for canceling orders
func (s *ZerodhaBroker) CancelOrder(ctx context.Context, request *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	// Call Zerodha API directly
	url := fmt.Sprintf("%s/orders/%s/%s", s.baseURL, "regular", request.OrderID)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", s.apiKey, s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var zerodhaResponse dto.ZerodhaCancelOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&zerodhaResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	// Convert Zerodha response to common response
	return s.convertCancelToOrderResponse(request, &zerodhaResponse), nil
}

// Helper conversion methods

func (s *ZerodhaBroker) convertToPlaceOrderRequest(request *dto.PlaceOrderRequest) *dto.ZerodhaPlaceOrderRequest {
	return &dto.ZerodhaPlaceOrderRequest{
		TradingSymbol:   request.Instrument,
		Exchange:        request.Exchange,
		TransactionType: request.TransactionType,
		OrderType:       request.OrderType,
		Quantity:        request.Quantity,
		Price:           request.Price,
		Product:         request.ProductType,
		Validity:        "DAY", // Default validity
		Variety:         "regular", // Default variety
		Tag:             request.AlgoID,
	}
}

func (s *ZerodhaBroker) convertToModifyOrderRequest(request *dto.ModifyOrderRequest) *dto.ZerodhaModifyOrderRequest {
	zerodhaRequest := &dto.ZerodhaModifyOrderRequest{
		OrderID: request.OrderID,
	}
	
	if request.Quantity != nil {
		zerodhaRequest.Quantity = *request.Quantity
	}
	if request.Price != nil {
		zerodhaRequest.Price = *request.Price
	}
	if request.OrderType != nil {
		zerodhaRequest.OrderType = *request.OrderType
	}
	if request.TriggerPrice != nil {
		zerodhaRequest.TriggerPrice = *request.TriggerPrice
	}
	
	return zerodhaRequest
}

func (s *ZerodhaBroker) convertToOrderResponse(request *dto.PlaceOrderRequest, zerodhaResponse *dto.ZerodhaPlaceOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		ID:                  request.ID,
		UserID:              request.UserID,
		PositionID:          request.PositionID,
		Exchange:            request.Exchange,
		Segment:             request.Segment,
		Instrument:          request.Instrument,
		TransactionType:     request.TransactionType,
		Quantity:            request.Quantity,
		Price:               request.Price,
		State:               "SENT", // Zerodha response typically just has OrderID
		ProductType:         request.ProductType,
		OrderType:           request.OrderType,
		BrokerOrderID:       zerodhaResponse.OrderID,
		BrokerOrderResponse: zerodhaResponse,
		AlgoID:              request.AlgoID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

func (s *ZerodhaBroker) convertModifyToOrderResponse(request *dto.ModifyOrderRequest, zerodhaResponse *dto.ZerodhaModifyOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: zerodhaResponse,
		UpdatedAt:           time.Now(),
		State:               "MODIFIED", // Default state for modified orders
	}
}

func (s *ZerodhaBroker) convertCancelToOrderResponse(request *dto.CancelOrderRequest, zerodhaResponse *dto.ZerodhaCancelOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: zerodhaResponse,
		UpdatedAt:           time.Now(),
		State:               "CANCELLED", // Default state for cancelled orders
	}
}