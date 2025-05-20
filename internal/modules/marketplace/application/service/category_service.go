// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/service/category_service.go
package service

import (
	"context"
	"fmt"

	"algo/internal/modules/marketplace/domain/model"
	"algo/internal/modules/marketplace/domain/repository"

	"github.com/google/uuid"
)

// CategoryService handles marketplace category business logic
type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// GetAllCategories retrieves all marketplace categories
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if category == nil {
		return nil, ErrCategoryNotFound
	}

	return category, nil
}

// GetCategoriesByIDs retrieves categories by their IDs
func (s *CategoryService) GetCategoriesByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Category, error) {
	if len(ids) == 0 {
		return []model.Category{}, nil
	}

	categories, err := s.categoryRepo.FindByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by IDs: %w", err)
	}

	return categories, nil
}
