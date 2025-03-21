package stock

import (
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
	name string,
	dateRangeStr string,
	granularityStr string,
	lastNRecords int32,
) *shared.StockData {
	instrumentId, err := service.repo.GetInstrumentId(name)
	if err != nil {
		panic(err)
	}
	dateRange := service.ConvertToDateRange(dateRangeStr)
	granularity := service.ConvertToGranularity(granularityStr)

	stockData, err := service.repo.Get(instrumentId, dateRange, granularity, lastNRecords)
	if err != nil {
		panic(err)
	}

	return stockData
}

func (s *Service) Insert(master shared.InstrumentMaster, instruments []*shared.OHLCV) error {
	// Insert the master record
	latestID, err := s.masterRepo.GetInstrument(master.Symbol)
	if err != nil {
		println(err)
	}

	master.ID = latestID + 1

	s.masterRepo.Insert(&master)

	s.repo.Insert(instruments)

	return nil
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

func (service *InstrumentService) InsertMaster(
	master *shared.InstrumentMaster,
) error {
	// master.ID = util.NewUUID()
	// master.InstrumentID = 1
	service.masterRepo.Insert(master)
	println("added to master")
	return nil
}

func (service *InstrumentService) InsertData(
	instruments []shared.OHLCV,
) error {
	// service.repo.Insert(instruments)
	return nil
}
