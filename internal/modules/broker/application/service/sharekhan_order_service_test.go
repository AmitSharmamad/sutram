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

func TestSharekhanBroker_PlaceOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.SharekhanPlaceOrderResponse{
			OrderNumber:     "SK12345",
			OrderType:       "MARKET",
			TradingSymbol:   "RELIANCE",
			Exchange:        "NSE",
			TransactionType: "BUY",
			ProductCode:     "CNC",
			Quantity:        10,
			Price:           2500.00,
			Validity:        "DAY",
			OrderStatus:     "PENDING",
			OrderTime:       time.Now().Format(time.RFC3339),
			Message:         "Order placed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewSharekhanBroker(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.SharekhanPlaceOrderRequest{
		TradingSymbol:   "RELIANCE",
		Exchange:        "NSE",
		TransactionType: "BUY",
		OrderType:       "MARKET",
		ProductCode:     "CNC",
		Quantity:        10,
		Price:           2500.00,
		Validity:        "DAY",
	}

	response, err := service.PlaceSharekhanOrder(context.Background(), request)

	require.NoError(t, err)
	assert.Equal(t, "SK12345", response.OrderNumber)
	assert.Equal(t, "PENDING", response.OrderStatus)
	assert.Equal(t, "BUY", response.TransactionType)
	assert.Equal(t, "RELIANCE", response.TradingSymbol)
}

func TestSharekhanBroker_ModifyOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/orders/SK12345", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.SharekhanModifyOrderResponse{
			OrderNumber:     "SK12345",
			OrderType:       "LIMIT",
			TradingSymbol:   "RELIANCE",
			Exchange:        "NSE",
			TransactionType: "BUY",
			ProductCode:     "CNC",
			Quantity:        15,
			Price:           2400.00,
			Validity:        "DAY",
			OrderStatus:     "MODIFIED",
			ModifiedTime:    time.Now().Format(time.RFC3339),
			Message:         "Order modified successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewSharekhanBroker(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.SharekhanModifyOrderRequest{
		OrderNumber: "SK12345",
		OrderType:   "LIMIT",
		Quantity:    15,
		Price:       2400.00,
		Validity:    "DAY",
	}

	response, err := service.ModifySharekhanOrder(context.Background(), request)

	require.NoError(t, err)
	assert.Equal(t, "SK12345", response.OrderNumber)
	assert.Equal(t, "MODIFIED", response.OrderStatus)
	assert.Equal(t, 15, response.Quantity)
	assert.Equal(t, 2400.00, response.Price)
}

func TestSharekhanBroker_CancelOrder(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/orders/SK12345", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := dto.SharekhanCancelOrderResponse{
			OrderNumber:   "SK12345",
			OrderStatus:   "CANCELLED",
			CancelledTime: time.Now().Format(time.RFC3339),
			Message:       "Order cancelled successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewSharekhanBroker(server.URL, "test-api-key", "test-token")

	// Test request
	request := &dto.SharekhanCancelOrderRequest{
		OrderNumber: "SK12345",
	}

	response, err := service.CancelSharekhanOrder(context.Background(), request)

	require.NoError(t, err)
	assert.Equal(t, "SK12345", response.OrderNumber)
	assert.Equal(t, "CANCELLED", response.OrderStatus)
}

func TestSharekhanBroker_GetOrders(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		orders := []dto.SharekhanOrderResponse{
			{
				OrderNumber:     "SK12345",
				TradingSymbol:   "RELIANCE",
				Exchange:        "NSE",
				TransactionType: "BUY",
				OrderType:       "MARKET",
				ProductCode:     "CNC",
				Quantity:        10,
				Price:           2500.00,
				OrderStatus:     "COMPLETE",
				Validity:        "DAY",
				OrderTime:       time.Now(),
			},
		}

		response := dto.SharekhanOrderListResponse{
			Orders:  orders,
			Status:  "SUCCESS",
			Message: "Orders retrieved successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewSharekhanBroker(server.URL, "test-api-key", "test-token")

	response, err := service.GetOrders(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Len(t, response.Orders, 1)
	assert.Equal(t, "SK12345", response.Orders[0].OrderNumber)
	assert.Equal(t, "RELIANCE", response.Orders[0].TradingSymbol)
}

func TestSharekhanBroker_GetOrderByID(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orders/SK12345", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		order := dto.SharekhanOrderResponse{
			OrderNumber:     "SK12345",
			TradingSymbol:   "RELIANCE",
			Exchange:        "NSE",
			TransactionType: "BUY",
			OrderType:       "MARKET",
			ProductCode:     "CNC",
			Quantity:        10,
			Price:           2500.00,
			OrderStatus:     "COMPLETE",
			Validity:        "DAY",
			OrderTime:       time.Now(),
		}

		response := dto.SharekhanOrderByIDResponse{
			Order:   order,
			Status:  "SUCCESS",
			Message: "Order retrieved successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewSharekhanBroker(server.URL, "test-api-key", "test-token")

	response, err := service.GetOrderByID(context.Background(), "SK12345")

	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Equal(t, "SK12345", response.Order.OrderNumber)
	assert.Equal(t, "RELIANCE", response.Order.TradingSymbol)
}

func TestSharekhanBroker_GetTrades(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/trades", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		trades := []dto.SharekhanTradeResponse{
			{
				TradeNumber:     "T12345",
				OrderNumber:     "SK12345",
				TradingSymbol:   "RELIANCE",
				Exchange:        "NSE",
				TransactionType: "BUY",
				ProductCode:     "CNC",
				Quantity:        10,
				Price:           2500.00,
				TradeValue:      25000.00,
				TradeTime:       time.Now(),
				Brokerage:       10.00,
				NetValue:        25010.00,
			},
		}

		response := dto.SharekhanTradeListResponse{
			Trades:  trades,
			Status:  "SUCCESS",
			Message: "Trades retrieved successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewSharekhanBroker(server.URL, "test-api-key", "test-token")

	response, err := service.GetTrades(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Len(t, response.Trades, 1)
	assert.Equal(t, "T12345", response.Trades[0].TradeNumber)
	assert.Equal(t, "SK12345", response.Trades[0].OrderNumber)
}

func TestSharekhanBroker_GetOrderBook(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/orderbook", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		orders := []dto.SharekhanOrderResponse{
			{
				OrderNumber:     "SK12345",
				TradingSymbol:   "RELIANCE",
				Exchange:        "NSE",
				TransactionType: "BUY",
				OrderType:       "MARKET",
				ProductCode:     "CNC",
				Quantity:        10,
				Price:           2500.00,
				OrderStatus:     "COMPLETE",
				Validity:        "DAY",
				OrderTime:       time.Now(),
			},
		}

		response := dto.SharekhanOrderBookResponse{
			OrderBook: orders,
			Status:    "SUCCESS",
			Message:   "Order book retrieved successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewSharekhanBroker(server.URL, "test-api-key", "test-token")

	response, err := service.GetOrderBook(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Len(t, response.OrderBook, 1)
	assert.Equal(t, "SK12345", response.OrderBook[0].OrderNumber)
}

func TestSharekhanBroker_GetTradeBook(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/tradebook", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		trades := []dto.SharekhanTradeResponse{
			{
				TradeNumber:     "T12345",
				OrderNumber:     "SK12345",
				TradingSymbol:   "RELIANCE",
				Exchange:        "NSE",
				TransactionType: "BUY",
				ProductCode:     "CNC",
				Quantity:        10,
				Price:           2500.00,
				TradeValue:      25000.00,
				TradeTime:       time.Now(),
				Brokerage:       10.00,
				NetValue:        25010.00,
			},
		}

		response := dto.SharekhanTradeBookResponse{
			TradeBook: trades,
			Status:    "SUCCESS",
			Message:   "Trade book retrieved successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create service with mock server
	service := NewSharekhanBroker(server.URL, "test-api-key", "test-token")

	response, err := service.GetTradeBook(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Len(t, response.TradeBook, 1)
	assert.Equal(t, "T12345", response.TradeBook[0].TradeNumber)
}