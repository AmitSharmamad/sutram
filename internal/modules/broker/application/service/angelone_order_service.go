package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"algo/internal/modules/broker/application/dto"
)

// AngelOneOrderService handles AngelOne Smart API order operations
type AngelOneOrderService struct {
	baseURL     string
	apiKey      string
	accessToken string
	httpClient  *http.Client
}

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
func (s *AngelOneOrderService) PlaceOrder(ctx context.Context, request *dto.AngelOnePlaceOrderRequest) (*dto.AngelOnePlaceOrderResponse, error) {
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
func (s *AngelOneOrderService) ModifyOrder(ctx context.Context, request *dto.AngelOneModifyOrderRequest) (*dto.AngelOneModifyOrderResponse, error) {
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
func (s *AngelOneOrderService) CancelOrder(ctx context.Context, request *dto.AngelOneCancelOrderRequest) (*dto.AngelOneCancelOrderResponse, error) {
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