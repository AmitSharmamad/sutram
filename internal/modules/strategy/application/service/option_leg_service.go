package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/domain/model"
	"algo/internal/modules/strategy/domain/repository"

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
	CreateOptionLeg(ctx context.Context, req *dto.CreateOptionLegRequest) (*dto.OptionLegResponse, error)

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
}

// optionLegService implements the OptionLegService interface
type optionLegService struct {
	optionLegRepo repository.OptionLegRepository
}

// NewOptionLegService creates a new option leg service
func NewOptionLegService(optionLegRepo repository.OptionLegRepository) OptionLegService {
	return &optionLegService{
		optionLegRepo: optionLegRepo,
	}
}

// CreateOptionLeg creates a new option leg
func (s *optionLegService) CreateOptionLeg(ctx context.Context, req *dto.CreateOptionLegRequest) (*dto.OptionLegResponse, error) {
	// Validate the request
	if err := s.ValidateOptionLeg(ctx, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if leg sequence already exists for the strategy
	existingLeg, err := s.optionLegRepo.GetByStrategyIDAndSequence(ctx, req.StrategyID, req.LegSequence)
	if err == nil && existingLeg != nil {
		return nil, ErrDuplicateOrder
	}

	// Convert DTO to domain model
	optionLeg := req.ToModel()

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
	if req.Segment == model.FO {
		if req.OptionType == nil {
			return model.ErrInvalidOptionDetails
		}
		if req.StrikePrice == nil && req.StrikeType == model.StrikePriceTypeCustom {
			return model.ErrInvalidStrikePrice
		}
	}

	// Validate EQ segment requirements
	if req.Segment == model.EQ {
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
