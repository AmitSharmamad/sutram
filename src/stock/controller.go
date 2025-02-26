package stock

import (
	"algo/src/shared"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewControllerProv(
	service *Service,
) *Controller {
	return &Controller{
		service: service,
	}
}

func (controller *Controller) Get(gctx *gin.Context) {
	instrumentName := gctx.Query("instrument")

	if instrumentName == "" {
		panic("instrument name is required")
	}

	stockData := controller.service.Get(instrumentName, 600)
	out := controller.service.Converter(*stockData)

	shared.SuccessResp(gctx, out)
}
