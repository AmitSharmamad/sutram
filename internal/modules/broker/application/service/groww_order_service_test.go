package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"algo/internal/modules/broker/application/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrowwOrderService_PlaceOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.GrowwPlaceOrderResponse{
			OrderID:     "GRW12345",
			Symbol:      "RELIANCE",
			Exchange:    "NSE",
			Segment:     "EQT",
			Side:        "BUY",
			OrderType:   "MARKET",
			ProductType: "DELIVERY",
			Validity:    "DAY",
			Quantity:    10,
			Price:       2500.00,
			Status:      "PENDING",
			Message:     "Order placed successfully",
			OrderTime:   time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewGrowwOrderService(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.GrowwPlaceOrderRequest{
		Symbol:      "RELIANCE",
		Exchange:    "NSE",
		Segment:     "EQT",
		Side:        "BUY",
		OrderType:   "MARKET",
		ProductType: "DELIVERY",
		Validity:    "DAY",
		Quantity:    10,
		Price:       2500.00,
	}

	response, err := service.PlaceGrowwOrder(context.Background(), request)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "GRW12345", response.OrderID)
	assert.Equal(t, "RELIANCE", response.Symbol)
	assert.Equal(t, "PENDING", response.Status)
}

func TestGrowwOrderService_ModifyOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/orders/GRW12345", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.GrowwModifyOrderResponse{
			OrderID:   "GRW12345",
			Status:    "MODIFIED",
			Message:   "Order modified successfully",
			UpdatedAt: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewGrowwOrderService(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.GrowwModifyOrderRequest{
		OrderID:   "GRW12345",
		OrderType: "LIMIT",
		Quantity:  15,
		Price:     2600.00,
		Validity:  "DAY",
	}

	response, err := service.ModifyGrowwOrder(context.Background(), request)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "GRW12345", response.OrderID)
	assert.Equal(t, "MODIFIED", response.Status)
}

func TestGrowwOrderService_CancelOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/orders/GRW12345", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.GrowwCancelOrderResponse{
			OrderID:   "GRW12345",
			Status:    "CANCELLED",
			Message:   "Order cancelled successfully",
			UpdatedAt: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewGrowwOrderService(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.GrowwCancelOrderRequest{
		OrderID: "GRW12345",
	}

	response, err := service.CancelGrowwOrder(context.Background(), request)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "GRW12345", response.OrderID)
	assert.Equal(t, "CANCELLED", response.Status)
}

func TestGrowwOrderService_GetOrders(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.GrowwOrderListResponse{
			Orders: []dto.GrowwOrderStatusResponse{
				{
					OrderID:     "GRW12345",
					Symbol:      "RELIANCE",
					Exchange:    "NSE",
					Side:        "BUY",
					Status:      "COMPLETE",
					Quantity:    10,
					FilledQty:   10,
					Price:       2500.00,
					OrderTime:   time.Now(),
				},
			},
			Status:  "SUCCESS",
			Message: "Orders retrieved successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewGrowwOrderService(server.URL, "test-api-key", "test-token")

	response, err := service.GetOrders(context.Background())

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Len(t, response.Orders, 1)
	assert.Equal(t, "GRW12345", response.Orders[0].OrderID)
}

func TestGrowwOrderService_GetOrderByID(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders/GRW12345", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.GrowwOrderByIDResponse{
			Order: dto.GrowwOrderStatusResponse{
				OrderID:   "GRW12345",
				Symbol:    "RELIANCE",
				Exchange:  "NSE",
				Side:      "BUY",
				Status:    "COMPLETE",
				Quantity:  10,
				FilledQty: 10,
				Price:     2500.00,
				OrderTime: time.Now(),
			},
			Status:  "SUCCESS",
			Message: "Order retrieved successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewGrowwOrderService(server.URL, "test-api-key", "test-token")

	response, err := service.GetOrderByID(context.Background(), "GRW12345")

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Equal(t, "GRW12345", response.Order.OrderID)
}

func TestGrowwOrderService_GetTrades(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/trades", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.GrowwTradeListResponse{
			Trades: []dto.GrowwTradeResponse{
				{
					TradeID:   "TRD12345",
					OrderID:   "GRW12345",
					Symbol:    "RELIANCE",
					Exchange:  "NSE",
					Side:      "BUY",
					Quantity:  10,
					Price:     2500.00,
					TradeTime: time.Now(),
				},
			},
			Status:  "SUCCESS",
			Message: "Trades retrieved successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewGrowwOrderService(server.URL, "test-api-key", "test-token")

	response, err := service.GetTrades(context.Background())

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Len(t, response.Trades, 1)
	assert.Equal(t, "TRD12345", response.Trades[0].TradeID)
}