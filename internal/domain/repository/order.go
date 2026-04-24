package repository

import (
	"context"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, filter OrderFilter) ([]*entity.Order, int, error)
	Update(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	List(ctx context.Context, filter OrderFilter) ([]*entity.Order, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type OrderFilter struct {
	Status string
	Limit  int
	Offset int
}
