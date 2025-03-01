package stock

import (
	"algo/src/shared"
	"time"
)

type Service struct {
	repo *Repo
}

func NewServiceProv(repo *Repo) *Service {
	return &Service{
		repo: repo,
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
	default:
		// Default to LastMonth
		dateRange = shared.LastMonth
	}

	return dateRange
}

func (service *Service) Get(
	name string,
	dateRangeStr string,
	lastNRecords int32,
) *shared.StockData {

	instrumentId, err := service.repo.GetInstrumentId(name)
	if err != nil {
		panic(err)
	}
	dateRange := service.ConvertToDateRange(dateRangeStr)

	stockData, err := service.repo.Get(instrumentId, dateRange, lastNRecords)
	if err != nil {
		panic(err)
	}

	return stockData
}
