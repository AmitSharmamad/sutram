package stock

import (
	"algo/src/shared"
	"algo/src/web"
	"strconv"

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

	// Get pagination parameters with defaults
	page := 1
	limit := 10
	if p := gctx.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if l := gctx.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	instruments, total, err := controller.service.masterRepo.GetAll(exchange, segment, page, limit)
	if err != nil {
		web.InternalServerErrorResp(gctx, err.Error(), "INSTRUMENT_FETCH_ERROR")
		return
	}

	shared.PaginatedResp(gctx, instruments, total, page, limit)
}
