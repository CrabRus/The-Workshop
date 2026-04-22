package category

import (
	"context"
	"fmt"
	"time"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/crabrus/the-workshop/pkg/utils"
	"github.com/google/uuid"
)

type CategoryService interface {
	List(ctx context.Context, filter repository.CategoryFilter) (*CategoryListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	Create(ctx context.Context, req CategoryRequest) (*entity.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	categoryRepo repository.CategoryRepository
}

func NewService(categoryRepo repository.CategoryRepository) CategoryService {
	return &service{categoryRepo: categoryRepo}
}

// Create implements CategoryService.
func (s *service) Create(ctx context.Context, req CategoryRequest) (*entity.Category, error) {
	role, ok := utils.GetRoleFromContext(ctx)
	if !ok || role != "admin" {
		return nil, ErrAdminOnly
	}

	if err := ValidateCategoryCreate(req); err != nil {
		return nil, err
	}

	existingCategory, err := s.categoryRepo.GetByName(ctx, *req.Name)
	if err == nil && existingCategory != nil {
		return nil, ErrCategoryAlreadyExists
	}

	category := &entity.Category{
		Name:        *req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// Delete implements CategoryService.
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	role, ok := utils.GetRoleFromContext(ctx)
	if !ok || role != "admin" {
		return ErrAdminOnly
	}
	return s.categoryRepo.Delete(ctx, id)
}

// GetByID implements CategoryService.
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return category, nil
}

// List implements CategoryService.
func (s *service) List(ctx context.Context, filter repository.CategoryFilter) (*CategoryListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	categories, total, err := s.categoryRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return &CategoryListResponse{
		Categories: categories,
		Total:      total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	}, nil
}
