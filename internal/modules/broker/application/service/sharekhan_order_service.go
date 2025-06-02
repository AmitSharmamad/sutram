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

// SharekhanOrderService handles Sharekhan order API operations
type SharekhanOrderService struct {
	baseURL     string
	apiKey      string
	accessToken string
	httpClient  *http.Client
}

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

// PlaceOrder places a new order through Sharekhan API
func (s *SharekhanOrderService) PlaceOrder(ctx context.Context, request *dto.SharekhanPlaceOrderRequest) (*dto.SharekhanPlaceOrderResponse, error) {
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

// ModifyOrder modifies an existing order through Sharekhan API
func (s *SharekhanOrderService) ModifyOrder(ctx context.Context, request *dto.SharekhanModifyOrderRequest) (*dto.SharekhanModifyOrderResponse, error) {
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

// CancelOrder cancels an existing order through Sharekhan API
func (s *SharekhanOrderService) CancelOrder(ctx context.Context, request *dto.SharekhanCancelOrderRequest) (*dto.SharekhanCancelOrderResponse, error) {
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