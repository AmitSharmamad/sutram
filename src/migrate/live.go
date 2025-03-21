package migrate

import "algo/src/shared"

type OHLCVSingleHandler interface {
	HandleSingle(*shared.OHLCV) error
}

// Handle processes a single data point using the provided data provider
func Handle(data map[string]string, provider DataProvider, handler OHLCVSingleHandler) error {
	record, err := provider.MapToOHLCV(data)
	if err != nil {
		return err
	}

	return handler.HandleSingle(record)
}
