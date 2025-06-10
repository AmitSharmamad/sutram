package service

import (
	"context"
	"fmt"
	"time"

	"algo/internal/modules/broker/application/dto"
	"algo/internal/modules/broker/domain/model"
	"algo/internal/modules/broker/domain/repository"

	"github.com/google/uuid"
)

// BrokerService handles broker-related business logic
type BrokerService struct {
	brokerRepo     repository.BrokerRepository
	brokerUserRepo repository.BrokerUserRepository
	orderRepo      repository.OrderRepository
	dhanService    *DhanOrderService
}

// NewBrokerService creates a new broker service
func NewBrokerService(
	brokerRepo repository.BrokerRepository,
	brokerUserRepo repository.BrokerUserRepository,
	orderRepo repository.OrderRepository,
	dhanService *DhanOrderService,
) *BrokerService {
	return &BrokerService{
		brokerRepo:     brokerRepo,
		brokerUserRepo: brokerUserRepo,
		orderRepo:      orderRepo,
		dhanService:    dhanService,
	}
}

// PlaceOrderWithDhan places an order through Dhan and saves it locally
func (s *BrokerService) PlaceOrderWithDhan(ctx context.Context, orderDTO *dto.PlaceOrderDTO) (*dto.OrderResponseDTO, error) {
	// Convert internal DTO to Dhan request
	dhanRequest := &dto.DhanPlaceOrderRequest{
		DhanClientID:    "default", // This should come from user configuration
		TransactionType: orderDTO.TransactionType,
		ExchangeSegment: fmt.Sprintf("%s_%s", orderDTO.Exchange, orderDTO.Segment),
		ProductType:     orderDTO.ProductType,
		OrderType:       orderDTO.OrderType,
		Validity:        "DAY", // Default validity
		SecurityID:      orderDTO.Instrument,
		Quantity:        orderDTO.Quantity,
		Price:           orderDTO.Price,
	}

	// Place order with Dhan
	dhanResponse, err := s.dhanService.PlaceDhanOrder(ctx, dhanRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to place order with Dhan: %w", err)
	}

	// Create local order record
	order := &model.Order{
		ID:                  orderDTO.ID,
		UserID:              orderDTO.UserID,
		PositionID:          orderDTO.PositionID,
		Exchange:            orderDTO.Exchange,
		Segment:             orderDTO.Segment,
		Instrument:          orderDTO.Instrument,
		TransactionType:     orderDTO.TransactionType,
		Quantity:            orderDTO.Quantity,
		Price:               orderDTO.Price,
		State:               "SENT",
		ProductType:         orderDTO.ProductType,
		OrderType:           orderDTO.OrderType,
		BrokerOrderID:       &dhanResponse.OrderID,
		BrokerOrderResponse: dhanResponse,
		AlgoID:              &orderDTO.AlgoID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Save order to database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Convert to response DTO
	return &dto.OrderResponseDTO{
		ID:                  order.ID,
		UserID:              order.UserID,
		PositionID:          order.PositionID,
		Exchange:            order.Exchange,
		Segment:             order.Segment,
		Instrument:          order.Instrument,
		TransactionType:     order.TransactionType,
		Quantity:            order.Quantity,
		Price:               order.Price,
		State:               order.State,
		ProductType:         order.ProductType,
		OrderType:           order.OrderType,
		BrokerOrderID:       *order.BrokerOrderID,
		BrokerOrderResponse: order.BrokerOrderResponse,
		AlgoID:              *order.AlgoID,
		CreatedAt:           order.CreatedAt,
		UpdatedAt:           order.UpdatedAt,
	}, nil
}

// GetOrderStatus retrieves order status from Dhan and updates local record
func (s *BrokerService) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (*dto.OrderResponseDTO, error) {
	// Get local order record
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	if order.BrokerOrderID == nil {
		return nil, fmt.Errorf("order has no broker order ID")
	}

	// Get status from Dhan
	dhanResponse, err := s.dhanService.GetOrderByID(ctx, *order.BrokerOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status from Dhan: %w", err)
	}

	// Update local record with latest status
	order.State = dhanResponse.Order.OrderStatus
	order.BrokerOrderResponse = dhanResponse.Order
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Convert to response DTO
	return &dto.OrderResponseDTO{
		ID:                  order.ID,
		UserID:              order.UserID,
		PositionID:          order.PositionID,
		Exchange:            order.Exchange,
		Segment:             order.Segment,
		Instrument:          order.Instrument,
		TransactionType:     order.TransactionType,
		Quantity:            order.Quantity,
		Price:               order.Price,
		State:               order.State,
		ProductType:         order.ProductType,
		OrderType:           order.OrderType,
		BrokerOrderID:       *order.BrokerOrderID,
		BrokerOrderResponse: order.BrokerOrderResponse,
		AlgoID:              *order.AlgoID,
		CreatedAt:           order.CreatedAt,
		UpdatedAt:           order.UpdatedAt,
	}, nil
}

// CancelOrderWithDhan cancels an order through Dhan and updates local record
func (s *BrokerService) CancelOrderWithDhan(ctx context.Context, orderID uuid.UUID) (*dto.OrderResponseDTO, error) {
	// Get local order record
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	if order.BrokerOrderID == nil {
		return nil, fmt.Errorf("order has no broker order ID")
	}

	// Cancel order with Dhan
	cancelRequest := &dto.DhanCancelOrderRequest{
		DhanClientID: "default", // This should come from user configuration
		OrderID:      *order.BrokerOrderID,
	}

	dhanResponse, err := s.dhanService.CancelDhanOrder(ctx, cancelRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order with Dhan: %w", err)
	}

	// Update local record
	order.State = dhanResponse.OrderStatus
	order.BrokerOrderResponse = dhanResponse
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Convert to response DTO
	return &dto.OrderResponseDTO{
		ID:                  order.ID,
		UserID:              order.UserID,
		PositionID:          order.PositionID,
		Exchange:            order.Exchange,
		Segment:             order.Segment,
		Instrument:          order.Instrument,
		TransactionType:     order.TransactionType,
		Quantity:            order.Quantity,
		Price:               order.Price,
		State:               order.State,
		ProductType:         order.ProductType,
		OrderType:           order.OrderType,
		BrokerOrderID:       *order.BrokerOrderID,
		BrokerOrderResponse: order.BrokerOrderResponse,
		AlgoID:              *order.AlgoID,
		CreatedAt:           order.CreatedAt,
		UpdatedAt:           order.UpdatedAt,
	}, nil
}