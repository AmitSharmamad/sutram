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

func TestHDFCSecOrderService_PlaceOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := dto.HDFCSecPlaceOrderResponse{
			OrderID:         "12345",
			OrderStatus:     "PENDING",
			TransactionType: "BUY",
			Exchange:        "NSE",
			Segment:         "EQ",
			ProductType:     "CNC",
			OrderType:       "MARKET",
			Validity:        "DAY",
			Symbol:          "RELIANCE",
			Quantity:        10,
			Price:           2500.00,
			OrderTime:       time.Now().Format(time.RFC3339),
			UpdateTime:      time.Now().Format(time.RFC3339),
			Message:         "Order placed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewHDFCSecOrderService(server.URL, "test-token")

	// Test request
	request := &dto.HDFCSecPlaceOrderRequest{
		ClientCode:      "TEST123",
		TransactionType: "BUY",
		Exchange:        "NSE",
		Segment:         "EQ",
		ProductType:     "CNC",
		OrderType:       "MARKET",
		Validity:        "DAY",
		Symbol:          "RELIANCE",
		Quantity:        10,
		Price:           2500.00,
	}

	// Execute test
	ctx := context.Background()
	response, err := service.PlaceOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "12345", response.OrderID)
	assert.Equal(t, "PENDING", response.OrderStatus)
	assert.Equal(t, "BUY", response.TransactionType)
	assert.Equal(t, 10, response.Quantity)
	assert.Equal(t, 2500.00, response.Price)
	assert.Equal(t, "Order placed successfully", response.Message)
}

func TestHDFCSecOrderService_GetOrders(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		orders := []dto.HDFCSecOrderStatusResponse{
			{
				OrderID:         "12345",
				OrderStatus:     "COMPLETED",
				TransactionType: "BUY",
				Exchange:        "NSE",
				Segment:         "EQ",
				ProductType:     "CNC",
				Symbol:          "RELIANCE",
				Quantity:        10,
				FilledQuantity:  10,
				Price:           2500.00,
				AveragePrice:    2500.00,
				OrderTime:       time.Now().Format(time.RFC3339),
			},
		}

		response := dto.HDFCSecOrderListResponse{
			Orders: orders,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewHDFCSecOrderService(server.URL, "test-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetOrders(ctx)

	// Assertions
	require.NoError(t, err)
	assert.Len(t, response.Orders, 1)
	assert.Equal(t, "12345", response.Orders[0].OrderID)
	assert.Equal(t, "COMPLETED", response.Orders[0].OrderStatus)
	assert.Equal(t, 10, response.Orders[0].Quantity)
	assert.Equal(t, 10, response.Orders[0].FilledQuantity)
}

func TestHDFCSecOrderService_CancelOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/orders/12345", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := dto.HDFCSecCancelOrderResponse{
			OrderID:     "12345",
			OrderStatus: "CANCELLED",
			UpdateTime:  time.Now().Format(time.RFC3339),
			Message:     "Order cancelled successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewHDFCSecOrderService(server.URL, "test-token")

	// Test request
	request := &dto.HDFCSecCancelOrderRequest{
		ClientCode: "TEST123",
		OrderID:    "12345",
	}

	// Execute test
	ctx := context.Background()
	response, err := service.CancelOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "12345", response.OrderID)
	assert.Equal(t, "CANCELLED", response.OrderStatus)
	assert.Equal(t, "Order cancelled successfully", response.Message)
}

func TestHDFCSecOrderService_ModifyOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/orders/12345", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		response := dto.HDFCSecModifyOrderResponse{
			OrderID:     "12345",
			OrderStatus: "MODIFIED",
			Quantity:    15,
			Price:       2600.00,
			Validity:    "DAY",
			UpdateTime:  time.Now().Format(time.RFC3339),
			Message:     "Order modified successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewHDFCSecOrderService(server.URL, "test-token")

	// Test request
	request := &dto.HDFCSecModifyOrderRequest{
		ClientCode: "TEST123",
		OrderID:    "12345",
		OrderType:  "LIMIT",
		Quantity:   15,
		Price:      2600.00,
		Validity:   "DAY",
	}

	// Execute test
	ctx := context.Background()
	response, err := service.ModifyOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "12345", response.OrderID)
	assert.Equal(t, "MODIFIED", response.OrderStatus)
	assert.Equal(t, 15, response.Quantity)
	assert.Equal(t, 2600.00, response.Price)
	assert.Equal(t, "Order modified successfully", response.Message)
}

func TestHDFCSecOrderService_GetOrderByID(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders/12345", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		order := dto.HDFCSecOrderStatusResponse{
			OrderID:         "12345",
			OrderStatus:     "COMPLETED",
			TransactionType: "BUY",
			Exchange:        "NSE",
			Segment:         "EQ",
			ProductType:     "CNC",
			Symbol:          "RELIANCE",
			Quantity:        10,
			FilledQuantity:  10,
			Price:           2500.00,
			AveragePrice:    2500.00,
			OrderTime:       time.Now().Format(time.RFC3339),
		}

		response := dto.HDFCSecOrderByIDResponse{
			Order: order,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewHDFCSecOrderService(server.URL, "test-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetOrderByID(ctx, "12345")

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "12345", response.Order.OrderID)
	assert.Equal(t, "COMPLETED", response.Order.OrderStatus)
	assert.Equal(t, 10, response.Order.Quantity)
	assert.Equal(t, 10, response.Order.FilledQuantity)
}

func TestHDFCSecOrderService_GetTrades(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/trades", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		trades := []dto.HDFCSecTradeResponse{
			{
				ClientCode:      "TEST123",
				OrderID:         "12345",
				TradeID:         "T123",
				TransactionType: "BUY",
				Exchange:        "NSE",
				Segment:         "EQ",
				ProductType:     "CNC",
				Symbol:          "RELIANCE",
				TradeQuantity:   10,
				TradePrice:      2500.00,
				TradeValue:      25000.00,
				TradeTime:       time.Now().Format(time.RFC3339),
				BrokerageAmount: 10.00,
				NetAmount:       24990.00,
			},
		}

		response := dto.HDFCSecTradeListResponse{
			Trades: trades,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewHDFCSecOrderService(server.URL, "test-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetTrades(ctx)

	// Assertions
	require.NoError(t, err)
	assert.Len(t, response.Trades, 1)
	assert.Equal(t, "T123", response.Trades[0].TradeID)
	assert.Equal(t, "12345", response.Trades[0].OrderID)
	assert.Equal(t, 10, response.Trades[0].TradeQuantity)
	assert.Equal(t, 2500.00, response.Trades[0].TradePrice)
}

func TestHDFCSecOrderService_GetPositions(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/positions", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		positions := []dto.HDFCSecPositionResponse{
			{
				ClientCode:       "TEST123",
				Exchange:         "NSE",
				Segment:          "EQ",
				ProductType:      "CNC",
				Symbol:           "RELIANCE",
				NetQuantity:      10,
				BuyQuantity:      10,
				SellQuantity:     0,
				AverageBuyPrice:  2500.00,
				AverageSellPrice: 0.00,
				NetAmount:        25000.00,
				RealizedPL:       0.00,
				UnrealizedPL:     100.00,
				MTM:              100.00,
				LTP:              2510.00,
				LastTradeTime:    time.Now().Format(time.RFC3339),
			},
		}

		response := dto.HDFCSecPositionListResponse{
			Positions: positions,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewHDFCSecOrderService(server.URL, "test-token")

	// Execute test
	ctx := context.Background()
	response, err := service.GetPositions(ctx)

	// Assertions
	require.NoError(t, err)
	assert.Len(t, response.Positions, 1)
	assert.Equal(t, "RELIANCE", response.Positions[0].Symbol)
	assert.Equal(t, 10, response.Positions[0].NetQuantity)
	assert.Equal(t, 2500.00, response.Positions[0].AverageBuyPrice)
	assert.Equal(t, 100.00, response.Positions[0].UnrealizedPL)
}