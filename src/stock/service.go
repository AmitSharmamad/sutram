package stock

import (
	"fmt"
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
