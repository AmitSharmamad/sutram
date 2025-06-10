package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"algo/internal/modules/broker/application/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDhanOrderService_BrokerInterface_PlaceOrder(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/orders", r.URL.Path)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"orderId": "12345",
			"orderStatus": "PENDING",
			"transactionType": "BUY",
			"quantity": 10,
			"price": 2500.00
		}`))
	}))
	defer server.Close()

	// Create service instance
	service := NewDhanBroker(server.URL, "test-token")

	// Prepare request using common DTO
	request := &dto.PlaceOrderRequest{
		ID:              uuid.New(),
		UserID:          uuid.New(),
		PositionID:      uuid.New(),
		Exchange:        "NSE",
		Segment:         "EQ",
		Instrument:      "RELIANCE",
		TransactionType: "BUY",
		Quantity:        10,
		Price:           2500.00,
		ProductType:     "CNC",
		OrderType:       "LIMIT",
	}

	// Execute test
	ctx := context.Background()
	response, err := service.PlaceOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, request.ID, response.ID)
	assert.Equal(t, request.UserID, response.UserID)
	assert.Equal(t, request.PositionID, response.PositionID)
	assert.Equal(t, "PENDING", response.State)
	assert.Equal(t, "12345", response.BrokerOrderID)
	assert.Equal(t, "BUY", response.TransactionType)
	assert.Equal(t, 10, response.Quantity)
	assert.Equal(t, 2500.00, response.Price)
}

func TestDhanOrderService_BrokerInterface_UpdateOrder(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/orders/12345", r.URL.Path)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"orderId": "12345",
			"orderStatus": "MODIFIED",
			"quantity": 20,
			"price": 2600.00
		}`))
	}))
	defer server.Close()

	// Create service instance
	service := NewDhanBroker(server.URL, "test-token")

	// Prepare request using common DTO
	newQuantity := 20
	newPrice := 2600.00
	request := &dto.ModifyOrderRequest{
		OrderID:  "12345",
		Quantity: &newQuantity,
		Price:    &newPrice,
	}

	// Execute test
	ctx := context.Background()
	response, err := service.UpdateOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "12345", response.BrokerOrderID)
	assert.Equal(t, "MODIFIED", response.State)
}

func TestDhanOrderService_BrokerInterface_CancelOrder(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/orders/12345", r.URL.Path)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"orderId": "12345",
			"orderStatus": "CANCELLED"
		}`))
	}))
	defer server.Close()

	// Create service instance
	service := NewDhanBroker(server.URL, "test-token")

	// Prepare request using common DTO
	request := &dto.CancelOrderRequest{
		OrderID: "12345",
	}

	// Execute test
	ctx := context.Background()
	response, err := service.CancelOrder(ctx, request)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "12345", response.BrokerOrderID)
	assert.Equal(t, "CANCELLED", response.State)
}