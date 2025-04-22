package dvendor

import (
	"algo/src/shared"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ZerodhaDataVendor implements the DataVendor interface for Zerodha
type ZerodhaDataVendor struct {
	config    VendorConfig
	client    *http.Client
	isConnect bool
}

// NewZerodhaDataVendor creates a new instance of the Zerodha data vendor
func NewZerodhaDataVendor(config VendorConfig) *ZerodhaDataVendor {
	return &ZerodhaDataVendor{
		config: config,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Connect establishes connection with the Zerodha data vendor
func (z *ZerodhaDataVendor) Connect(ctx context.Context) error {
	// Validate that configuration parameters have been set
	if z.config.ClientID == "" || z.config.Token == "" || z.config.BaseURL == "" {
		return errors.New("missing required configuration for Zerodha data vendor")
	}

	// Validate the token by making a simple API call to get user profile
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/user/profile", z.config.BaseURL),
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
func (z *ZerodhaDataVendor) setupHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", z.config.ClientID, z.config.Token))
}

// GetQuote gets real-time quote for an instrument
func (z *ZerodhaDataVendor) GetQuote(ctx context.Context, instrumentID int32) (*shared.OHLCV, error) {
	if !z.isConnect {
		return nil, errors.New("not connected to Zerodha data vendor")
	}

	// For Zerodha, we need to format the instrument ID as a trading symbol
	tradingSymbol := fmt.Sprintf("NSE:%d", instrumentID)

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/market-data/quote?tradingsymbol=%s", z.config.BaseURL, tradingSymbol),
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
func (z *ZerodhaDataVendor) Subscribe(ctx context.Context, instrumentIDs []int32) error {
	if !z.isConnect {
		return errors.New("not connected to Zerodha data vendor")
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
		fmt.Sprintf("%s/market-data/subscribe", z.config.BaseURL),
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
func (z *ZerodhaDataVendor) Unsubscribe(ctx context.Context, instrumentIDs []int32) error {
	if !z.isConnect {
		return errors.New("not connected to Zerodha data vendor")
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
		fmt.Sprintf("%s/market-data/unsubscribe", z.config.BaseURL),
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
