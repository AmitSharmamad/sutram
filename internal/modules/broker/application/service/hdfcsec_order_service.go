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

// HDFCSecBroker handles HDFC Securities order API operations
type HDFCSecBroker struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

// Ensure HDFCSecBroker implements the Broker interface
var _ domain.Broker = (*HDFCSecBroker)(nil)

// NewHDFCSecBroker creates a new HDFC Securities order service
func NewHDFCSecBroker(baseURL, accessToken string) *HDFCSecBroker {
	return &HDFCSecBroker{
		baseURL:     baseURL,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceHDFCSecOrder places a new order through HDFC Securities API
func (s *HDFCSecBroker) PlaceHDFCSecOrder(ctx context.Context, request *dto.HDFCSecPlaceOrderRequest) (*dto.HDFCSecPlaceOrderResponse, error) {
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

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecPlaceOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// ModifyHDFCSecOrder modifies an existing order through HDFC Securities API
func (s *HDFCSecBroker) ModifyHDFCSecOrder(ctx context.Context, request *dto.HDFCSecModifyOrderRequest) (*dto.HDFCSecModifyOrderResponse, error) {
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

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecModifyOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CancelHDFCSecOrder cancels an existing order through HDFC Securities API
func (s *HDFCSecBroker) CancelHDFCSecOrder(ctx context.Context, request *dto.HDFCSecCancelOrderRequest) (*dto.HDFCSecCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, request.OrderID)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecCancelOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrders retrieves all orders for the client
func (s *HDFCSecBroker) GetOrders(ctx context.Context) (*dto.HDFCSecOrderListResponse, error) {
	url := fmt.Sprintf("%s/orders", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecOrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrderByID retrieves a specific order by ID
func (s *HDFCSecBroker) GetOrderByID(ctx context.Context, orderID string) (*dto.HDFCSecOrderByIDResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, orderID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecOrderByIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTrades retrieves all trades for the client
func (s *HDFCSecBroker) GetTrades(ctx context.Context) (*dto.HDFCSecTradeListResponse, error) {
	url := fmt.Sprintf("%s/trades", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecTradeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrderBook retrieves the order book for the client
func (s *HDFCSecBroker) GetOrderBook(ctx context.Context) (*dto.HDFCSecOrderListResponse, error) {
	url := fmt.Sprintf("%s/orderbook", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecOrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTradeBook retrieves the trade book for the client
func (s *HDFCSecBroker) GetTradeBook(ctx context.Context) (*dto.HDFCSecTradeListResponse, error) {
	url := fmt.Sprintf("%s/tradebook", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecTradeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetPositions retrieves all positions for the client
func (s *HDFCSecBroker) GetPositions(ctx context.Context) (*dto.HDFCSecPositionListResponse, error) {
	url := fmt.Sprintf("%s/positions", s.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response dto.HDFCSecPositionListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// Broker interface implementation

// PlaceOrder implements the Broker interface for placing orders
func (s *HDFCSecBroker) PlaceOrder(ctx context.Context, request *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to HDFC Securities-specific request
	hdfcRequest := s.convertToPlaceOrderRequest(request)
	
	// Call the existing HDFC Securities-specific PlaceOrder method
	hdfcResponse, err := s.PlaceHDFCSecOrder(ctx, hdfcRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert HDFC Securities response to common response
	return s.convertToOrderResponse(request, hdfcResponse), nil
}

// UpdateOrder implements the Broker interface for modifying orders
func (s *HDFCSecBroker) UpdateOrder(ctx context.Context, request *dto.ModifyOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to HDFC Securities-specific request
	hdfcRequest := s.convertToModifyOrderRequest(request)
	
	// Call the existing HDFC Securities-specific ModifyOrder method
	hdfcResponse, err := s.ModifyHDFCSecOrder(ctx, hdfcRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert HDFC Securities response to common response
	return s.convertModifyToOrderResponse(request, hdfcResponse), nil
}

// CancelOrder implements the Broker interface for canceling orders
func (s *HDFCSecBroker) CancelOrder(ctx context.Context, request *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to HDFC Securities-specific request
	hdfcRequest := &dto.HDFCSecCancelOrderRequest{
		OrderID: request.OrderID,
	}
	
	// Call the existing HDFC Securities-specific CancelOrder method
	hdfcResponse, err := s.CancelHDFCSecOrder(ctx, hdfcRequest)
	if err != nil {
		return nil, err
	}
	
	// Convert HDFC Securities response to common response
	return s.convertCancelToOrderResponse(request, hdfcResponse), nil
}

// Helper conversion methods

func (s *HDFCSecBroker) convertToPlaceOrderRequest(request *dto.PlaceOrderRequest) *dto.HDFCSecPlaceOrderRequest {
	return &dto.HDFCSecPlaceOrderRequest{
		ClientCode:      "", // This would need to be populated from user context
		Symbol:          request.Instrument,
		Exchange:        request.Exchange,
		Segment:         request.Segment,
		TransactionType: request.TransactionType,
		OrderType:       request.OrderType,
		ProductType:     request.ProductType,
		Quantity:        request.Quantity,
		Price:           request.Price,
		Validity:        "DAY", // Default validity
		DisclosedQty:    0,     // Default
		TriggerPrice:    0,     // Default
	}
}

func (s *HDFCSecBroker) convertToModifyOrderRequest(request *dto.ModifyOrderRequest) *dto.HDFCSecModifyOrderRequest {
	hdfcRequest := &dto.HDFCSecModifyOrderRequest{
		OrderID: request.OrderID,
	}
	
	if request.Quantity != nil {
		hdfcRequest.Quantity = *request.Quantity
	}
	if request.Price != nil {
		hdfcRequest.Price = *request.Price
	}
	if request.OrderType != nil {
		hdfcRequest.OrderType = *request.OrderType
	}
	if request.TriggerPrice != nil {
		hdfcRequest.TriggerPrice = *request.TriggerPrice
	}
	
	return hdfcRequest
}

func (s *HDFCSecBroker) convertToOrderResponse(request *dto.PlaceOrderRequest, hdfcResponse *dto.HDFCSecPlaceOrderResponse) *dto.OrderResponse {
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
		State:               hdfcResponse.OrderStatus,
		ProductType:         request.ProductType,
		OrderType:           request.OrderType,
		BrokerOrderID:       hdfcResponse.OrderID,
		BrokerOrderResponse: hdfcResponse,
		AlgoID:              request.AlgoID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

func (s *HDFCSecBroker) convertModifyToOrderResponse(request *dto.ModifyOrderRequest, hdfcResponse *dto.HDFCSecModifyOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: hdfcResponse,
		UpdatedAt:           time.Now(),
		State:               hdfcResponse.OrderStatus,
	}
}

func (s *HDFCSecBroker) convertCancelToOrderResponse(request *dto.CancelOrderRequest, hdfcResponse *dto.HDFCSecCancelOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: hdfcResponse,
		UpdatedAt:           time.Now(),
		State:               hdfcResponse.OrderStatus,
	}
}