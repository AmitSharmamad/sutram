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

// ZerodhaOrderService handles Zerodha Kite order API operations
type ZerodhaOrderService struct {
	baseURL     string
	apiKey      string
	accessToken string
	httpClient  *http.Client
}

// NewZerodhaOrderService creates a new Zerodha order service
func NewZerodhaOrderService(baseURL, apiKey, accessToken string) *ZerodhaOrderService {
	return &ZerodhaOrderService{
		baseURL:     baseURL,
		apiKey:      apiKey,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceOrder places a new order through Zerodha Kite API
func (s *ZerodhaOrderService) PlaceOrder(ctx context.Context, request *dto.ZerodhaPlaceOrderRequest) (*dto.ZerodhaPlaceOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", s.baseURL, "regular")
	
	// Convert to form data as Zerodha expects form-encoded data
	formData := s.convertToFormData(request)
	
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

	var response dto.ZerodhaPlaceOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// ModifyOrder modifies an existing order through Zerodha Kite API
func (s *ZerodhaOrderService) ModifyOrder(ctx context.Context, request *dto.ZerodhaModifyOrderRequest) (*dto.ZerodhaModifyOrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s/%s", s.baseURL, "regular", request.OrderID)
	
	// Convert to form data as Zerodha expects form-encoded data
	formData := s.convertModifyToFormData(request)
	
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

	var response dto.ZerodhaModifyOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CancelOrder cancels an existing order through Zerodha Kite API
func (s *ZerodhaOrderService) CancelOrder(ctx context.Context, request *dto.ZerodhaCancelOrderRequest) (*dto.ZerodhaCancelOrderResponse, error) {
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

	var response dto.ZerodhaCancelOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetOrders retrieves all orders for the client
func (s *ZerodhaOrderService) GetOrders(ctx context.Context) (*dto.ZerodhaOrderListResponse, error) {
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
func (s *ZerodhaOrderService) GetOrderHistory(ctx context.Context, orderID string) (*dto.ZerodhaOrderHistoryResponse, error) {
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
func (s *ZerodhaOrderService) GetTrades(ctx context.Context) (*dto.ZerodhaTradeListResponse, error) {
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
func (s *ZerodhaOrderService) GetPositions(ctx context.Context) (*dto.ZerodhaPositionListResponse, error) {
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

func (s *ZerodhaOrderService) convertToFormData(request *dto.ZerodhaPlaceOrderRequest) string {
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

func (s *ZerodhaOrderService) convertModifyToFormData(request *dto.ZerodhaModifyOrderRequest) string {
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