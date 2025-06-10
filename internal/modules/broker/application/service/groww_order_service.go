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

// GrowwOrderService handles Groww order API operations
type GrowwOrderService struct {
	baseURL     string
	apiKey      string
	accessToken string
	httpClient  *http.Client
}

// Ensure GrowwOrderService implements the Broker interface
var _ domain.Broker = (*GrowwOrderService)(nil)

// NewGrowwOrderService creates a new Groww order service
func NewGrowwOrderService(baseURL, apiKey, accessToken string) *GrowwOrderService {
	return &GrowwOrderService{
		baseURL:     baseURL,
		apiKey:      apiKey,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceOrder places a new order through Groww API
func (s *GrowwOrderService) PlaceGrowwOrder(ctx context.Context, request *dto.GrowwPlaceOrderRequest) (*dto.GrowwPlaceOrderResponse, error) {
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

	var response dto.GrowwPlaceOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// ModifyOrder modifies an existing order through Groww API
func (s *GrowwOrderService) ModifyGrowwOrder(ctx context.Context, request *dto.GrowwModifyOrderRequest) (*dto.GrowwModifyOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, request.OrderID)
	
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

	var response dto.GrowwModifyOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CancelOrder cancels an existing order through Groww API
func (s *GrowwOrderService) CancelGrowwOrder(ctx context.Context, request *dto.GrowwCancelOrderRequest) (*dto.GrowwCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, request.OrderID)
	
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

	var response dto.GrowwCancelOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrders retrieves all orders for the client
func (s *GrowwOrderService) GetOrders(ctx context.Context) (*dto.GrowwOrderListResponse, error) {
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

	var response dto.GrowwOrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrderByID retrieves a specific order by ID
func (s *GrowwOrderService) GetOrderByID(ctx context.Context, orderID string) (*dto.GrowwOrderByIDResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, orderID)
	
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

	var response dto.GrowwOrderByIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTrades retrieves all trades for the client
func (s *GrowwOrderService) GetTrades(ctx context.Context) (*dto.GrowwTradeListResponse, error) {
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

	var response dto.GrowwTradeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// Broker interface implementation

// PlaceOrder implements the Broker interface for placing orders
func (s *GrowwOrderService) PlaceOrder(ctx context.Context, request *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Groww-specific request
	growwRequest := s.convertToPlaceOrderRequest(request)
	
	// Call the existing Groww-specific PlaceOrder method
	growwResponse, err := s.PlaceGrowwOrder(ctx, growwRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Groww response to common response
	return s.convertToOrderResponse(request, growwResponse), nil
}

// UpdateOrder implements the Broker interface for modifying orders
func (s *GrowwOrderService) UpdateOrder(ctx context.Context, request *dto.ModifyOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Groww-specific request
	growwRequest := s.convertToModifyOrderRequest(request)
	
	// Call the existing Groww-specific ModifyOrder method
	growwResponse, err := s.ModifyGrowwOrder(ctx, growwRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Groww response to common response
	return s.convertModifyToOrderResponse(request, growwResponse), nil
}

// CancelOrder implements the Broker interface for canceling orders
func (s *GrowwOrderService) CancelOrder(ctx context.Context, request *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Groww-specific request
	growwRequest := &dto.GrowwCancelOrderRequest{
		OrderID: request.OrderID,
	}
	
	// Call the existing Groww-specific CancelOrder method
	growwResponse, err := s.CancelGrowwOrder(ctx, growwRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Groww response to common response
	return s.convertCancelToOrderResponse(request, growwResponse), nil
}

// Helper conversion methods

func (s *GrowwOrderService) convertToPlaceOrderRequest(request *dto.PlaceOrderRequest) *dto.GrowwPlaceOrderRequest {
	return &dto.GrowwPlaceOrderRequest{
		Symbol:          request.Instrument,
		Exchange:        request.Exchange,
		Segment:         request.Segment,
		Side:            request.TransactionType,
		OrderType:       request.OrderType,
		ProductType:     request.ProductType,
		Quantity:        request.Quantity,
		Price:           request.Price,
		Validity:        "DAY", // Default validity
		DisclosedQty:    0,     // Default
		StopPrice:       0,     // Default
	}
}

func (s *GrowwOrderService) convertToModifyOrderRequest(request *dto.ModifyOrderRequest) *dto.GrowwModifyOrderRequest {
	growwRequest := &dto.GrowwModifyOrderRequest{
		OrderID:   request.OrderID,
		OrderType: "LIMIT",  // Default
		Quantity:  1,        // Default
		Validity:  "DAY",    // Default
	}
	
	if request.Quantity != nil {
		growwRequest.Quantity = *request.Quantity
	}
	if request.Price != nil {
		growwRequest.Price = *request.Price
	}
	if request.OrderType != nil {
		growwRequest.OrderType = *request.OrderType
	}
	if request.TriggerPrice != nil {
		growwRequest.StopPrice = *request.TriggerPrice
	}
	
	return growwRequest
}

func (s *GrowwOrderService) convertToOrderResponse(request *dto.PlaceOrderRequest, growwResponse *dto.GrowwPlaceOrderResponse) *dto.OrderResponse {
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
		State:               "SENT", // Default state for placed orders
		ProductType:         request.ProductType,
		OrderType:           request.OrderType,
		BrokerOrderID:       growwResponse.OrderID,
		BrokerOrderResponse: growwResponse,
		AlgoID:              request.AlgoID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

func (s *GrowwOrderService) convertModifyToOrderResponse(request *dto.ModifyOrderRequest, growwResponse *dto.GrowwModifyOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: growwResponse,
		UpdatedAt:           time.Now(),
		State:               growwResponse.Status, // Use the Status field
	}
}

func (s *GrowwOrderService) convertCancelToOrderResponse(request *dto.CancelOrderRequest, growwResponse *dto.GrowwCancelOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: growwResponse,
		UpdatedAt:           time.Now(),
		State:               growwResponse.Status, // Use the Status field
	}
}
