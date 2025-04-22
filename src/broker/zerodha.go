package broker

import (
	env "algo/src/core/config"
	"algo/src/shared"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Zerodha implements both OrderBroker and DataVendor interfaces
type Zerodha struct {
	clientID    string
	accessToken string
	baseURL     string
	client      *http.Client
	isConnect   bool
	envProv     *env.EnvProv
}

// NewZerodha creates a new instance of the Zerodha broker with configuration from environment variables
func NewZerodha(envProv *env.EnvProv) *Zerodha {
	return &Zerodha{
		clientID:    envProv.GetProp(env.ZERODHA_CLIENT_ID),
		accessToken: envProv.GetProp(env.ZERODHA_ACCESS_TOKEN),
		baseURL:     envProv.GetProp(env.ZERODHA_BASE_URL),
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		envProv: envProv,
	}
}

// Connect establishes connection with the Zerodha broker
func (z *Zerodha) Connect(ctx context.Context) error {
	// Validate that environment variables have been set
	if z.clientID == "" || z.accessToken == "" || z.baseURL == "" {
		return errors.New("missing required environment variables for Zerodha broker")
	}

	// Validate the token by making a simple API call to get user profile
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/user/profile", z.baseURL),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	z.setupHeaders(req)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Zerodha: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate with Zerodha, status code: %d", resp.StatusCode)
	}

	z.isConnect = true
	return nil
}

// setupHeaders sets up the common headers required for Zerodha API calls
func (z *Zerodha) setupHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", z.clientID, z.accessToken))
}

// PlaceOrder places a new order and returns the order ID
func (z *Zerodha) PlaceOrder(ctx context.Context, order Order) (string, error) {
	if !z.isConnect {
		return "", errors.New("not connected to Zerodha broker")
	}

	// Convert our internal order to Zerodha's order format
	zerodhaOrder := map[string]interface{}{
		"exchange":         "NSE", // Assuming NSE; might need to be determined from the instrument ID
		"tradingsymbol":    fmt.Sprintf("NSE:%d", order.InstrumentID),
		"transaction_type": convertDirectionZerodha(order.Direction),
		"quantity":         order.Quantity,
		"product":          "MIS", // Margin Intraday Square-off, can be parameterized
		"validity":         "DAY",
	}

	// Set order type and price
	if order.Type == MARKET {
		zerodhaOrder["order_type"] = "MARKET"
	} else if order.Type == LIMIT {
		zerodhaOrder["order_type"] = "LIMIT"
		if order.Price != nil {
			zerodhaOrder["price"] = *order.Price
		} else {
			return "", errors.New("price is required for limit orders")
		}
	}

	jsonData, err := json.Marshal(zerodhaOrder)
	if err != nil {
		return "", fmt.Errorf("failed to marshal order data: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/orders/regular", z.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create order request: %w", err)
	}

	z.setupHeaders(req)

	resp, err := z.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send order request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Data struct {
			OrderID string `json:"order_id"`
		} `json:"data"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode order response: %w", err)
	}

	if resp.StatusCode != http.StatusOK || response.Status != "success" {
		return "", fmt.Errorf("order placement failed, status code: %d", resp.StatusCode)
	}

	return response.Data.OrderID, nil
}

// CancelOrder cancels an existing order
func (z *Zerodha) CancelOrder(ctx context.Context, orderID string) error {
	if !z.isConnect {
		return errors.New("not connected to Zerodha broker")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf("%s/orders/regular/%s", z.baseURL, orderID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create cancel order request: %w", err)
	}

	z.setupHeaders(req)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send cancel order request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Message != "" {
			return fmt.Errorf("cancel order failed: %s", errResp.Message)
		}
		return fmt.Errorf("cancel order failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// GetOrderStatus retrieves the current status of an order
func (z *Zerodha) GetOrderStatus(ctx context.Context, orderID string) (OrderStatus, error) {
	if !z.isConnect {
		return "", errors.New("not connected to Zerodha broker")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/orders/%s", z.baseURL, orderID),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create get order status request: %w", err)
	}

	z.setupHeaders(req)

	resp, err := z.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send get order status request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get order status failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode order status response: %w", err)
	}

	return convertZerodhaStatus(response.Data.Status), nil
}

// GetPositions retrieves all current positions
func (z *Zerodha) GetPositions(ctx context.Context) ([]shared.Position, error) {
	if !z.isConnect {
		return nil, errors.New("not connected to Zerodha broker")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/portfolio/positions", z.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create get positions request: %w", err)
	}

	z.setupHeaders(req)

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get positions request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get positions failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data []struct {
			TradingSymbol string  `json:"tradingsymbol"`
			Quantity      int     `json:"quantity"`
			AveragePrice  float64 `json:"average_price"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode positions response: %w", err)
	}

	positions := make([]shared.Position, 0, len(response.Data))
	for _, pos := range response.Data {
		// Extract instrument ID from trading symbol (simplified approach)
		// In a real implementation, you might need a more sophisticated mapping
		// between Zerodha's symbols and your own instrument IDs
		var instrumentID int64
		_, err := fmt.Sscanf(pos.TradingSymbol, "NSE:%d", &instrumentID)
		if err != nil {
			continue // Skip this position if we can't parse the instrument ID
		}

		direction := shared.LONG
		if pos.Quantity < 0 {
			direction = shared.SHORT
		}

		entryPrice := pos.AveragePrice

		// Create a new position using shared.NewPosition
		position := shared.NewPosition(uuid.New(), direction)
		position.UpdateOnEntryOrder(&entryPrice, &time.Time{})

		positions = append(positions, *position)
	}

	return positions, nil
}

// GetQuote gets real-time quote for an instrument
func (z *Zerodha) GetQuote(ctx context.Context, instrumentID int32) (*shared.OHLCV, error) {
	if !z.isConnect {
		return nil, errors.New("not connected to Zerodha broker")
	}

	// For Zerodha, we need to format the instrument ID as a trading symbol
	tradingSymbol := fmt.Sprintf("NSE:%d", instrumentID)

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/market-data/quote?tradingsymbol=%s", z.baseURL, tradingSymbol),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create get quote request: %w", err)
	}

	z.setupHeaders(req)

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get quote request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get quote failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data map[string]struct {
			LastPrice float64 `json:"last_price"`
			OHLC      struct {
				Open  float64 `json:"open"`
				High  float64 `json:"high"`
				Low   float64 `json:"low"`
				Close float64 `json:"close"`
			} `json:"ohlc"`
			Volume int64 `json:"volume"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode quote response: %w", err)
	}

	// Look for our trading symbol in the response
	quoteData, found := response.Data[tradingSymbol]
	if !found {
		return nil, fmt.Errorf("no quote data found for instrument: %d", instrumentID)
	}

	return &shared.OHLCV{
		Open:   quoteData.OHLC.Open,
		High:   quoteData.OHLC.High,
		Low:    quoteData.OHLC.Low,
		Close:  quoteData.LastPrice,
		Volume: quoteData.Volume,
	}, nil
}

// Subscribe to real-time updates for instruments
func (z *Zerodha) Subscribe(ctx context.Context, instrumentIDs []int32) error {
	if !z.isConnect {
		return errors.New("not connected to Zerodha broker")
	}

	// Convert int32 IDs to strings in the format Zerodha expects
	tradingSymbols := make([]string, len(instrumentIDs))
	for i, id := range instrumentIDs {
		tradingSymbols[i] = fmt.Sprintf("NSE:%d", id)
	}

	subscriptionData := map[string]interface{}{
		"tradingsymbols": tradingSymbols,
		"mode":           "full", // full market depth
	}

	jsonData, err := json.Marshal(subscriptionData)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription data: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/market-data/subscribe", z.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription request: %w", err)
	}

	z.setupHeaders(req)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send subscription request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("subscription failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// Unsubscribe from real-time updates
func (z *Zerodha) Unsubscribe(ctx context.Context, instrumentIDs []int32) error {
	if !z.isConnect {
		return errors.New("not connected to Zerodha broker")
	}

	// Convert int32 IDs to strings in the format Zerodha expects
	tradingSymbols := make([]string, len(instrumentIDs))
	for i, id := range instrumentIDs {
		tradingSymbols[i] = fmt.Sprintf("NSE:%d", id)
	}

	unsubscriptionData := map[string]interface{}{
		"tradingsymbols": tradingSymbols,
	}

	jsonData, err := json.Marshal(unsubscriptionData)
	if err != nil {
		return fmt.Errorf("failed to marshal unsubscription data: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/market-data/unsubscribe", z.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create unsubscription request: %w", err)
	}

	z.setupHeaders(req)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send unsubscription request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsubscription failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// Helper functions

// convertDirectionZerodha converts our internal direction to Zerodha's transaction type
func convertDirectionZerodha(direction shared.Direction) string {
	if direction == shared.LONG {
		return "BUY"
	}
	return "SELL"
}

// convertZerodhaStatus converts Zerodha status to our internal order status
func convertZerodhaStatus(zerodhaStatus string) OrderStatus {
	switch zerodhaStatus {
	case "COMPLETE":
		return EXECUTED
	case "REJECTED":
		return REJECTED
	case "CANCELLED":
		return CANCELLED
	case "PENDING", "OPEN":
		return PENDING
	default:
		return PENDING
	}
}
