package stock

import (
	"fmt"
	"math"
	"time"

	"algo/src/shared"
)

type Service struct {
	masterRepo *MasterRepo
	repo       *Repo
}

func NewServiceProv(
	masterRepo *MasterRepo,
	repo *Repo,
) *Service {
	return &Service{
		masterRepo: masterRepo,
		repo:       repo,
	}
}

func (service *Service) Converter(sd shared.StockData) *shared.StockColumn {
	// If no data, return empty column structure
	if len(sd) == 0 {
		return &shared.StockColumn{
			Timestamp: []time.Time{},
			Open:      []float64{},
			High:      []float64{},
			Close:     []float64{},
			Low:       []float64{},
			Volume:    []int64{},
			CreatedAt: []time.Time{},
		}
	}

	// Get the number of records
	n := len(sd)

	// Create the column structure
	columns := &shared.StockColumn{
		Timestamp: make([]time.Time, n),
		Open:      make([]float64, n),
		High:      make([]float64, n),
		Close:     make([]float64, n),
		Low:       make([]float64, n),
		Volume:    make([]int64, n),
		CreatedAt: make([]time.Time, n),
	}

	// Fill each column with values from the StockData
	for i, stock := range sd {
		columns.Timestamp[i] = stock.Timestamp
		columns.Open[i] = stock.Open
		columns.High[i] = stock.High
		columns.Close[i] = stock.Close
		columns.Low[i] = stock.Low
		columns.Volume[i] = stock.Volume
		columns.CreatedAt[i] = stock.CreatedAt
	}

	return columns
}

func (service *Service) ConvertToDateRange(dateRangeStr string) shared.DateRange {
	var dateRange shared.DateRange

	switch dateRangeStr {
	case "last_day":
		dateRange = shared.LastDay
	case "last_week":
		dateRange = shared.LastWeek
	case "last_month":
		dateRange = shared.LastMonth
	case "last_6_months":
		dateRange = shared.Last6Months
	case "last_year":
		dateRange = shared.LastYear
	case "year_to_date":
		dateRange = shared.YearToDate
	case "last_3_years":
		dateRange = shared.Last3Years
	case "last_5_years":
		dateRange = shared.Last5Years
	case "":
		dateRange = shared.All
	default:
		// Default to LastMonth
		dateRange = shared.LastMonth
	}

	return dateRange
}

func (service *Service) ConvertToGranularity(granularityStr string) shared.StockGranularity {
	var granularity shared.StockGranularity

	switch granularityStr {
	case "1m":
		granularity = shared.StockGranularity1m
	case "3m":
		granularity = shared.StockGranularity3m
	case "5m":
		granularity = shared.StockGranularity5m
	case "15m":
		granularity = shared.StockGranularity15m
	case "30m":
		granularity = shared.StockGranularity30m
	case "1h":
		granularity = shared.StockGranularity1h
	case "3h":
		granularity = shared.StockGranularity3h
	case "1d":
		granularity = shared.StockGranularity1d
	case "1M":
		granularity = shared.StockGranularity1M
	default:
		granularity = shared.StockGranularity1m
	}

	return granularity
}

func (service *Service) Get(
	id int32,
	startPeriod shared.CustomTime,
	endPeriod shared.CustomTime,
	timeFrame string,
	lookBack int32,
) *shared.StockData {
	dateRangeData, err := service.repo.Get(id, startPeriod, endPeriod, timeFrame)
	if err != nil {
		// Since interface doesn't allow returning error, we return empty StockData
		return &shared.StockData{}
	}

	lookBackData, err := service.repo.GetLookBack(id, timeFrame, startPeriod, lookBack)
	if err != nil {
		// Since interface doesn't allow returning error, we return date range data only
		return dateRangeData
	}

	for i, j := 0, len(*lookBackData)-1; i < j; i, j = i+1, j-1 {
		(*lookBackData)[i], (*lookBackData)[j] = (*lookBackData)[j], (*lookBackData)[i]
	}

	data := append(make(shared.StockData, 0), *lookBackData...)
	data = append(data, *dateRangeData...)

	return &data
}

func (s *Service) Insert(master shared.InstrumentMaster, instruments []*shared.OHLCV) error {
	uniqueInstruments := make(map[string]struct{})

	for _, instrument := range instruments {
		uniqueInstruments[instrument.Instrument] = struct{}{}
	}

	for instrument := range uniqueInstruments {
		m, err := s.masterRepo.GetInstrument(instrument)
		if err != nil {
			fmt.Printf("error in getting instrument %+v\n", err)
		}
		if m == nil {
			latestID, err := s.masterRepo.GetLatestID()
			if err != nil {
				panic("error in getting instrument id")
			}
			id := latestID + 1
			s.masterRepo.Insert(&shared.InstrumentMaster{
				ID:             id,
				Exchange:       master.Exchange,
				Segment:        master.Segment,
				Symbol:         instrument,
				Instrument:     instrument,
				InstrumentID:   id,
				InstrumentName: instrument,
			})
		}
	}

	instrumentIDMap := make(map[string]uint64)
	masterInstruments, total, _ := s.masterRepo.GetAll("", "", 1, 1000) // Get first 1000 instruments

	// If we have more than 1000 instruments, fetch the rest
	if total > 1000 {
		remainingPages := (total + 999) / 1000
		for page := 2; page <= remainingPages; page++ {
			moreInstruments, _, _ := s.masterRepo.GetAll("", "", page, 1000)
			masterInstruments = append(masterInstruments, moreInstruments...)
		}
	}

	for _, masterInstrument := range masterInstruments {
		if (masterInstrument.Instrument == "") || (masterInstrument.ID == 0) {
			println(masterInstrument.Instrument, masterInstrument.ID)
			panic("data is missing")
		}
		instrumentIDMap[masterInstrument.Instrument] = masterInstrument.ID
	}

	for _, instrument := range instruments {
		instrument.ID = int64(instrumentIDMap[instrument.Instrument])
	}

	err := s.repo.Insert(instruments)
	if err != nil {
		fmt.Printf("error in inserting instrument %+v\n", err)
	}

	return nil
}

func (s *Service) GetStocksByIDs(ids []uint64) ([]shared.InstrumentMaster, error) {
	return s.masterRepo.GetByIDs(ids)
}

type InstrumentService struct {
	masterRepo *MasterRepo
	repo       *Repo
}

func NewInstrumentServiceProv(
	masterRepo *MasterRepo,
	repo *Repo,
) *InstrumentService {
	return &InstrumentService{
		masterRepo: masterRepo,
		repo:       repo,
	}
}

// TODO: this is not tested at all
// CalculateHeikenAshi calculates Heiken-Ashi candles from standard OHLCV data.
// It requires the input data to be sorted chronologically.
func CalculateHeikenAshi(ohlcvData shared.StockData) (shared.StockData, error) {
	if len(ohlcvData) == 0 {
		return nil, fmt.Errorf("input ohlcvData slice is empty")
	}

	haCandles := make(shared.StockData, len(ohlcvData))

	// Calculate the first Heiken-Ashi candle
	firstOhlc := ohlcvData[0]
	haCandles[0].Timestamp = firstOhlc.Timestamp
	haCandles[0].Close = (firstOhlc.Open + firstOhlc.High + firstOhlc.Low + firstOhlc.Close) / 4
	// For the very first candle, HA Open is often calculated like this:
	haCandles[0].Open = (firstOhlc.Open + firstOhlc.Close) / 2
	haCandles[0].High = max(firstOhlc.High, max(haCandles[0].Open, haCandles[0].Close))
	haCandles[0].Low = min(firstOhlc.Low, min(haCandles[0].Open, haCandles[0].Close))

	// Calculate subsequent Heiken-Ashi candles
	for i := 1; i < len(ohlcvData); i++ {
		currentOhlc := ohlcvData[i]
		prevHA := haCandles[i-1]

		haCandles[i].Timestamp = currentOhlc.Timestamp
		haCandles[i].Close = (currentOhlc.Open + currentOhlc.High + currentOhlc.Low + currentOhlc.Close) / 4
		haCandles[i].Open = (prevHA.Open + prevHA.Close) / 2
		haCandles[i].High = max(currentOhlc.High, max(haCandles[i].Open, haCandles[i].Close))
		haCandles[i].Low = min(currentOhlc.Low, min(haCandles[i].Open, haCandles[i].Close))
	}

	return haCandles, nil
}

// TODO: this is not tested at all
// CalculateRenko calculates Renko bricks from standard OHLCV data using closing prices.
// It requires the input data to be sorted chronologically.
// brickSize determines the minimum price movement required to form a new brick.
func CalculateRenko(ohlcvData []shared.Stock, brickSize float64) ([]shared.Stock, error) {
	if len(ohlcvData) == 0 {
		return nil, fmt.Errorf("input ohlcvData slice is empty")
	}
	if brickSize <= 0 {
		return nil, fmt.Errorf("brickSize must be positive")
	}

	var renkoBricks []shared.Stock
	if len(ohlcvData) == 0 {
		return renkoBricks, nil // Return empty slice if no input data
	}

	// --- Initial Brick Determination ---
	// We need a starting point. Often, the first brick is established
	// based on the first significant move or simply the first close.
	// Here, we'll establish the first "level" based on the first close,
	// but won't draw a brick until a brickSize move occurs from it.
	// A more robust approach might analyze the first few candles.
	// For simplicity, let's set the initial level based on the first close.
	lastBrickCloseLevel := ohlcvData[0].Close
	// Or, alternatively, start with an actual first brick:
	firstBrick := shared.Stock{
		Timestamp: ohlcvData[0].Timestamp,
		Open:      ohlcvData[0].Close, // Or Open price
		Close:     ohlcvData[0].Close, // Initial close defines the level
		Direction: 0,                  // Undefined until next move
	}
	renkoBricks = append(renkoBricks, firstBrick)
	lastBrickCloseLevel = firstBrick.Close

	// --- Iterate through remaining data ---
	for i := 1; i < len(ohlcvData); i++ {
		currentPrice := ohlcvData[i].Close
		currentTimestamp := ohlcvData[i].Timestamp

		// Calculate potential next brick levels
		priceMovement := currentPrice - lastBrickCloseLevel
		direction := 0
		numBricks := 0

		if priceMovement >= brickSize {
			direction = 1 // Up brick
			numBricks = int(math.Floor(priceMovement / brickSize))
		} else if priceMovement <= -brickSize {
			direction = -1 // Down brick
			numBricks = int(math.Floor(math.Abs(priceMovement) / brickSize))
		}

		// Add new bricks if movement threshold is met
		if numBricks > 0 {
			for j := 0; j < numBricks; j++ {
				brickOpen := lastBrickCloseLevel + float64(j*direction)*brickSize
				brickClose := brickOpen + float64(direction)*brickSize
				newBrick := shared.Stock{
					Timestamp: currentTimestamp, // Mark with timestamp when threshold was crossed
					Open:      brickOpen,
					Close:     brickClose,
					Direction: direction,
				}
				renkoBricks = append(renkoBricks, newBrick)
				lastBrickCloseLevel = brickClose // Update the level for the next comparison
			}
		}
		// If numBricks is 0, the price didn't move enough to form a new brick,
		// so we keep the lastBrickCloseLevel as is and continue.
	}

	// Handle the case where the first brick wasn't immediately added
	// (if using the 'level' approach instead of adding the first brick right away)
	if len(renkoBricks) == 0 && len(ohlcvData) > 0 {
		// Optionally add an initial brick based on the first close if no movement occurred
		// This behavior might need adjustment based on specific requirements.
		initialBrick := shared.Stock{
			Timestamp: ohlcvData[0].Timestamp,
			Open:      ohlcvData[0].Close, // Or Open
			Close:     ohlcvData[0].Close,
			Direction: 0, // No direction determined yet
		}
		// Decide if you want to add this placeholder or not.
		renkoBricks = append(renkoBricks, initialBrick)
	}

	return renkoBricks, nil
}
