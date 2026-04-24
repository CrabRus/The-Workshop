package postgres

import (
	"context"
	"fmt"

	database "github.com/crabrus/the-workshop/internal/db"
	"github.com/crabrus/the-workshop/internal/domain/entity"
	repository "github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/google/uuid"
)

type orderItemRepo struct {
	db *database.DB
}

func NewOrderItemRepository(db *database.DB) repository.OrderItemRepository {
	return &orderItemRepo{db: db}
}

// Create implements repository.OrderItemRepository.
func (r *orderItemRepo) Create(ctx context.Context, item *entity.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	return r.db.QueryRowContext(
		ctx, query,
		item.ID, item.OrderID, item.ProductID, item.Quantity, item.Price,
	).Scan(&item.CreatedAt)
}

// GetByOrderID implements repository.OrderItemRepository.
func (r *orderItemRepo) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entity.OrderItem, error) {
	var items []*entity.OrderItem

	query := `
		SELECT id, order_id, product_id, quantity, price, created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	err := r.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	return items, nil
}

// Delete implements repository.OrderItemRepository.
func (r *orderItemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM order_items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order item not found")
	}

	return nil
}

// DeleteByOrderID implements repository.OrderItemRepository.
func (r *orderItemRepo) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	query := `DELETE FROM order_items WHERE order_id = $1`

	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order items by order_id: %w", err)
	}

	return nil
}
