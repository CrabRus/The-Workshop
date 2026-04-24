package repository

import (
	"context"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/google/uuid"
)

type OrderItemRepository interface {
	Create(ctx context.Context, item *entity.OrderItem) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entity.OrderItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error
}

type OrderItemFilter struct {
	OrderID uuid.UUID
	Limit   int
	Offset  int
}
