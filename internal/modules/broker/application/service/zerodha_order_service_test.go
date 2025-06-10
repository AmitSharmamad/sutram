package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"algo/internal/modules/broker/application/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZerodhaBroker_PlaceOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/orders/regular", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Equal(t, "3", r.Header.Get("X-Kite-Version"))
		assert.Equal(t, "token test-api-key:test-access-token", r.Header.Get("Authorization"))

		// Check form data contains expected fields
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		formData := string(body)
		assert.Contains(t, formData, "tradingsymbol=RELIANCE")
		assert.Contains(t, formData, "exchange=NSE")
		assert.Contains(t, formData, "transaction_type=BUY")
		assert.Contains(t, formData, "quantity=10")

		response := dto.ZerodhaPlaceOrderResponse{
			OrderID: "240308000000001",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewZerodhaBroker(server.URL, "test-api-key", "test-access-token")

	// Test request
	request := &dto.ZerodhaPlaceOrderRequest{
		TradingSymbol:   "RELIANCE",
		Exchange:        "NSE",
		TransactionType: "BUY",
		OrderType:       "LIMIT",
		Quantity:        10,
		Price:           2500.00,
		Product:         "CNC",
		Validity:        "DAY",
	}

	// Execute test
	ctx := context.Background()
	response, err := service.PlaceZerodhaOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "240308000000001", response.OrderID)
}

func TestZerodhaBroker_ModifyOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/orders/regular/240308000000001", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Equal(t, "3", r.Header.Get("X-Kite-Version"))
		assert.Equal(t, "token test-api-key:test-access-token", r.Header.Get("Authorization"))

		response := dto.ZerodhaModifyOrderResponse{
			OrderID: "240308000000001",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewZerodhaBroker(server.URL, "test-api-key", "test-access-token")

	// Test request
	request := &dto.ZerodhaModifyOrderRequest{
		OrderID:  "240308000000001",
		Quantity: 15,
		Price:    2600.00,
	}

	// Execute test
	ctx := context.Background()
	response, err := service.ModifyZerodhaOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "240308000000001", response.OrderID)
}

func TestZerodhaBroker_CancelOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/orders/regular/240308000000001", r.URL.Path)
		assert.Equal(t, "3", r.Header.Get("X-Kite-Version"))
		assert.Equal(t, "token test-api-key:test-access-token", r.Header.Get("Authorization"))

		response := dto.ZerodhaCancelOrderResponse{
			OrderID: "240308000000001",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewZerodhaBroker(server.URL, "test-api-key", "test-access-token")

	// Test request
	request := &dto.ZerodhaCancelOrderRequest{
		OrderID: "240308000000001",
	}

	// Execute test
	ctx := context.Background()
	response, err := service.CancelZerodhaOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "240308000000001", response.OrderID)
}

func TestZerodhaBroker_GetOrders(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "3", r.Header.Get("X-Kite-Version"))
		assert.Equal(t, "token test-api-key:test-access-token", r.Header.Get("Authorization"))

		response := dto.ZerodhaOrderListResponse{
			Status: "success",
			Data: []dto.ZerodhaOrderResponse{
				{
					OrderID:         "240308000000001",
					Status:          "COMPLETE",
					TradingSymbol:   "RELIANCE",
					Exchange:        "NSE",
					TransactionType: "BUY",
					OrderType:       "LIMIT",
					Product:         "CNC",
					Quantity:        10,
					Price:           2500.00,
					AveragePrice:    2499.50,
					FilledQuantity:  10,
					PendingQuantity: 0,
					OrderTimestamp:  time.Now(),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewZerodhaBroker(server.URL, "test-api-key", "test-access-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetOrders(ctx)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "240308000000001", response.Data[0].OrderID)
	assert.Equal(t, "COMPLETE", response.Data[0].Status)
	assert.Equal(t, "RELIANCE", response.Data[0].TradingSymbol)
}

func TestZerodhaBroker_GetOrderHistory(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders/240308000000001", r.URL.Path)
		assert.Equal(t, "3", r.Header.Get("X-Kite-Version"))
		assert.Equal(t, "token test-api-key:test-access-token", r.Header.Get("Authorization"))

		response := dto.ZerodhaOrderHistoryResponse{
			Status: "success",
			Data: []dto.ZerodhaOrderResponse{
				{
					OrderID:         "240308000000001",
					Status:          "COMPLETE",
					TradingSymbol:   "RELIANCE",
					Exchange:        "NSE",
					TransactionType: "BUY",
					OrderType:       "LIMIT",
					Product:         "CNC",
					Quantity:        10,
					Price:           2500.00,
					AveragePrice:    2499.50,
					FilledQuantity:  10,
					PendingQuantity: 0,
					OrderTimestamp:  time.Now(),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewZerodhaBroker(server.URL, "test-api-key", "test-access-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetOrderHistory(ctx, "240308000000001")

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "240308000000001", response.Data[0].OrderID)
	assert.Equal(t, "COMPLETE", response.Data[0].Status)
}

func TestZerodhaBroker_GetTrades(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/trades", r.URL.Path)
		assert.Equal(t, "3", r.Header.Get("X-Kite-Version"))
		assert.Equal(t, "token test-api-key:test-access-token", r.Header.Get("Authorization"))

		response := dto.ZerodhaTradeListResponse{
			Status: "success",
			Data: []dto.ZerodhaTradeResponse{
				{
					TradeID:           "12345",
					OrderID:           "240308000000001",
					TradingSymbol:     "RELIANCE",
					Exchange:          "NSE",
					TransactionType:   "BUY",
					Product:           "CNC",
					AveragePrice:      2499.50,
					Quantity:          10,
					FillTimestamp:     time.Now(),
					ExchangeTimestamp: time.Now(),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewZerodhaBroker(server.URL, "test-api-key", "test-access-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetTrades(ctx)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "12345", response.Data[0].TradeID)
	assert.Equal(t, "240308000000001", response.Data[0].OrderID)
	assert.Equal(t, "RELIANCE", response.Data[0].TradingSymbol)
}

func TestZerodhaBroker_GetPositions(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/portfolio/positions", r.URL.Path)
		assert.Equal(t, "3", r.Header.Get("X-Kite-Version"))
		assert.Equal(t, "token test-api-key:test-access-token", r.Header.Get("Authorization"))

		response := dto.ZerodhaPositionListResponse{
			Status: "success",
			Data: struct {
				Net []dto.ZerodhaPositionResponse `json:"net"`
				Day []dto.ZerodhaPositionResponse `json:"day"`
			}{
				Net: []dto.ZerodhaPositionResponse{
					{
						TradingSymbol: "RELIANCE",
						Exchange:      "NSE",
						Product:       "CNC",
						Quantity:      10,
						AveragePrice:  2499.50,
						LastPrice:     2520.00,
						PNL:           205.00,
					},
				},
				Day: []dto.ZerodhaPositionResponse{},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewZerodhaBroker(server.URL, "test-api-key", "test-access-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetPositions(ctx)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Len(t, response.Data.Net, 1)
	assert.Equal(t, "RELIANCE", response.Data.Net[0].TradingSymbol)
	assert.Equal(t, "NSE", response.Data.Net[0].Exchange)
	assert.Equal(t, 10, response.Data.Net[0].Quantity)
}

func TestZerodhaBroker_ConvertToFormData(t *testing.T) {
	service := NewZerodhaBroker("http://test.com", "test-api-key", "test-access-token")

	request := &dto.ZerodhaPlaceOrderRequest{
		TradingSymbol:   "RELIANCE",
		Exchange:        "NSE",
		TransactionType: "BUY",
		OrderType:       "LIMIT",
		Quantity:        10,
		Price:           2500.00,
		Product:         "CNC",
		Validity:        "DAY",
		Tag:             "test-tag",
	}

	formData := service.convertToFormData(request)

	// Check that all required fields are present
	assert.Contains(t, formData, "tradingsymbol=RELIANCE")
	assert.Contains(t, formData, "exchange=NSE")
	assert.Contains(t, formData, "transaction_type=BUY")
	assert.Contains(t, formData, "order_type=LIMIT")
	assert.Contains(t, formData, "quantity=10")
	assert.Contains(t, formData, "price=2500.00")
	assert.Contains(t, formData, "product=CNC")
	assert.Contains(t, formData, "validity=DAY")
	assert.Contains(t, formData, "tag=test-tag")

	// Check format (should be URL encoded)
	assert.True(t, strings.Contains(formData, "="))
	assert.True(t, strings.Contains(formData, "&"))
}