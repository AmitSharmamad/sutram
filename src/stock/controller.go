package stock

import (
	"algo/src/shared"
	"strconv"
	"strings"

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
		panic(shared.NewBadRequest(err.Error(), "INSTRUMENT_FETCH_ERROR"))
	}

	shared.PaginatedResp(gctx, instruments, total, page, limit)
}

func (controller *Controller) GetStocksByIDs(gctx *gin.Context) {
	// Get the IDs from query parameter
	idsStr := gctx.Query("ids")
	if idsStr == "" {
		panic(shared.NewBadRequest("No IDs provided", "MISSING_IDS"))
	}

	// Split the comma-separated IDs and convert to uint64
	idStrs := strings.Split(idsStr, ",")
	ids := make([]uint64, 0, len(idStrs))

	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			panic(shared.NewBadRequest("Invalid ID format: "+idStr, "INVALID_ID"))
		}
		ids = append(ids, id)
	}

	// Get the stocks
	stocks, err := controller.service.GetStocksByIDs(ids)
	if err != nil {
		panic(shared.NewBadRequest(err.Error(), "STOCKS_FETCH_ERROR"))
	}

	shared.SuccessResp(gctx, stocks)
}
