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

// DhanBroker handles Dhan order API operations
type DhanBroker struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

// Ensure DhanBroker implements the Broker interface
var _ domain.Broker = (*DhanBroker)(nil)

// NewDhanBroker creates a new Dhan broker
func NewDhanBroker(baseURL, accessToken string) *DhanBroker {
	return &DhanBroker{
		baseURL:     baseURL,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceOrder implements the Broker interface for placing orders
func (b *DhanBroker) PlaceOrder(ctx context.Context, request *dto.PlaceOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Dhan-specific request
	dhanRequest := b.convertToPlaceOrderRequest(request)
	
	// Call Dhan API directly
	url := fmt.Sprintf("%s/orders", b.baseURL)
	
	jsonData, err := json.Marshal(dhanRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var dhanResponse dto.DhanPlaceOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&dhanResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	// Convert Dhan response to common response
	return b.convertToOrderResponse(request, &dhanResponse), nil
}

// UpdateOrder implements the Broker interface for modifying orders
func (b *DhanBroker) UpdateOrder(ctx context.Context, request *dto.ModifyOrderRequest) (*dto.OrderResponse, error) {
	// Convert common request to Dhan-specific request
	dhanRequest := b.convertToModifyOrderRequest(request)
	
	// Call Dhan API directly
	url := fmt.Sprintf("%s/orders/%s", b.baseURL, request.OrderID)
	
	jsonData, err := json.Marshal(dhanRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var dhanResponse dto.DhanModifyOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&dhanResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	// Convert Dhan response to common response
	return b.convertModifyToOrderResponse(request, &dhanResponse), nil
}

// CancelOrder implements the Broker interface for canceling orders
func (b *DhanBroker) CancelOrder(ctx context.Context, request *dto.CancelOrderRequest) (*dto.OrderResponse, error) {
	// Call Dhan API directly
	url := fmt.Sprintf("%s/orders/%s", b.baseURL, request.OrderID)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var dhanResponse dto.DhanCancelOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&dhanResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	// Convert Dhan response to common response
	return b.convertCancelToOrderResponse(request, &dhanResponse), nil
}

// GetOrders retrieves all orders for the client
func (b *DhanBroker) GetOrders(ctx context.Context) (*dto.DhanOrderListResponse, error) {
	url := fmt.Sprintf("%s/orders", b.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
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
func (b *DhanBroker) GetOrderByID(ctx context.Context, orderID string) (*dto.DhanOrderByIDResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", b.baseURL, orderID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
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
func (b *DhanBroker) GetTrades(ctx context.Context) (*dto.DhanTradeListResponse, error) {
	url := fmt.Sprintf("%s/trades", b.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
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
func (b *DhanBroker) GetOrderBook(ctx context.Context) (*dto.DhanOrderListResponse, error) {
	url := fmt.Sprintf("%s/orderbook", b.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
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
func (b *DhanBroker) GetTradeBook(ctx context.Context) (*dto.DhanTradeListResponse, error) {
	url := fmt.Sprintf("%s/tradebook", b.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
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

// Helper conversion methods

func (b *DhanBroker) convertToPlaceOrderRequest(request *dto.PlaceOrderRequest) *dto.DhanPlaceOrderRequest {
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

func (b *DhanBroker) convertToModifyOrderRequest(request *dto.ModifyOrderRequest) *dto.DhanModifyOrderRequest {
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

func (b *DhanBroker) convertToOrderResponse(request *dto.PlaceOrderRequest, dhanResponse *dto.DhanPlaceOrderResponse) *dto.OrderResponse {
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

func (b *DhanBroker) convertModifyToOrderResponse(request *dto.ModifyOrderRequest, dhanResponse *dto.DhanModifyOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: dhanResponse,
		UpdatedAt:           time.Now(),
		State:               dhanResponse.OrderStatus,
	}
}

func (b *DhanBroker) convertCancelToOrderResponse(request *dto.CancelOrderRequest, dhanResponse *dto.DhanCancelOrderResponse) *dto.OrderResponse {
	return &dto.OrderResponse{
		BrokerOrderID:       request.OrderID,
		BrokerOrderResponse: dhanResponse,
		UpdatedAt:           time.Now(),
		State:               dhanResponse.OrderStatus,
	}
}

// Legacy broker-specific methods for backward compatibility
// These are now wrappers around the internal implementation

// PlaceDhanOrder places a new order through Dhan API (legacy method)
func (b *DhanBroker) PlaceDhanOrder(ctx context.Context, request *dto.DhanPlaceOrderRequest) (*dto.DhanPlaceOrderResponse, error) {
	url := fmt.Sprintf("%s/orders", b.baseURL)
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
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

// ModifyDhanOrder modifies an existing order through Dhan API (legacy method)
func (b *DhanBroker) ModifyDhanOrder(ctx context.Context, request *dto.DhanModifyOrderRequest) (*dto.DhanModifyOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", b.baseURL, request.OrderID)
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
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

// CancelDhanOrder cancels an existing order through Dhan API (legacy method)
func (b *DhanBroker) CancelDhanOrder(ctx context.Context, request *dto.DhanCancelOrderRequest) (*dto.DhanCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", b.baseURL, request.OrderID)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("access-token", b.accessToken)

	resp, err := b.httpClient.Do(req)
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