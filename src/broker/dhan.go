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
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Dhan implements both OrderBroker and DataVendor interfaces
type Dhan struct {
	clientID     string
	accessToken  string
	baseURL      string
	dhanClientID string
	client       *http.Client
	isConnect    bool
	envProv      *env.EnvProv
}

// NewDhan creates a new instance of the Dhan broker with configuration from environment variables
func NewDhan(envProv *env.EnvProv) *Dhan {
	return &Dhan{
		clientID:     envProv.GetProp(env.DHAN_CLIENT_ID),
		accessToken:  envProv.GetProp(env.DHAN_ACCESS_TOKEN),
		baseURL:      envProv.GetProp(env.DHAN_BASE_URL),
		dhanClientID: envProv.GetProp(env.DHAN_DHAN_CLIENT_ID),
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		envProv: envProv,
	}
}

// Connect establishes connection with the Dhan broker
func (d *Dhan) Connect(ctx context.Context) error {
	// Validate that environment variables have been set
	if d.clientID == "" || d.accessToken == "" || d.baseURL == "" || d.dhanClientID == "" {
		return errors.New("missing required environment variables for Dhan broker")
	}

	// Validate the token by making a simple API call to get user profile
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/user/profile", d.baseURL),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	d.setupHeaders(req)

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Dhan: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate with Dhan, status code: %d", resp.StatusCode)
	}

	d.isConnect = true
	return nil
}

// setupHeaders sets up the common headers required for Dhan API calls
func (d *Dhan) setupHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access-token", d.accessToken)
}

// PlaceOrder places a new order and returns the order ID
func (d *Dhan) PlaceOrder(ctx context.Context, order Order) (string, error) {
	if !d.isConnect {
		return "", errors.New("not connected to Dhan broker")
	}

	// Generate a correlation ID for this order
	correlationID := fmt.Sprintf("order_%s", time.Now().Format("20060102150405"))

	// Convert our internal order to Dhan's order format
	dhanOrder := map[string]interface{}{
		"dhanClientId":      d.dhanClientID,
		"correlationId":     correlationID,
		"transactionType":   convertDirection(order.Direction),
		"exchangeSegment":   "NSE_EQ",   // Assuming NSE Equity; might need to be determined from the instrument ID
		"productType":       "INTRADAY", // Can be parameterized: INTRADAY, DELIVERY, etc.
		"validity":          "DAY",      // Default validity: DAY
		"securityId":        strconv.FormatInt(int64(order.InstrumentID), 10),
		"quantity":          strconv.Itoa(order.Quantity),
		"disclosedQuantity": "",    // Optional, leaving empty
		"afterMarketOrder":  false, // Default false, can be parameterized
	}

	// Set order type and price
	if order.Type == MARKET {
		dhanOrder["orderType"] = "MARKET"
		dhanOrder["price"] = ""
		dhanOrder["triggerPrice"] = ""
	} else if order.Type == LIMIT {
		dhanOrder["orderType"] = "LIMIT"
		if order.Price != nil {
			dhanOrder["price"] = strconv.FormatFloat(*order.Price, 'f', 2, 64)
		} else {
			return "", errors.New("price is required for limit orders")
		}
		dhanOrder["triggerPrice"] = ""
	}

	jsonData, err := json.Marshal(dhanOrder)
	if err != nil {
		return "", fmt.Errorf("failed to marshal order data: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/orders", d.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create order request: %w", err)
	}

	d.setupHeaders(req)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send order request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		OrderID     string `json:"orderId"`
		OrderStatus string `json:"orderStatus"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode order response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("order placement failed, status code: %d", resp.StatusCode)
	}

	return response.OrderID, nil
}

// CancelOrder cancels an existing order
func (d *Dhan) CancelOrder(ctx context.Context, orderID string) error {
	if !d.isConnect {
		return errors.New("not connected to Dhan broker")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf("%s/orders/%s", d.baseURL, orderID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create cancel order request: %w", err)
	}

	d.setupHeaders(req)

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send cancel order request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
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
func (d *Dhan) GetOrderStatus(ctx context.Context, orderID string) (OrderStatus, error) {
	if !d.isConnect {
		return "", errors.New("not connected to Dhan broker")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/orders/%s", d.baseURL, orderID),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create get order status request: %w", err)
	}

	d.setupHeaders(req)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send get order status request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get order status failed with status code: %d", resp.StatusCode)
	}

	var orderResp struct {
		OrderStatus string `json:"orderStatus"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return "", fmt.Errorf("failed to decode order status response: %w", err)
	}

	return convertDhanStatus(orderResp.OrderStatus), nil
}

// GetPositions retrieves all current positions
func (d *Dhan) GetPositions(ctx context.Context) ([]shared.Position, error) {
	if !d.isConnect {
		return nil, errors.New("not connected to Dhan broker")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/portfolio/positions", d.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create get positions request: %w", err)
	}

	d.setupHeaders(req)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get positions request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get positions failed with status code: %d", resp.StatusCode)
	}

	var positionsResp []struct {
		SecurityID   string  `json:"securityId"`
		Quantity     int     `json:"quantity"`
		AveragePrice float64 `json:"averagePrice"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&positionsResp); err != nil {
		return nil, fmt.Errorf("failed to decode positions response: %w", err)
	}

	positions := make([]shared.Position, 0, len(positionsResp))
	for _, pos := range positionsResp {
		// Security ID is not currently used in the Position struct, but parsing it to validate format
		_, err := strconv.ParseInt(pos.SecurityID, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse security ID: %w", err)
		}

		direction := shared.LONG
		// Just check sign for direction, no need to modify the quantity here
		// The actual quantity is handled by the position object
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
func (d *Dhan) GetQuote(ctx context.Context, instrumentID int32) (*shared.OHLCV, error) {
	if !d.isConnect {
		return nil, errors.New("not connected to Dhan broker")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/market-quote/%d", d.baseURL, instrumentID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create get quote request: %w", err)
	}

	d.setupHeaders(req)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get quote request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get quote failed with status code: %d", resp.StatusCode)
	}

	var quoteResp struct {
		LastPrice float64 `json:"lastPrice"`
		Open      float64 `json:"open"`
		High      float64 `json:"high"`
		Low       float64 `json:"low"`
		Close     float64 `json:"close"`
		Volume    int64   `json:"volume"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&quoteResp); err != nil {
		return nil, fmt.Errorf("failed to decode quote response: %w", err)
	}

	return &shared.OHLCV{
		Open:   quoteResp.Open,
		High:   quoteResp.High,
		Low:    quoteResp.Low,
		Close:  quoteResp.LastPrice,
		Volume: quoteResp.Volume,
	}, nil
}

// Subscribe to real-time updates for instruments
func (d *Dhan) Subscribe(ctx context.Context, instrumentIDs []int32) error {
	if !d.isConnect {
		return errors.New("not connected to Dhan broker")
	}

	// Convert int32 slice to string slice for API
	instrumentIDStrs := make([]string, len(instrumentIDs))
	for i, id := range instrumentIDs {
		instrumentIDStrs[i] = strconv.FormatInt(int64(id), 10)
	}

	subscriptionData := map[string]interface{}{
		"securityIds": instrumentIDStrs,
	}

	jsonData, err := json.Marshal(subscriptionData)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription data: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/live-market-feed", d.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription request: %w", err)
	}

	d.setupHeaders(req)

	resp, err := d.client.Do(req)
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
func (d *Dhan) Unsubscribe(ctx context.Context, instrumentIDs []int32) error {
	if !d.isConnect {
		return errors.New("not connected to Dhan broker")
	}

	// Convert int32 slice to string slice for API
	instrumentIDStrs := make([]string, len(instrumentIDs))
	for i, id := range instrumentIDs {
		instrumentIDStrs[i] = strconv.FormatInt(int64(id), 10)
	}

	unsubscriptionData := map[string]interface{}{
		"securityIds": instrumentIDStrs,
	}

	jsonData, err := json.Marshal(unsubscriptionData)
	if err != nil {
		return fmt.Errorf("failed to marshal unsubscription data: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf("%s/live-market-feed", d.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create unsubscription request: %w", err)
	}

	d.setupHeaders(req)

	resp, err := d.client.Do(req)
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

// convertDirection converts our internal direction to Dhan's transaction type
func convertDirection(direction shared.Direction) string {
	if direction == shared.LONG {
		return "BUY"
	}
	return "SELL"
}

// convertDhanStatus converts Dhan status to our internal order status
func convertDhanStatus(dhanStatus string) OrderStatus {
	switch dhanStatus {
	case "TRADED":
		return EXECUTED
	case "REJECTED":
		return REJECTED
	case "CANCELLED":
		return CANCELLED
	case "PENDING", "TRANSIT":
		return PENDING
	default:
		return PENDING
	}
}
