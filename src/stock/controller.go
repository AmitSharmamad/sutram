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
	exchange := shared.Exchange(gctx.Query("exchange"))
	segment := shared.Segment(gctx.Query("segment"))

	// Get all instruments from master repo
	instruments, err := controller.service.masterRepo.GetAll()
	if err != nil {
		web.InternalServerErrorResp(gctx, err.Error(), "INSTRUMENT_FETCH_ERROR")
		return
	}

	// Apply filters if provided
	filteredInstruments := make([]shared.InstrumentMaster, 0)
	for _, instrument := range instruments {
		if (exchange == "" || instrument.Exchange == exchange) &&
			(segment == "" || instrument.Segment == segment) {
			filteredInstruments = append(filteredInstruments, instrument)
		}
	}

	shared.SuccessResp(gctx, filteredInstruments)
}
