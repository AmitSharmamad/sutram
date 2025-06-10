package http

import (
	"net/http"

	"algo/internal/modules/broker/application/dto"
	"algo/internal/modules/broker/application/service"

	"github.com/gin-gonic/gin"
)

// BrokerController handles HTTP requests for broker operations
type BrokerController struct {
	dhanBroker *service.DhanBroker
}

// NewBrokerController creates a new broker controller
func NewBrokerController(dhanBroker *service.DhanBroker) *BrokerController {
	return &BrokerController{
		dhanBroker: dhanBroker,
	}
}

// PlaceOrder handles placing a new order through Dhan
// @Summary Place order through Dhan
// @Description Place a new order through Dhan broker API
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dto.DhanPlaceOrderRequest true "Order details"
// @Success 200 {object} dto.DhanPlaceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/broker/dhan/orders [post]
func (c *BrokerController) PlaceOrder(ctx *gin.Context) {
	var request dto.DhanPlaceOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.dhanBroker.PlaceDhanOrder(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ModifyOrder handles modifying an existing order through Dhan
// @Summary Modify order through Dhan
// @Description Modify an existing order through Dhan broker API
// @Tags orders
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Param order body dto.DhanModifyOrderRequest true "Order modification details"
// @Success 200 {object} dto.DhanModifyOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/broker/dhan/orders/{orderID} [put]
func (c *BrokerController) ModifyOrder(ctx *gin.Context) {
	var request dto.DhanModifyOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the order ID from URL matches the request
	orderID := ctx.Param("orderID")
	request.OrderID = orderID

	response, err := c.dhanBroker.ModifyDhanOrder(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// CancelOrder handles canceling an order through Dhan
// @Summary Cancel order through Dhan
// @Description Cancel an existing order through Dhan broker API
// @Tags orders
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Param request body dto.DhanCancelOrderRequest true "Cancel order request"
// @Success 200 {object} dto.DhanCancelOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/broker/dhan/orders/{orderID} [delete]
func (c *BrokerController) CancelOrder(ctx *gin.Context) {
	var request dto.DhanCancelOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the order ID from URL matches the request
	orderID := ctx.Param("orderID")
	request.OrderID = orderID

	response, err := c.dhanBroker.CancelDhanOrder(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetOrders handles getting all orders from Dhan
// @Summary Get all orders from Dhan
// @Description Retrieve all orders for the client from Dhan broker API
// @Tags orders
// @Produce json
// @Success 200 {object} dto.DhanOrderListResponse
// @Failure 500 {object} map[string]string
// @Router /api/broker/dhan/orders [get]
func (c *BrokerController) GetOrders(ctx *gin.Context) {
	response, err := c.dhanBroker.GetOrders(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetOrderByID handles getting a specific order by ID from Dhan
// @Summary Get order by ID from Dhan
// @Description Retrieve a specific order by ID from Dhan broker API
// @Tags orders
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} dto.DhanOrderByIDResponse
// @Failure 500 {object} map[string]string
// @Router /api/broker/dhan/orders/{orderID} [get]
func (c *BrokerController) GetOrderByID(ctx *gin.Context) {
	orderID := ctx.Param("orderID")
	
	response, err := c.dhanBroker.GetOrderByID(ctx.Request.Context(), orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTrades handles getting all trades from Dhan
// @Summary Get all trades from Dhan
// @Description Retrieve all trades for the client from Dhan broker API
// @Tags trades
// @Produce json
// @Success 200 {object} dto.DhanTradeListResponse
// @Failure 500 {object} map[string]string
// @Router /api/broker/dhan/trades [get]
func (c *BrokerController) GetTrades(ctx *gin.Context) {
	response, err := c.dhanBroker.GetTrades(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetOrderBook handles getting the order book from Dhan
// @Summary Get order book from Dhan
// @Description Retrieve the order book for the client from Dhan broker API
// @Tags orders
// @Produce json
// @Success 200 {object} dto.DhanOrderListResponse
// @Failure 500 {object} map[string]string
// @Router /api/broker/dhan/orderbook [get]
func (c *BrokerController) GetOrderBook(ctx *gin.Context) {
	response, err := c.dhanBroker.GetOrderBook(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTradeBook handles getting the trade book from Dhan
// @Summary Get trade book from Dhan
// @Description Retrieve the trade book for the client from Dhan broker API
// @Tags trades
// @Produce json
// @Success 200 {object} dto.DhanTradeListResponse
// @Failure 500 {object} map[string]string
// @Router /api/broker/dhan/tradebook [get]
func (c *BrokerController) GetTradeBook(ctx *gin.Context) {
	response, err := c.dhanBroker.GetTradeBook(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}