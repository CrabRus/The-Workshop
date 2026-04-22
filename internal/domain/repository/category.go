package repository

import (
	"context"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(ctx context.Context, product *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter CategoryFilter) ([]*entity.Category, int, error)
}

type CategoryFilter struct {
	Search  string
	Limit   int
	Offset  int
	OrderBy string
}
