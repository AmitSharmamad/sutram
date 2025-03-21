package handler

import (
	"algo/src/shared"
)

// DBHandler implements OHLCVHandler interface for database operations
type DBHandler struct {
	stockService shared.StockService
}

func NewDBHandler(
	stockService shared.StockService,
) *DBHandler {
	return &DBHandler{
		stockService: stockService,
	}
}

func (h *DBHandler) HandleBatch(
	master shared.InstrumentMaster,
	data []*shared.OHLCV,
) error {
	h.stockService.Insert(master, data)
	return nil
}
