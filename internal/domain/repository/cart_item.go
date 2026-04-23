package repository

import (
	"context"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/google/uuid"
)

type CartItemRepository interface {
	Create(ctx context.Context, cartItem *entity.CartItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CartItem, error)
	Update(ctx context.Context, cartItem *entity.CartItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter CartItemFilter) ([]*entity.CartItem, int, error)
}

type CartItemFilter struct {
	Search          string
	UserID          *string
	OrderByPrice    string
	OrderByquantity string
	OrderBy         string
}
