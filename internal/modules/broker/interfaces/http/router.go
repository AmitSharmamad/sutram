package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the broker routes
func RegisterRoutes(router *gin.RouterGroup, controller *BrokerController) {
	// Dhan order API routes
	dhanGroup := router.Group("/dhan")
	{
		// Order management routes
		dhanGroup.POST("/orders", controller.PlaceOrder)
		dhanGroup.GET("/orders", controller.GetOrders)
		dhanGroup.GET("/orders/:orderID", controller.GetOrderByID)
		dhanGroup.PUT("/orders/:orderID", controller.ModifyOrder)
		dhanGroup.DELETE("/orders/:orderID", controller.CancelOrder)
		
		// Trade routes
		dhanGroup.GET("/trades", controller.GetTrades)
		
		// Order book and trade book routes
		dhanGroup.GET("/orderbook", controller.GetOrderBook)
		dhanGroup.GET("/tradebook", controller.GetTradeBook)
	}
}