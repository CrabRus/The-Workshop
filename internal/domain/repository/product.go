package repository

import (
	"context"

	entity "github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetByCategoryID(ctx context.Context, category_id uuid.UUID) (*entity.Product, error)
	GetByName(ctx context.Context, name string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter ProductFilter) ([]*entity.Product, int, error)
}

type ProductFilter struct {
	Search     string
	CategoryID *string
	Limit      int
	Offset     int
	OrderBy    string
}
