package stock

import (
	"algo/src/shared"
	"algo/src/web"

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

func (controller *Controller) GetStocks(gctx *gin.Context) {
	exchange := gctx.Query("exchange")
	segment := gctx.Query("segment")

	instruments, err := controller.service.masterRepo.GetAll(exchange, segment)
	if err != nil {
		web.InternalServerErrorResp(gctx, err.Error(), "INSTRUMENT_FETCH_ERROR")
		return
	}

	shared.SuccessResp(gctx, instruments)
}
