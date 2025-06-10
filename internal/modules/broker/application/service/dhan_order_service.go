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

// DhanOrderService handles Dhan order API operations
type DhanOrderService struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

// Ensure DhanOrderService implements the Broker interface
var _ domain.Broker = (*DhanOrderService)(nil)

// NewDhanOrderService creates a new Dhan order service
func NewDhanOrderService(baseURL, accessToken string) *DhanOrderService {
	return &DhanOrderService{
		baseURL:     baseURL,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceDhanOrder places a new order through Dhan API
func (s *DhanOrderService) PlaceDhanOrder(ctx context.Context, request *dto.DhanPlaceOrderRequest) (*dto.DhanPlaceOrderResponse, error) {
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
	req.Header.Set("access-token", s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.DhanPlaceOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// ModifyDhanOrder modifies an existing order through Dhan API
func (s *DhanOrderService) ModifyDhanOrder(ctx context.Context, request *dto.DhanModifyOrderRequest) (*dto.DhanModifyOrderResponse, error) {
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
	req.Header.Set("access-token", s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.DhanModifyOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CancelDhanOrder cancels an existing order through Dhan API
func (s *DhanOrderService) CancelDhanOrder(ctx context.Context, request *dto.DhanCancelOrderRequest) (*dto.DhanCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, request.OrderID)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.DhanCancelOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrders retrieves all orders for the client
func (s *DhanOrderService) GetOrders(ctx context.Context) (*dto.DhanOrderListResponse, error) {
	url := fmt.Sprintf("%s/orders", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.DhanOrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrderByID retrieves a specific order by ID
func (s *DhanOrderService) GetOrderByID(ctx context.Context, orderID string) (*dto.DhanOrderByIDResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, orderID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.DhanOrderByIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTrades retrieves all trades for the client
func (s *DhanOrderService) GetTrades(ctx context.Context) (*dto.DhanTradeListResponse, error) {
	url := fmt.Sprintf("%s/trades", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.DhanTradeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrderBook retrieves the order book for the client
func (s *DhanOrderService) GetOrderBook(ctx context.Context) (*dto.DhanOrderListResponse, error) {
	url := fmt.Sprintf("%s/orderbook", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.DhanOrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTradeBook retrieves the trade book for the client
func (s *DhanOrderService) GetTradeBook(ctx context.Context) (*dto.DhanTradeListResponse, error) {
	url := fmt.Sprintf("%s/tradebook", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.DhanTradeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// Broker interface implementation

// PlaceOrder implements the Broker interface for placing orders
func (s *DhanOrderService) PlaceOrder(ctx context.Context, request *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Dhan-specific request
	dhanRequest := s.convertToPlaceOrderRequest(request)
	
	// Call the existing Dhan-specific PlaceOrder method
	dhanResponse, err := s.PlaceDhanOrder(ctx, dhanRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Dhan response to common response
	return s.convertToOrderResponse(request, dhanResponse), nil
}

// UpdateOrder implements the Broker interface for modifying orders
func (s *DhanOrderService) UpdateOrder(ctx context.Context, request *dto.ModifyOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Dhan-specific request
	dhanRequest := s.convertToModifyOrderRequest(request)
	
	// Call the existing Dhan-specific ModifyOrder method
	dhanResponse, err := s.ModifyDhanOrder(ctx, dhanRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Dhan response to common response
	return s.convertModifyToOrderResponse(request, dhanResponse), nil
}

// CancelOrder implements the Broker interface for canceling orders
func (s *DhanOrderService) CancelOrder(ctx context.Context, request *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Dhan-specific request
	dhanRequest := &dto.DhanCancelOrderRequest{
		OrderID: request.OrderID,
	}
	
	// Call the existing Dhan-specific CancelOrder method
	dhanResponse, err := s.CancelDhanOrder(ctx, dhanRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert Dhan response to common response
	return s.convertCancelToOrderResponse(request, dhanResponse), nil
}

// Helper conversion methods

func (s *DhanOrderService) convertToPlaceOrderRequest(request *dto.PlaceOrderRequest) *dto.DhanPlaceOrderRequest {
	return &dto.DhanPlaceOrderRequest{
		DhanClientID:       "", // This would need to be populated from user context
		TransactionType:    request.TransactionType,
		ExchangeSegment:    fmt.Sprintf("%s_%s", request.Exchange, request.Segment),
		ProductType:        request.ProductType,
		OrderType:          request.OrderType,
		Validity:           "DAY", // Default validity
		SecurityID:         request.Instrument,
		Quantity:           request.Quantity,
		Price:              request.Price,
		DisclosedQuantity:  0,
		TriggerPrice:       0,
		AfterMarketOrder:   false,
		AmoTime:           "",
		BolBlockDeal:      false,
		CorrelationID:     request.AlgoID,
	}
}

func (s *DhanOrderService) convertToModifyOrderRequest(request *dto.ModifyOrderRequest) *dto.DhanModifyOrderRequest {
	dhanRequest := &dto.DhanModifyOrderRequest{
		OrderID: request.OrderID,
	}
	
	if request.Quantity != nil {
		dhanRequest.Quantity = *request.Quantity
	}
	if request.Price != nil {
		dhanRequest.Price = *request.Price
	}
	if request.OrderType != nil {
		dhanRequest.OrderType = *request.OrderType
	}
	if request.TriggerPrice != nil {
		dhanRequest.TriggerPrice = *request.TriggerPrice
	}
	
	return dhanRequest
}

func (s *DhanOrderService) convertToOrderResponse(request *dto.PlaceOrderRequest, dhanResponse *dto.DhanPlaceOrderResponse) *dto.OrderResponse {
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
		State:               dhanResponse.OrderStatus,
		ProductType:         request.ProductType,
		OrderType:           request.OrderType,
		BrokerOrderID:       dhanResponse.OrderID,
		BrokerOrderResponse: dhanResponse,
		AlgoID:              request.AlgoID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

func (s *DhanOrderService) convertModifyToOrderResponse(request *dto.ModifyOrderRequest, dhanResponse *dto.DhanModifyOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: dhanResponse,
		UpdatedAt:           time.Now(),
		State:               dhanResponse.OrderStatus,
	}
}

func (s *DhanOrderService) convertCancelToOrderResponse(request *dto.CancelOrderRequest, dhanResponse *dto.DhanCancelOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: dhanResponse,
		UpdatedAt:           time.Now(),
		State:               dhanResponse.OrderStatus,
	}
}