package product

import (
	"context"
	"fmt"
	"time"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/crabrus/the-workshop/pkg/utils"
	"github.com/google/uuid"
)

type ProductService interface {
	List(ctx context.Context, filter repository.ProductFilter) (*ProductListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	Create(ctx context.Context, req ProductRequest) (*entity.Product, error)
	Update(ctx context.Context, id uuid.UUID, req ProductRequest) (*entity.Product, error)
	Delete(ctx context.Context, id uuid.UUID, isAdmin bool) error
}

type service struct {
	productRepo repository.ProductRepository
}

func NewService(productRepo repository.ProductRepository) ProductService {
	return &service{productRepo: productRepo}
}

// List implements ProductService.
func (s *service) List(ctx context.Context, filter repository.ProductFilter) (*ProductListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	products, total, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return &ProductListResponse{
		Products: products,
		Total:    total,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil
}

// Create implements ProductService.
func (s *service) Create(ctx context.Context, req ProductRequest) (*entity.Product, error) {
	role, ok := utils.GetRoleFromContext(ctx)
	if !ok || role != "admin" {
		return nil, ErrAdminOnly
	}

	if err := ValidateProductCreate(req); err != nil {
		return nil, err
	}

	existingProduct, err := s.productRepo.GetByName(ctx, *req.Name)
	if err == nil && existingProduct != nil {
		return nil, ErrProductAlreadyExists
	}

	product := &entity.Product{
		Name:        *req.Name,
		Description: derefString(req.Description),
		Price:       *req.Price,
		Stock:       *req.Stock,
		CategoryID:  *req.CategoryID,
		Image_url:   derefString(req.ImageURL),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// Delete implements ProductService.
func (s *service) Delete(ctx context.Context, id uuid.UUID, isAdmin bool) error {
	role, ok := utils.GetRoleFromContext(ctx)
	if !ok || role != "admin" {
		return ErrAdminOnly
	}
	return s.productRepo.Delete(ctx, id)
}

// GetByID implements ProductService.
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

// Update implements ProductService.
func (s *service) Update(ctx context.Context, id uuid.UUID, req ProductRequest) (*entity.Product, error) {
	role, ok := utils.GetRoleFromContext(ctx)
	if !ok || role != "admin" {
		return nil, ErrAdminOnly
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if err := ValidateProductUpdate(req); err != nil {
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}

	if req.Description != nil {
		product.Description = *req.Description
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}

	if req.ImageURL != nil {
		product.Image_url = *req.ImageURL
	}

	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}
