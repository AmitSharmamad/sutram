package http

import (
	"net/http"

	"algo/internal/modules/broker/application/dto"
	"algo/internal/modules/broker/domain"

	"github.com/gin-gonic/gin"
)

// UnifiedBrokerController demonstrates how the controller could be updated to use the common Broker interface
// This is an example implementation showing how to leverage the new common interface
type UnifiedBrokerController struct {
	// Map of broker name to broker service implementing the common interface
	brokers map[string]domain.Broker
}

// NewUnifiedBrokerController creates a new unified broker controller
func NewUnifiedBrokerController(brokers map[string]domain.Broker) *UnifiedBrokerController {
	return &UnifiedBrokerController{
		brokers: brokers,
	}
}

// PlaceOrder handles placing a new order through any broker using the common interface
// @Summary Place order through specified broker
// @Description Place a new order through any supported broker using the common interface
// @Tags orders
// @Accept json
// @Produce json
// @Param broker path string true "Broker name (dhan, zerodha, sharekhan, hdfcsec, angelone, groww)"
// @Param order body dto.PlaceOrderRequest true "Order details using common format"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/broker/{broker}/orders [post]
func (c *UnifiedBrokerController) PlaceOrder(ctx *gin.Context) {
	brokerName := ctx.Param("broker")
	
	broker, exists := c.brokers[brokerName]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Broker not supported: " + brokerName})
		return
	}

	var request dto.PlaceOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := broker.PlaceOrder(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateOrder handles modifying an existing order through any broker using the common interface
// @Summary Modify order through specified broker
// @Description Modify an existing order through any supported broker using the common interface
// @Tags orders
// @Accept json
// @Produce json
// @Param broker path string true "Broker name (dhan, zerodha, sharekhan, hdfcsec, angelone, groww)"
// @Param order body dto.ModifyOrderRequest true "Order modification details"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/broker/{broker}/orders [put]
func (c *UnifiedBrokerController) UpdateOrder(ctx *gin.Context) {
	brokerName := ctx.Param("broker")
	
	broker, exists := c.brokers[brokerName]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Broker not supported: " + brokerName})
		return
	}

	var request dto.ModifyOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := broker.UpdateOrder(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// CancelOrder handles canceling an order through any broker using the common interface
// @Summary Cancel order through specified broker
// @Description Cancel an existing order through any supported broker using the common interface
// @Tags orders
// @Accept json
// @Produce json
// @Param broker path string true "Broker name (dhan, zerodha, sharekhan, hdfcsec, angelone, groww)"
// @Param order body dto.CancelOrderRequest true "Order cancellation details"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/broker/{broker}/orders [delete]
func (c *UnifiedBrokerController) CancelOrder(ctx *gin.Context) {
	brokerName := ctx.Param("broker")
	
	broker, exists := c.brokers[brokerName]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Broker not supported: " + brokerName})
		return
	}

	var request dto.CancelOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := broker.CancelOrder(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}