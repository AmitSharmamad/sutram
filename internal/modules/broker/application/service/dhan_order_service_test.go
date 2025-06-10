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

func TestDhanOrderService_PlaceOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-token", r.Header.Get("access-token"))

		response := dto.DhanPlaceOrderResponse{
			OrderID:         "12345",
			DrmOrderID:      "DRM12345",
			OrderStatus:     "PENDING",
			TransactionType: "BUY",
			ExchangeSegment: "NSE_EQ",
			ProductType:     "CNC",
			OrderType:       "MARKET",
			Validity:        "DAY",
			TradingSymbol:   "RELIANCE",
			SecurityID:      "2885",
			Quantity:        10,
			Price:           2500.00,
			CreateTime:      time.Now().Format(time.RFC3339),
			UpdateTime:      time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewDhanOrderService(server.URL, "test-token")

	// Test request
	request := &dto.DhanPlaceOrderRequest{
		DhanClientID:    "TEST123",
		TransactionType: "BUY",
		ExchangeSegment: "NSE_EQ",
		ProductType:     "CNC",
		OrderType:       "MARKET",
		Validity:        "DAY",
		SecurityID:      "2885",
		Quantity:        10,
		Price:           2500.00,
	}

	// Execute test
	ctx := context.Background()
	response, err := service.PlaceDhanOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "12345", response.OrderID)
	assert.Equal(t, "PENDING", response.OrderStatus)
	assert.Equal(t, "BUY", response.TransactionType)
	assert.Equal(t, 10, response.Quantity)
	assert.Equal(t, 2500.00, response.Price)
}

func TestDhanOrderService_GetOrders(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "test-token", r.Header.Get("access-token"))

		response := dto.DhanOrderListResponse{
			Orders: []dto.DhanOrderStatusResponse{
				{
					OrderID:         "12345",
					OrderStatus:     "EXECUTED",
					TransactionType: "BUY",
					ExchangeSegment: "NSE_EQ",
					ProductType:     "CNC",
					OrderType:       "MARKET",
					Quantity:        10,
					Price:           2500.00,
					FilledQty:       10,
					AveragePrice:    2498.50,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewDhanOrderService(server.URL, "test-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetOrders(ctx)

	// Assertions
	require.NoError(t, err)
	require.Len(t, response.Orders, 1)
	assert.Equal(t, "12345", response.Orders[0].OrderID)
	assert.Equal(t, "EXECUTED", response.Orders[0].OrderStatus)
	assert.Equal(t, 10, response.Orders[0].FilledQty)
}

func TestDhanOrderService_CancelOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/orders/12345", r.URL.Path)
		assert.Equal(t, "test-token", r.Header.Get("access-token"))

		response := dto.DhanCancelOrderResponse{
			OrderID:     "12345",
			OrderStatus: "CANCELLED",
			UpdateTime:  time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewDhanOrderService(server.URL, "test-token")

	// Test request
	request := &dto.DhanCancelOrderRequest{
		DhanClientID: "TEST123",
		OrderID:      "12345",
	}

	// Execute test
	ctx := context.Background()
	response, err := service.CancelDhanOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "12345", response.OrderID)
	assert.Equal(t, "CANCELLED", response.OrderStatus)
}