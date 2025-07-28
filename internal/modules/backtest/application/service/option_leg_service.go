package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"algo/internal/modules/backtest/application/dto"
	"algo/internal/modules/backtest/domain/model"
	"algo/internal/modules/backtest/domain/repository"
	"algo/internal/modules/instrument/application/service"
	instrumentModel "algo/internal/modules/instrument/domain/model"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrOptionLegNotFound       = errors.New("option leg not found")
	ErrInvalidOptionLegRequest = errors.New("invalid option leg request")
	ErrUnauthorizedOptionLeg   = errors.New("unauthorized access to option leg")
	ErrOptionLegCreateFailed   = errors.New("failed to create option leg")
	ErrOptionLegUpdateFailed   = errors.New("failed to update option leg")
	ErrOptionLegDeleteFailed   = errors.New("failed to delete option leg")
	ErrDuplicateOrder          = errors.New("leg sequence already exists for this strategy")
)

// OptionLegService defines the service for option leg operations
type OptionLegService interface {
	// CreateOptionLeg creates a new option leg
	CreateOptionLeg(ctx context.Context, strategyID uuid.UUID, req *dto.CreateOptionLegRequest) (*dto.OptionLegResponse, error)

	// GetOptionLegByID retrieves an option leg by its ID
	GetOptionLegByID(ctx context.Context, id uuid.UUID) (*dto.OptionLegResponse, error)

	// GetOptionLegsByStrategyID retrieves all option legs for a strategy
	GetOptionLegsByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*dto.OptionLegResponse, error)

	// UpdateOptionLeg updates an existing option leg
	UpdateOptionLeg(ctx context.Context, id uuid.UUID, req *dto.UpdateOptionLegRequest) (*dto.OptionLegResponse, error)

	// DeleteOptionLeg deletes an option leg
	DeleteOptionLeg(ctx context.Context, id uuid.UUID) error

	// DeleteOptionLegsByStrategyID deletes all option legs for a strategy
	DeleteOptionLegsByStrategyID(ctx context.Context, strategyID uuid.UUID) error

	// GetOptionLegsByIndex retrieves option legs by index/underlying
	GetOptionLegsByIndex(ctx context.Context, selectIndex string) ([]*dto.OptionLegResponse, error)

	// SearchOptionLegs searches option legs with filter criteria
	SearchOptionLegs(ctx context.Context, filter *dto.OptionLegFilter) ([]*dto.OptionLegResponse, int64, error)

	// ValidateOptionLeg validates option leg data
	ValidateOptionLeg(ctx context.Context, req *dto.CreateOptionLegRequest) error

	// ResolveExchangeSecurityID resolves the exchange_security_id for an option leg based on spot price
	ResolveExchangeSecurityID(ctx context.Context, legID uuid.UUID, spotPrice float64) (int32, error)
}

// optionLegService implements the OptionLegService interface
type optionLegService struct {
	optionLegRepo     repository.OptionLegRepository
	instrumentService service.InstrumentService
}

// NewOptionLegService creates a new option leg service
func NewOptionLegService(optionLegRepo repository.OptionLegRepository, instrumentService service.InstrumentService) OptionLegService {
	return &optionLegService{
		optionLegRepo:     optionLegRepo,
		instrumentService: instrumentService,
	}
}

// CreateOptionLeg creates a new option leg
func (s *optionLegService) CreateOptionLeg(ctx context.Context, strategyID uuid.UUID, req *dto.CreateOptionLegRequest) (*dto.OptionLegResponse, error) {
	// Validate the request
	if err := s.ValidateOptionLeg(ctx, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if leg sequence already exists for the strategy
	existingLeg, err := s.optionLegRepo.GetByStrategyIDAndSequence(ctx, strategyID, req.LegSequence)
	if err == nil && existingLeg != nil {
		return nil, ErrDuplicateOrder
	}

	// Convert DTO to domain model
	optionLeg := req.ToModel(strategyID)

	// Calculate total quantity if lot size is provided
	if optionLeg.LotSize != nil {
		totalQuantity := optionLeg.TotalLot * *optionLeg.LotSize
		optionLeg.TotalQuantity = &totalQuantity
	}

	// Save to database
	if err := s.optionLegRepo.Create(ctx, optionLeg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOptionLegCreateFailed, err)
	}

	// Return response DTO
	return dto.ToResponse(optionLeg), nil
}

// GetOptionLegByID retrieves an option leg by its ID
func (s *optionLegService) GetOptionLegByID(ctx context.Context, id uuid.UUID) (*dto.OptionLegResponse, error) {
	optionLeg, err := s.optionLegRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get option leg: %w", err)
	}

	if optionLeg == nil {
		return nil, ErrOptionLegNotFound
	}

	return dto.ToResponse(optionLeg), nil
}

// GetOptionLegsByStrategyID retrieves all option legs for a strategy
func (s *optionLegService) GetOptionLegsByStrategyID(ctx context.Context, strategyID uuid.UUID) ([]*dto.OptionLegResponse, error) {
	optionLegs, err := s.optionLegRepo.GetByStrategyID(ctx, strategyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get option legs for strategy: %w", err)
	}

	responses := make([]*dto.OptionLegResponse, len(optionLegs))
	for i, optionLeg := range optionLegs {
		responses[i] = dto.ToResponse(optionLeg)
	}

	return responses, nil
}

// UpdateOptionLeg updates an existing option leg
func (s *optionLegService) UpdateOptionLeg(ctx context.Context, id uuid.UUID, req *dto.UpdateOptionLegRequest) (*dto.OptionLegResponse, error) {
	// Get existing option leg
	optionLeg, err := s.optionLegRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get option leg: %w", err)
	}

	if optionLeg == nil {
		return nil, ErrOptionLegNotFound
	}

	// Update fields if provided
	if req.LegSequence != nil {
		// Check if new leg sequence conflicts with existing ones
		if *req.LegSequence != optionLeg.LegSequence {
			existingLeg, err := s.optionLegRepo.GetByStrategyIDAndSequence(ctx, optionLeg.StrategyID, *req.LegSequence)
			if err == nil && existingLeg != nil {
				return nil, ErrDuplicateOrder
			}
		}
		optionLeg.LegSequence = *req.LegSequence
	}

	if req.SelectIndex != nil {
		optionLeg.SelectIndex = *req.SelectIndex
	}
	if req.Segment != nil {
		optionLeg.Segment = *req.Segment
	}
	if req.OptionType != nil {
		optionLeg.OptionType = req.OptionType
	}
	if req.ActionType != nil {
		optionLeg.ActionType = *req.ActionType
	}
	if req.StrikeType != nil {
		optionLeg.StrikeType = *req.StrikeType
	}
	if req.StrikePrice != nil {
		optionLeg.StrikePrice = req.StrikePrice
	}
	if req.StrikeOffset != nil {
		optionLeg.StrikeOffset = req.StrikeOffset
	}
	if req.TotalLot != nil {
		optionLeg.TotalLot = *req.TotalLot
	}
	if req.LotSize != nil {
		optionLeg.LotSize = req.LotSize
	}
	if req.OrderType != nil {
		optionLeg.OrderType = *req.OrderType
	}
	if req.TickOffset != nil {
		optionLeg.TickOffset = req.TickOffset
	}
	if req.ExpiryType != nil {
		optionLeg.ExpiryType = *req.ExpiryType
	}
	if req.ExpiryCycle != nil {
		optionLeg.ExpiryCycle = *req.ExpiryCycle
	}
	if req.ExpiryDate != nil {
		optionLeg.ExpiryDate = req.ExpiryDate
	}
	if req.DaysToExpiry != nil {
		optionLeg.DaysToExpiry = req.DaysToExpiry
	}

	// Recalculate total quantity if lot size or total lot changed
	if req.LotSize != nil || req.TotalLot != nil {
		if optionLeg.LotSize != nil {
			totalQuantity := optionLeg.TotalLot * *optionLeg.LotSize
			optionLeg.TotalQuantity = &totalQuantity
		}
	}

	// Validate the updated option leg
	if err := optionLeg.IsValid(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update timestamp
	optionLeg.UpdatedAt = time.Now()

	// Save to database
	if err := s.optionLegRepo.Update(ctx, optionLeg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOptionLegUpdateFailed, err)
	}

	return dto.ToResponse(optionLeg), nil
}

// DeleteOptionLeg deletes an option leg
func (s *optionLegService) DeleteOptionLeg(ctx context.Context, id uuid.UUID) error {
	// Check if option leg exists
	optionLeg, err := s.optionLegRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get option leg: %w", err)
	}

	if optionLeg == nil {
		return ErrOptionLegNotFound
	}

	// Delete from database
	if err := s.optionLegRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrOptionLegDeleteFailed, err)
	}

	return nil
}

// DeleteOptionLegsByStrategyID deletes all option legs for a strategy
func (s *optionLegService) DeleteOptionLegsByStrategyID(ctx context.Context, strategyID uuid.UUID) error {
	if err := s.optionLegRepo.DeleteByStrategyID(ctx, strategyID); err != nil {
		return fmt.Errorf("%w: %v", ErrOptionLegDeleteFailed, err)
	}

	return nil
}

// GetOptionLegsByIndex retrieves option legs by index/underlying
func (s *optionLegService) GetOptionLegsByIndex(ctx context.Context, selectIndex string) ([]*dto.OptionLegResponse, error) {
	optionLegs, err := s.optionLegRepo.GetByIndex(ctx, selectIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get option legs by index: %w", err)
	}

	responses := make([]*dto.OptionLegResponse, len(optionLegs))
	for i, optionLeg := range optionLegs {
		responses[i] = dto.ToResponse(optionLeg)
	}

	return responses, nil
}

// SearchOptionLegs searches option legs with filter criteria
func (s *optionLegService) SearchOptionLegs(ctx context.Context, filter *dto.OptionLegFilter) ([]*dto.OptionLegResponse, int64, error) {
	// Set default pagination if not provided
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}

	offset := (filter.Page - 1) * filter.PerPage

	optionLegs, total, err := s.optionLegRepo.List(ctx, offset, filter.PerPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search option legs: %w", err)
	}

	responses := make([]*dto.OptionLegResponse, len(optionLegs))
	for i, optionLeg := range optionLegs {
		responses[i] = dto.ToResponse(optionLeg)
	}

	return responses, total, nil
}

// ValidateOptionLeg validates option leg data
func (s *optionLegService) ValidateOptionLeg(ctx context.Context, req *dto.CreateOptionLegRequest) error {
	// Validate F&O segment requirements
	if req.Segment == model.SegmentFO {
		if req.OptionType == nil {
			return model.ErrInvalidOptionDetails
		}
	}

	// Validate EQ segment requirements
	if req.Segment == model.SegmentEQ {
		if req.OptionType != nil || req.StrikePrice != nil {
			return model.ErrInvalidEquityDetails
		}
	}

	// Validate leg sequence
	if req.LegSequence <= 0 {
		return model.ErrInvalidOrder
	}

	// Validate total lot
	if req.TotalLot <= 0 {
		return model.ErrInvalidTotalLot
	}

	return nil
}

// ResolveExchangeSecurityID resolves the exchange_security_id for an option leg based on spot price
func (s *optionLegService) ResolveExchangeSecurityID(ctx context.Context, legID uuid.UUID, spotPrice float64) (int32, error) {
	// Get the option leg details
	optionLeg, err := s.optionLegRepo.GetByID(ctx, legID)
	if err != nil {
		return 0, fmt.Errorf("failed to get option leg: %w", err)
	}

	if optionLeg == nil {
		return 0, ErrOptionLegNotFound
	}

	// Only process F&O segment options
	if optionLeg.Segment != model.SegmentFO {
		return 0, fmt.Errorf("exchange security ID resolution is only supported for F&O segment")
	}

	if optionLeg.OptionType == nil {
		return 0, fmt.Errorf("option type is required for F&O segment")
	}

	// Step 1: Resolve the nearest strike price
	nearestStrike := s.calculateNearestStrike(spotPrice)

	// Step 2: Format expiry month from expiry date or expiry type
	expiryMonth, err := s.formatExpiryMonth(optionLeg)
	if err != nil {
		return 0, fmt.Errorf("failed to format expiry month: %w", err)
	}

	// Step 3: Prepare trading symbol in format "[index]-[expirymonth]-[strikePrice]-[CE/PE]"
	optionTypeStr := "CE"
	if *optionLeg.OptionType == model.OptionTypePut {
		optionTypeStr = "PE"
	}

	tradingSymbol := fmt.Sprintf("%s-%s-%.0f-%s",
		strings.ToUpper(optionLeg.SelectIndex),
		expiryMonth,
		nearestStrike,
		optionTypeStr)

	// Step 4: Query instrument master table using the trading symbol
	instrumentMasters, _, err := s.instrumentService.GetInstrumentMastersWithFilters(ctx, "", instrumentModel.SegmentDerivatives, tradingSymbol, 1, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to query instrument master: %w", err)
	}

	if len(instrumentMasters) == 0 {
		return 0, fmt.Errorf("no instrument found for trading symbol: %s", tradingSymbol)
	}

	return instrumentMasters[0].ExchangeSecurityID, nil
}

// calculateNearestStrike calculates the nearest strike price based on spot price
// For Indian markets, strikes are typically rounded to nearest 50 or 100
func (s *optionLegService) calculateNearestStrike(spotPrice float64) float64 {
	// Determine the strike interval based on spot price range
	var interval float64

	if spotPrice < 500 {
		interval = 10 // For lower priced stocks/indices
	} else if spotPrice < 1000 {
		interval = 25
	} else if spotPrice < 3000 {
		interval = 50
	} else if spotPrice < 10000 {
		interval = 100
	} else {
		interval = 200 // For higher priced indices like Nifty 50
	}

	// Round to nearest interval
	return math.Round(spotPrice/interval) * interval
}

// formatExpiryMonth formats the expiry month based on expiry date or expiry type
func (s *optionLegService) formatExpiryMonth(optionLeg *model.OptionLeg) (string, error) {
	now := time.Now()

	if optionLeg.ExpiryDate != nil {
		// Use the provided expiry date
		return s.getExpiryString(*optionLeg.ExpiryDate), nil
	}

	// Calculate expiry based on expiry type
	var expiryDate time.Time

	switch optionLeg.ExpiryType {
	case model.ExpiryTypeWeekly:
		// Find next Thursday
		daysUntilThursday := (4 - int(now.Weekday()) + 7) % 7
		if daysUntilThursday == 0 && now.Weekday() == time.Thursday {
			daysUntilThursday = 7 // Next week if today is Thursday
		}
		expiryDate = now.AddDate(0, 0, daysUntilThursday)

	case model.ExpiryTypeMonthly:
		// Find last Thursday of current month
		expiryDate = s.getLastThursdayOfMonth(now.Year(), now.Month())
		if expiryDate.Before(now) {
			// Move to next month if current month's expiry has passed
			nextMonth := now.AddDate(0, 1, 0)
			expiryDate = s.getLastThursdayOfMonth(nextMonth.Year(), nextMonth.Month())
		}

	case model.ExpiryTypeQuarterly:
		// Find last Thursday of current quarter
		expiryDate = s.getLastThursdayOfQuarter(now)
		if expiryDate.Before(now) {
			// Move to next quarter if current quarter's expiry has passed
			nextQuarter := now.AddDate(0, 3, 0)
			expiryDate = s.getLastThursdayOfQuarter(nextQuarter)
		}

	default:
		return "", fmt.Errorf("unsupported expiry type: %s", optionLeg.ExpiryType)
	}

	return s.getExpiryString(expiryDate), nil
}

// getExpiryString formats date to expiry string format (e.g., "25JAN", "25FEB")
func (s *optionLegService) getExpiryString(date time.Time) string {
	monthNames := []string{
		"JAN", "FEB", "MAR", "APR", "MAY", "JUN",
		"JUL", "AUG", "SEP", "OCT", "NOV", "DEC",
	}

	year := fmt.Sprintf("%02d", date.Year()%100) // Last 2 digits of year
	month := monthNames[date.Month()-1]

	return fmt.Sprintf("%s%s", year, month)
}

// getLastThursdayOfMonth finds the last Thursday of a given month
func (s *optionLegService) getLastThursdayOfMonth(year int, month time.Month) time.Time {
	// Start from the last day of the month
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	// Work backwards to find Thursday
	for lastDay.Weekday() != time.Thursday {
		lastDay = lastDay.AddDate(0, 0, -1)
	}

	return lastDay
}

// getLastThursdayOfQuarter finds the last Thursday of the current quarter
func (s *optionLegService) getLastThursdayOfQuarter(date time.Time) time.Time {
	year := date.Year()
	month := date.Month()

	// Determine the last month of the current quarter
	var quarterEndMonth time.Month
	if month <= 3 {
		quarterEndMonth = 3 // March
	} else if month <= 6 {
		quarterEndMonth = 6 // June
	} else if month <= 9 {
		quarterEndMonth = 9 // September
	} else {
		quarterEndMonth = 12 // December
	}

	return s.getLastThursdayOfMonth(year, quarterEndMonth)
}
