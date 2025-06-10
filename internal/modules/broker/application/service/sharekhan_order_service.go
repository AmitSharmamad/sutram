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

// SharekhanOrderService handles Sharekhan order API operations
type SharekhanOrderService struct {
	baseURL     string
	apiKey      string
	accessToken string
	httpClient  *http.Client
}

// Ensure SharekhanOrderService implements the Broker interface
var _ domain.Broker = (*SharekhanOrderService)(nil)

// NewSharekhanOrderService creates a new Sharekhan order service
func NewSharekhanOrderService(baseURL, apiKey, accessToken string) *SharekhanOrderService {
	return &SharekhanOrderService{
		baseURL:     baseURL,
		apiKey:      apiKey,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceSharekhanOrder places a new order through Sharekhan API
func (s *SharekhanOrderService) PlaceSharekhanOrder(ctx context.Context, request *dto.SharekhanPlaceOrderRequest) (*dto.SharekhanPlaceOrderResponse, error) {
	url := fmt.Sprintf("%s/orders", s.baseURL)
	
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
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.SharekhanPlaceOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// ModifySharekhanOrder modifies an existing order through Sharekhan API
func (s *SharekhanOrderService) ModifySharekhanOrder(ctx context.Context, request *dto.SharekhanModifyOrderRequest) (*dto.SharekhanModifyOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, request.OrderNumber)
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.SharekhanModifyOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CancelSharekhanOrder cancels an existing order through Sharekhan API
func (s *SharekhanOrderService) CancelSharekhanOrder(ctx context.Context, request *dto.SharekhanCancelOrderRequest) (*dto.SharekhanCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, request.OrderNumber)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.SharekhanCancelOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrders retrieves all orders for the client
func (s *SharekhanOrderService) GetOrders(ctx context.Context) (*dto.SharekhanOrderListResponse, error) {
	url := fmt.Sprintf("%s/orders", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.SharekhanOrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrderByID retrieves a specific order by ID
func (s *SharekhanOrderService) GetOrderByID(ctx context.Context, orderNumber string) (*dto.SharekhanOrderByIDResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, orderNumber)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.SharekhanOrderByIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTrades retrieves all trades for the client
func (s *SharekhanOrderService) GetTrades(ctx context.Context) (*dto.SharekhanTradeListResponse, error) {
	url := fmt.Sprintf("%s/trades", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.SharekhanTradeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrderBook retrieves the order book for the client
func (s *SharekhanOrderService) GetOrderBook(ctx context.Context) (*dto.SharekhanOrderBookResponse, error) {
	url := fmt.Sprintf("%s/orderbook", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.SharekhanOrderBookResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTradeBook retrieves the trade book for the client
func (s *SharekhanOrderService) GetTradeBook(ctx context.Context) (*dto.SharekhanTradeBookResponse, error) {
	url := fmt.Sprintf("%s/tradebook", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.SharekhanTradeBookResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// Broker interface implementation

// PlaceOrder implements the Broker interface for placing orders
func (s *SharekhanOrderService) PlaceOrder(ctx context.Context, request *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Sharekhan-specific request
	sharekhanRequest := s.convertToPlaceOrderRequest(request)
	
	// Call the existing Sharekhan-specific PlaceOrder method
	sharekhanResponse, err := s.PlaceSharekhanOrder(ctx, sharekhanRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Sharekhan response to common response
	return s.convertToOrderResponse(request, sharekhanResponse), nil
}

// UpdateOrder implements the Broker interface for modifying orders
func (s *SharekhanOrderService) UpdateOrder(ctx context.Context, request *dto.ModifyOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Sharekhan-specific request
	sharekhanRequest := s.convertToModifyOrderRequest(request)
	
	// Call the existing Sharekhan-specific ModifyOrder method
	sharekhanResponse, err := s.ModifySharekhanOrder(ctx, sharekhanRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Sharekhan response to common response
	return s.convertModifyToOrderResponse(request, sharekhanResponse), nil
}

// CancelOrder implements the Broker interface for canceling orders
func (s *SharekhanOrderService) CancelOrder(ctx context.Context, request *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Sharekhan-specific request
	sharekhanRequest := &dto.SharekhanCancelOrderRequest{
		OrderNumber: request.OrderID,
	}
	
	// Call the existing Sharekhan-specific CancelOrder method
	sharekhanResponse, err := s.CancelSharekhanOrder(ctx, sharekhanRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Sharekhan response to common response
	return s.convertCancelToOrderResponse(request, sharekhanResponse), nil
}

// Helper conversion methods

func (s *SharekhanOrderService) convertToPlaceOrderRequest(request *dto.PlaceOrderRequest) *dto.SharekhanPlaceOrderRequest {
	return &dto.SharekhanPlaceOrderRequest{
		TradingSymbol:   request.Instrument,
		Exchange:        request.Exchange,
		TransactionType: request.TransactionType,
		OrderType:       request.OrderType,
		ProductCode:     request.ProductType,
		Quantity:        request.Quantity,
		Price:           request.Price,
		Validity:        "DAY", // Default validity
		DisclosedQty:    0,     // Default
		TriggerPrice:    0,     // Default
	}
}

func (s *SharekhanOrderService) convertToModifyOrderRequest(request *dto.ModifyOrderRequest) *dto.SharekhanModifyOrderRequest {
	sharekhanRequest := &dto.SharekhanModifyOrderRequest{
		OrderNumber: request.OrderID,
	}
	
	if request.Quantity != nil {
		sharekhanRequest.Quantity = *request.Quantity
	}
	if request.Price != nil {
		sharekhanRequest.Price = *request.Price
	}
	if request.OrderType != nil {
		sharekhanRequest.OrderType = *request.OrderType
	}
	if request.TriggerPrice != nil {
		sharekhanRequest.TriggerPrice = *request.TriggerPrice
	}
	
	return sharekhanRequest
}

func (s *SharekhanOrderService) convertToOrderResponse(request *dto.PlaceOrderRequest, sharekhanResponse *dto.SharekhanPlaceOrderResponse) *dto.OrderResponse {
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
		State:               sharekhanResponse.OrderStatus,
		ProductType:         request.ProductType,
		OrderType:           request.OrderType,
		BrokerOrderID:       sharekhanResponse.OrderNumber,
		BrokerOrderResponse: sharekhanResponse,
		AlgoID:              request.AlgoID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

func (s *SharekhanOrderService) convertModifyToOrderResponse(request *dto.ModifyOrderRequest, sharekhanResponse *dto.SharekhanModifyOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: sharekhanResponse,
		UpdatedAt:           time.Now(),
		State:               sharekhanResponse.OrderStatus,
	}
}

func (s *SharekhanOrderService) convertCancelToOrderResponse(request *dto.CancelOrderRequest, sharekhanResponse *dto.SharekhanCancelOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: sharekhanResponse,
		UpdatedAt:           time.Now(),
		State:               sharekhanResponse.OrderStatus,
	}
}