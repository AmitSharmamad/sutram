package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the instrument module routes
func RegisterRoutes(router *gin.RouterGroup, controller *InstrumentController) {
	// Instrument routes
	instruments := router.Group("/instruments")
	{
		instruments.GET("/:id", controller.GetInstrumentByID)
		instruments.PUT("/:id", controller.UpdateInstrument)
		instruments.DELETE("/:id", controller.DeleteInstrument)
		instruments.POST("", controller.CreateInstrument)
		instruments.GET("/symbol/:symbol", controller.GetInstrumentsBySymbol)
		instruments.POST("/timerange", controller.GetInstrumentsByTimeRange)
	}

	// Instrument master routes
	masters := router.Group("/instrument-masters")
	{
		masters.GET("", controller.GetAllInstrumentMasters)
		masters.GET("/:id", controller.GetInstrumentMasterByID)
		masters.PUT("/:id", controller.UpdateInstrumentMaster)
		masters.DELETE("/:id", controller.DeleteInstrumentMaster)
		masters.POST("", controller.CreateInstrumentMaster)
	}
}
