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

// AngelOneOrderService handles AngelOne Smart API order operations
type AngelOneOrderService struct {
	baseURL     string
	apiKey      string
	accessToken string
	httpClient  *http.Client
}

// Ensure AngelOneOrderService implements the Broker interface
var _ domain.Broker = (*AngelOneOrderService)(nil)

// NewAngelOneOrderService creates a new AngelOne order service
func NewAngelOneOrderService(baseURL, apiKey, accessToken string) *AngelOneOrderService {
	return &AngelOneOrderService{
		baseURL:     baseURL,
		apiKey:      apiKey,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceOrder places a new order through AngelOne Smart API
func (s *AngelOneOrderService) PlaceAngelOneOrder(ctx context.Context, request *dto.AngelOnePlaceOrderRequest) (*dto.AngelOnePlaceOrderResponse, error) {
	url := fmt.Sprintf("%s/rest/secure/angelbroking/order/v1/placeOrder", s.baseURL)
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")
	req.Header.Set("X-ClientLocalIP", "127.0.0.1")
	req.Header.Set("X-ClientPublicIP", "127.0.0.1")
	req.Header.Set("X-MACAddress", "00:00:00:00:00:00")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-PrivateKey", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response dto.AngelOnePlaceOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API request failed: %s (error code: %s)", response.Message, response.ErrorCode)
	}

	return &response, nil
}

// ModifyOrder modifies an existing order through AngelOne Smart API
func (s *AngelOneOrderService) ModifyAngelOneOrder(ctx context.Context, request *dto.AngelOneModifyOrderRequest) (*dto.AngelOneModifyOrderResponse, error) {
	url := fmt.Sprintf("%s/rest/secure/angelbroking/order/v1/modifyOrder", s.baseURL)
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")
	req.Header.Set("X-ClientLocalIP", "127.0.0.1")
	req.Header.Set("X-ClientPublicIP", "127.0.0.1")
	req.Header.Set("X-MACAddress", "00:00:00:00:00:00")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-PrivateKey", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response dto.AngelOneModifyOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API request failed: %s (error code: %s)", response.Message, response.ErrorCode)
	}

	return &response, nil
}

// CancelOrder cancels an existing order through AngelOne Smart API
func (s *AngelOneOrderService) CancelAngelOneOrder(ctx context.Context, request *dto.AngelOneCancelOrderRequest) (*dto.AngelOneCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/rest/secure/angelbroking/order/v1/cancelOrder", s.baseURL)
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")
	req.Header.Set("X-ClientLocalIP", "127.0.0.1")
	req.Header.Set("X-ClientPublicIP", "127.0.0.1")
	req.Header.Set("X-MACAddress", "00:00:00:00:00:00")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-PrivateKey", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response dto.AngelOneCancelOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API request failed: %s (error code: %s)", response.Message, response.ErrorCode)
	}

	return &response, nil
}

// GetOrders retrieves all orders for the client
func (s *AngelOneOrderService) GetOrders(ctx context.Context) (*dto.AngelOneOrderListResponse, error) {
	url := fmt.Sprintf("%s/rest/secure/angelbroking/order/v1/getOrderBook", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")
	req.Header.Set("X-ClientLocalIP", "127.0.0.1")
	req.Header.Set("X-ClientPublicIP", "127.0.0.1")
	req.Header.Set("X-MACAddress", "00:00:00:00:00:00")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-PrivateKey", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response dto.AngelOneOrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API request failed: %s (error code: %s)", response.Message, response.ErrorCode)
	}

	return &response, nil
}

// GetOrderByID retrieves a specific order by ID
func (s *AngelOneOrderService) GetOrderByID(ctx context.Context, orderID string) (*dto.AngelOneOrderByIDResponse, error) {
	url := fmt.Sprintf("%s/rest/secure/angelbroking/order/v1/details/%s", s.baseURL, orderID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")
	req.Header.Set("X-ClientLocalIP", "127.0.0.1")
	req.Header.Set("X-ClientPublicIP", "127.0.0.1")
	req.Header.Set("X-MACAddress", "00:00:00:00:00:00")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-PrivateKey", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response dto.AngelOneOrderByIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API request failed: %s (error code: %s)", response.Message, response.ErrorCode)
	}

	return &response, nil
}

// GetTrades retrieves all trades for the client
func (s *AngelOneOrderService) GetTrades(ctx context.Context) (*dto.AngelOneTradeListResponse, error) {
	url := fmt.Sprintf("%s/rest/secure/angelbroking/order/v1/getTradeBook", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")
	req.Header.Set("X-ClientLocalIP", "127.0.0.1")
	req.Header.Set("X-ClientPublicIP", "127.0.0.1")
	req.Header.Set("X-MACAddress", "00:00:00:00:00:00")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-PrivateKey", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response dto.AngelOneTradeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API request failed: %s (error code: %s)", response.Message, response.ErrorCode)
	}

	return &response, nil
}

// GetOrderBook retrieves the order book for the client (alias for GetOrders)
func (s *AngelOneOrderService) GetOrderBook(ctx context.Context) (*dto.AngelOneOrderListResponse, error) {
	return s.GetOrders(ctx)
}

// GetTradeBook retrieves the trade book for the client (alias for GetTrades)
func (s *AngelOneOrderService) GetTradeBook(ctx context.Context) (*dto.AngelOneTradeListResponse, error) {
	return s.GetTrades(ctx)
}

// Broker interface implementation

// PlaceOrder implements the Broker interface for placing orders
func (s *AngelOneOrderService) PlaceOrder(ctx context.Context, request *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to AngelOne-specific request
	angelRequest := s.convertToPlaceOrderRequest(request)
	
	// Call the existing AngelOne-specific PlaceOrder method
	angelResponse, err := s.PlaceAngelOneOrder(ctx, angelRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert AngelOne response to common response
	return s.convertToOrderResponse(request, angelResponse), nil
}

// UpdateOrder implements the Broker interface for modifying orders
func (s *AngelOneOrderService) UpdateOrder(ctx context.Context, request *dto.ModifyOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to AngelOne-specific request
	angelRequest := s.convertToModifyOrderRequest(request)
	
	// Call the existing AngelOne-specific ModifyOrder method
	angelResponse, err := s.ModifyAngelOneOrder(ctx, angelRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert AngelOne response to common response
	return s.convertModifyToOrderResponse(request, angelResponse), nil
}

// CancelOrder implements the Broker interface for canceling orders
func (s *AngelOneOrderService) CancelOrder(ctx context.Context, request *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to AngelOne-specific request
	angelRequest := &dto.AngelOneCancelOrderRequest{
		OrderID: request.OrderID,
		Variety: "regular", // Default variety for AngelOne
	}
	
	// Call the existing AngelOne-specific CancelOrder method
	angelResponse, err := s.CancelAngelOneOrder(ctx, angelRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert AngelOne response to common response
	return s.convertCancelToOrderResponse(request, angelResponse), nil
}

// Helper conversion methods

func (s *AngelOneOrderService) convertToPlaceOrderRequest(request *dto.PlaceOrderRequest) *dto.AngelOnePlaceOrderRequest {
	return &dto.AngelOnePlaceOrderRequest{
		Variety:         "regular", // Default variety
		TradingSymbol:   request.Instrument,
		SymbolToken:     "",        // This would need to be populated from symbol mapping
		Exchange:        request.Exchange,
		TransactionType: request.TransactionType,
		OrderType:       request.OrderType,
		ProductType:     request.ProductType,
		Duration:        "DAY",     // Default duration
		Quantity:        fmt.Sprintf("%d", request.Quantity),
		Price:           fmt.Sprintf("%.2f", request.Price),
		DisclosedQuantity: "0",     // Default
		TriggerPrice:    "0",       // Default
		Tag:             request.AlgoID,
	}
}

func (s *AngelOneOrderService) convertToModifyOrderRequest(request *dto.ModifyOrderRequest) *dto.AngelOneModifyOrderRequest {
	angelRequest := &dto.AngelOneModifyOrderRequest{
		Variety:     "regular", // Default variety
		OrderID:     request.OrderID,
		OrderType:   "LIMIT",   // Default
		ProductType: "DELIVERY", // Default
		Duration:    "DAY",     // Default
	}
	
	if request.Quantity != nil {
		angelRequest.Quantity = fmt.Sprintf("%d", *request.Quantity)
	}
	if request.Price != nil {
		angelRequest.Price = fmt.Sprintf("%.2f", *request.Price)
	}
	if request.OrderType != nil {
		angelRequest.OrderType = *request.OrderType
	}
	if request.TriggerPrice != nil {
		angelRequest.TriggerPrice = fmt.Sprintf("%.2f", *request.TriggerPrice)
	}
	
	return angelRequest
}

func (s *AngelOneOrderService) convertToOrderResponse(request *dto.PlaceOrderRequest, angelResponse *dto.AngelOnePlaceOrderResponse) *dto.OrderResponse {
	state := "SENT"
	if !angelResponse.Status {
		state = "REJECTED"
	}
	
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
		State:               state,
		ProductType:         request.ProductType,
		OrderType:           request.OrderType,
		BrokerOrderID:       angelResponse.Data.OrderID,
		BrokerOrderResponse: angelResponse,
		AlgoID:              request.AlgoID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

func (s *AngelOneOrderService) convertModifyToOrderResponse(request *dto.ModifyOrderRequest, angelResponse *dto.AngelOneModifyOrderResponse) *dto.OrderResponse {
	state := "MODIFIED"
	if !angelResponse.Status {
		state = "MODIFY_REJECTED"
	}
	
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: angelResponse,
		UpdatedAt:           time.Now(),
		State:               state,
	}
}

func (s *AngelOneOrderService) convertCancelToOrderResponse(request *dto.CancelOrderRequest, angelResponse *dto.AngelOneCancelOrderResponse) *dto.OrderResponse {
	state := "CANCELLED"
	if !angelResponse.Status {
		state = "CANCEL_REJECTED"
	}
	
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: angelResponse,
		UpdatedAt:           time.Now(),
		State:               state,
	}
}
