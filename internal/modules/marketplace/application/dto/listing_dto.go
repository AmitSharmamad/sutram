// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/dto/listing_dto.go
package dto

import (
	"time"

	"algo/internal/modules/marketplace/domain/model"

	"github.com/google/uuid"
)

// CreateListingRequest represents a request to create a new marketplace listing
type CreateListingRequest struct {
	StrategyID       uuid.UUID `json:"strategy_id" binding:"required"`
	Title            string    `json:"title" binding:"required,min=5,max=255"`
	Description      string    `json:"description" binding:"required"`
	Price            float64   `json:"price" binding:"required,min=0"`
	SubscriptionType string    `json:"subscription_type" binding:"required,oneof=one-time monthly yearly"`
	CategoryIDs      []string  `json:"category_ids" binding:"required,dive,uuid4"`
	Features         []string  `json:"features" binding:"omitempty,dive,max=100"`
}

// UpdateListingRequest represents a request to update a marketplace listing
type UpdateListingRequest struct {
	Title            string   `json:"title" binding:"omitempty,min=5,max=255"`
	Description      string   `json:"description"`
	Price            float64  `json:"price" binding:"omitempty,min=0"`
	SubscriptionType string   `json:"subscription_type" binding:"omitempty,oneof=one-time monthly yearly"`
	CategoryIDs      []string `json:"category_ids" binding:"omitempty,dive,uuid4"`
	Features         []string `json:"features" binding:"omitempty,dive,max=100"`
	IsActive         *bool    `json:"is_active"`
}

// ListingResponse represents a marketplace listing response
type ListingResponse struct {
	ID               uuid.UUID            `json:"id"`
	StrategyID       uuid.UUID            `json:"strategy_id"`
	UserID           uuid.UUID            `json:"user_id"`
	SellerName       string               `json:"seller_name,omitempty"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	Price            float64              `json:"price"`
	SubscriptionType string               `json:"subscription_type"`
	IsFeatured       bool                 `json:"is_featured"`
	IsActive         bool                 `json:"is_active"`
	Views            int                  `json:"views"`
	Ratings          float32              `json:"ratings"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	Categories       []CategoryResponse   `json:"categories,omitempty"`
	Performance      *PerformanceResponse `json:"performance,omitempty"`
	Features         []string             `json:"features,omitempty"`
}

// ListingFilterRequest represents filters for listing API
type ListingFilterRequest struct {
	Category string `form:"category"`
	Sort     string `form:"sort" binding:"omitempty,oneof=newest price rating popularity"`
	Search   string `form:"search"`
	PriceMin string `form:"price_min"`
	PriceMax string `form:"price_max"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PerPage  int    `form:"per_page,default=20" binding:"min=1,max=100"`
}

// ToListingResponse converts a model.Listing to a ListingResponse
func ToListingResponse(listing *model.Listing) ListingResponse {
	var categoryResponses []CategoryResponse
	if len(listing.Categories) > 0 {
		categoryResponses = make([]CategoryResponse, len(listing.Categories))
		for i, category := range listing.Categories {
			categoryResponses[i] = CategoryResponse{
				ID:   category.ID,
				Name: category.Name,
			}
		}
	}

	var performanceResponse *PerformanceResponse
	if listing.Performance != nil {
		performanceResponse = &PerformanceResponse{
			SharpeRatio:        listing.Performance.SharpeRatio,
			SortinoRatio:       listing.Performance.SortinoRatio,
			TotalReturnPct:     listing.Performance.TotalReturnPct,
			AnnualizedReturn:   listing.Performance.AnnualizedReturn,
			MaxDrawdown:        listing.Performance.MaxDrawdown,
			WinRate:            listing.Performance.WinRate,
			ProfitFactor:       listing.Performance.ProfitFactor,
			BacktestPeriodDays: listing.Performance.BacktestPeriodDays,
		}
	}

	return ListingResponse{
		ID:               listing.ID,
		StrategyID:       listing.StrategyID,
		UserID:           listing.UserID,
		Title:            listing.Title,
		Description:      listing.Description,
		Price:            listing.Price,
		SubscriptionType: string(listing.SubscriptionType),
		IsFeatured:       listing.IsFeatured,
		IsActive:         listing.IsActive,
		Views:            listing.Views,
		Ratings:          listing.Ratings,
		CreatedAt:        listing.CreatedAt,
		UpdatedAt:        listing.UpdatedAt,
		Categories:       categoryResponses,
		Performance:      performanceResponse,
		Features:         listing.Features,
	}
}
