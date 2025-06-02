package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"algo/internal/modules/broker/application/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAngelOneOrderService_PlaceOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/rest/secure/angelbroking/order/v1/placeOrder", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-PrivateKey"))

		response := dto.AngelOnePlaceOrderResponse{
			Status:    true,
			Message:   "Order placed successfully",
			ErrorCode: "",
			Data: struct {
				Script   string `json:"script"`
				OrderID  string `json:"orderid"`
				UniqueID string `json:"uniqueid"`
			}{
				Script:   "RELIANCE-EQ",
				OrderID:  "221114000000372",
				UniqueID: "unique123",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewAngelOneOrderService(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.AngelOnePlaceOrderRequest{
		Variety:         "NORMAL",
		TradingSymbol:   "RELIANCE-EQ",
		SymbolToken:     "2885",
		TransactionType: "BUY",
		Exchange:        "NSE",
		OrderType:       "MARKET",
		ProductType:     "DELIVERY",
		Duration:        "DAY",
		Quantity:        "10",
	}

	// Execute test
	response, err := service.PlaceOrder(context.Background(), request)

	// Assertions
	require.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "Order placed successfully", response.Message)
	assert.Equal(t, "221114000000372", response.Data.OrderID)
	assert.Equal(t, "RELIANCE-EQ", response.Data.Script)
}

func TestAngelOneOrderService_PlaceOrder_Error(t *testing.T) {
	// Mock HTTP server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := dto.AngelOnePlaceOrderResponse{
			Status:    false,
			Message:   "Insufficient balance",
			ErrorCode: "AG8001",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewAngelOneOrderService(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.AngelOnePlaceOrderRequest{
		Variety:         "NORMAL",
		TradingSymbol:   "RELIANCE-EQ",
		SymbolToken:     "2885",
		TransactionType: "BUY",
		Exchange:        "NSE",
		OrderType:       "MARKET",
		ProductType:     "DELIVERY",
		Duration:        "DAY",
		Quantity:        "10",
	}

	// Execute test
	_, err := service.PlaceOrder(context.Background(), request)

	// Assertions
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Insufficient balance")
	assert.Contains(t, err.Error(), "AG8001")
}

func TestAngelOneOrderService_ModifyOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/rest/secure/angelbroking/order/v1/modifyOrder", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := dto.AngelOneModifyOrderResponse{
			Status:    true,
			Message:   "Order modified successfully",
			ErrorCode: "",
			Data: struct {
				OrderID  string `json:"orderid"`
				UniqueID string `json:"uniqueid"`
			}{
				OrderID:  "221114000000372",
				UniqueID: "unique123",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewAngelOneOrderService(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.AngelOneModifyOrderRequest{
		Variety:     "NORMAL",
		OrderID:     "221114000000372",
		OrderType:   "LIMIT",
		ProductType: "DELIVERY",
		Duration:    "DAY",
		Price:       "2600",
		Quantity:    "5",
	}

	// Execute test
	response, err := service.ModifyOrder(context.Background(), request)

	// Assertions
	require.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "Order modified successfully", response.Message)
	assert.Equal(t, "221114000000372", response.Data.OrderID)
}

func TestAngelOneOrderService_CancelOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/rest/secure/angelbroking/order/v1/cancelOrder", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := dto.AngelOneCancelOrderResponse{
			Status:    true,
			Message:   "Order cancelled successfully",
			ErrorCode: "",
			Data: struct {
				OrderID  string `json:"orderid"`
				UniqueID string `json:"uniqueid"`
			}{
				OrderID:  "221114000000372",
				UniqueID: "unique123",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewAngelOneOrderService(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.AngelOneCancelOrderRequest{
		Variety: "NORMAL",
		OrderID: "221114000000372",
	}

	// Execute test
	response, err := service.CancelOrder(context.Background(), request)

	// Assertions
	require.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "Order cancelled successfully", response.Message)
	assert.Equal(t, "221114000000372", response.Data.OrderID)
}

func TestAngelOneOrderService_GetOrders(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/rest/secure/angelbroking/order/v1/getOrderBook", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := dto.AngelOneOrderListResponse{
			Status:    true,
			Message:   "SUCCESS",
			ErrorCode: "",
			Data: []dto.AngelOneOrderResponse{
				{
					OrderID:         "221114000000372",
					TradingSymbol:   "RELIANCE-EQ",
					Exchange:        "NSE",
					ProductType:     "DELIVERY",
					OrderType:       "MARKET",
					TransactionType: "BUY",
					Quantity:        "10",
					Status:          "COMPLETE",
					Price:           "2500",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewAngelOneOrderService(server.URL, "test-api-key", "test-token")

	// Execute test
	response, err := service.GetOrders(context.Background())

	// Assertions
	require.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "221114000000372", response.Data[0].OrderID)
	assert.Equal(t, "RELIANCE-EQ", response.Data[0].TradingSymbol)
}

func TestAngelOneOrderService_GetOrderByID(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/rest/secure/angelbroking/order/v1/details/221114000000372", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := dto.AngelOneOrderByIDResponse{
			Status:    true,
			Message:   "SUCCESS",
			ErrorCode: "",
			Data: dto.AngelOneOrderResponse{
				OrderID:         "221114000000372",
				TradingSymbol:   "RELIANCE-EQ",
				Exchange:        "NSE",
				ProductType:     "DELIVERY",
				OrderType:       "MARKET",
				TransactionType: "BUY",
				Quantity:        "10",
				Status:          "COMPLETE",
				Price:           "2500",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewAngelOneOrderService(server.URL, "test-api-key", "test-token")

	// Execute test
	response, err := service.GetOrderByID(context.Background(), "221114000000372")

	// Assertions
	require.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, "221114000000372", response.Data.OrderID)
	assert.Equal(t, "RELIANCE-EQ", response.Data.TradingSymbol)
}

func TestAngelOneOrderService_GetTrades(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/rest/secure/angelbroking/order/v1/getTradeBook", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := dto.AngelOneTradeListResponse{
			Status:    true,
			Message:   "SUCCESS",
			ErrorCode: "",
			Data: []dto.AngelOneTradeResponse{
				{
					TradeID:         "123456789",
					OrderID:         "221114000000372",
					TradingSymbol:   "RELIANCE-EQ",
					Exchange:        "NSE",
					ProductType:     "DELIVERY",
					TransactionType: "BUY",
					TradedQuantity:  "10",
					TradedPrice:     "2500",
					TradedValue:     "25000",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewAngelOneOrderService(server.URL, "test-api-key", "test-token")

	// Execute test
	response, err := service.GetTrades(context.Background())

	// Assertions
	require.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "123456789", response.Data[0].TradeID)
	assert.Equal(t, "221114000000372", response.Data[0].OrderID)
}