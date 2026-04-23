package postgres

import (
	"context"
	"errors"
	"fmt"

	database "github.com/crabrus/the-workshop/internal/db"
	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrCartItemNotFound = errors.New("cart item not found")
)

type cartItemRepo struct {
	db *database.DB
}

func NewCartItemRepository(db *database.DB) repository.CartItemRepository {
	return &cartItemRepo{db: db}
}

// Create implements repository.CartItemRepository.
func (c *cartItemRepo) Create(ctx context.Context, cartItem *entity.CartItem) error {
	query := `
		INSERT INTO cart_items (id, user_id, product_id, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	if cartItem.ID == uuid.Nil {
		cartItem.ID = uuid.New()
	}

	return c.db.QueryRowContext(
		ctx, query,
		cartItem.ID, cartItem.UserID, cartItem.ProductID, cartItem.Quantity,
	).Scan(&cartItem.CreatedAt)
}

// Delete implements repository.CartItemRepository.
func (c *cartItemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id=$1`

	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		return ErrCartItemNotFound
	}
	return nil
}

// GetByID implements repository.CartItemRepository.
func (c *cartItemRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.CartItem, error) {
	var cartItem entity.CartItem

	query := `
		SELECT id, user_id, product_id, quantity, created_at
		FROM cart_items
		WHERE id = $1
	`

	err := c.db.GetContext(ctx, &cartItem, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart item by id: %w", err)
	}

	return &cartItem, nil
}

// List implements repository.CartItemRepository.
func (c *cartItemRepo) List(ctx context.Context, filter repository.CartItemFilter) ([]*entity.CartItem, int, error) {
	// 1. Формируем WHERE условие
	where := "WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if filter.UserID != nil && *filter.UserID != "" {
		userID, err := uuid.Parse(*filter.UserID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid user_id format: %w", err)
		}
		where += fmt.Sprintf(" AND user_id = $%d", argID)
		args = append(args, userID)
		argID++
	}

	// 2. Формируем ORDER BY
	orderBy := "created_at DESC"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	} else if filter.OrderByPrice != "" {
		orderBy = filter.OrderByPrice
	} else if filter.OrderByquantity != "" {
		orderBy = filter.OrderByquantity
	}

	// 3. Считаем Total (без LIMIT/OFFSET)
	var total int
	countQuery := "SELECT count(*) FROM cart_items " + where
	if err := c.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count cart items: %w", err)
	}

	// 4. Получаем данные
	selectQuery := fmt.Sprintf(
		"SELECT id, user_id, product_id, quantity, created_at FROM cart_items %s ORDER BY %s",
		where, orderBy,
	)

	var cartItems []*entity.CartItem
	if err := c.db.SelectContext(ctx, &cartItems, selectQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to list cart items: %w", err)
	}

	return cartItems, total, nil
}

// Update implements repository.CartItemRepository.
func (c *cartItemRepo) Update(ctx context.Context, cartItem *entity.CartItem) error {
	query := `
		UPDATE cart_items
		SET user_id = $1,
			product_id = $2,
			quantity = $3
		WHERE id = $4
	`

	result, err := c.db.ExecContext(
		ctx, query,
		cartItem.UserID, cartItem.ProductID, cartItem.Quantity, cartItem.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		return ErrCartItemNotFound
	}

	return nil
}
