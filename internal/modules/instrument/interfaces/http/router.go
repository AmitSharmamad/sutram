package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the instrument module routes
func RegisterRoutes(router *gin.RouterGroup, controller *InstrumentController) {
	// Instrument routes
	instruments := router.Group("/instrument")
	{
		instruments.GET("/:id", controller.GetInstrumentByExchangeSecurityID)
		instruments.PUT("/:id", controller.UpdateInstrument)
		instruments.DELETE("/:id", controller.DeleteInstrument)
		instruments.POST("", controller.CreateInstrument)
		instruments.POST("/timerange", controller.GetInstrumentsByTimeRange)
	}

	// Instrument master routes
	masters := router.Group("/instrument-master")
	{
		masters.GET("", controller.GetAllInstrumentMasters)
		masters.GET("/by-exchange-security-ids", controller.GetInstrumentMastersByExchangeSecurityIDs)
		masters.GET("/by-exchange-security-id/:exchange_security_id", controller.GetInstrumentMasterByExchangeSecurityID)
		masters.PUT("/:exchange_security_id", controller.UpdateInstrumentMaster)
		masters.DELETE("/:exchange_security_id", controller.DeleteInstrumentMaster)
		masters.POST("", controller.CreateInstrumentMaster)
	}
}
