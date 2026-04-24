package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	database "github.com/crabrus/the-workshop/internal/db"
	"github.com/crabrus/the-workshop/internal/domain/entity"
	repository "github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type orderRepo struct {
	db *database.DB
}

func NewOrderRepository(db *database.DB) repository.OrderRepository {
	return &orderRepo{db: db}
}

// Create implements repository.OrderRepository.
func (r *orderRepo) Create(ctx context.Context, order *entity.Order) error {
	query := `
		INSERT INTO orders (id, user_id, status, total_amount, shipping_address, payment_method)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}

	addressJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal shipping address: %w", err)
	}

	return r.db.QueryRowContext(
		ctx, query,
		order.ID, order.UserID, entity.OrderStatusPending, order.TotalAmount, addressJSON, order.PaymentMethod,
	).Scan(&order.CreatedAt, &order.UpdatedAt)
}

// GetByID implements repository.OrderRepository.
func (r *orderRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order

	query := `
		SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	return &order, nil
}

// GetByUserID implements repository.OrderRepository.
func (r *orderRepo) GetByUserID(ctx context.Context, userID uuid.UUID, filter repository.OrderFilter) ([]*entity.Order, int, error) {
	var orders []*entity.Order

	query := `
		SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
		FROM orders
		WHERE user_id = $1
	`

	if filter.Status != "" {
		query += ` AND status = $2`
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", 3) + ` OFFSET $` + fmt.Sprintf("%d", 4)

	args := []interface{}{userID, filter.Limit, filter.Offset}
	if filter.Status != "" {
		args = []interface{}{userID, filter.Status, filter.Limit, filter.Offset}
		query = `
			SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
			FROM orders
			WHERE user_id = $1 AND status = $2
			ORDER BY created_at DESC LIMIT $3 OFFSET $4
		`
	} else {
		query = `
			SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
			FROM orders
			WHERE user_id = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3
		`
	}

	err := r.db.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders by user: %w", err)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	countArgs := []interface{}{userID}
	if filter.Status != "" {
		countQuery += ` AND status = $2`
		countArgs = append(countArgs, filter.Status)
	}

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	return orders, total, nil
}

// Update implements repository.OrderRepository.
func (r *orderRepo) Update(ctx context.Context, order *entity.Order) error {
	addressJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal shipping address: %w", err)
	}

	query := `
		UPDATE orders
		SET status = $1, total_amount = $2, shipping_address = $3, payment_method = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		order.Status, order.TotalAmount, addressJSON, order.PaymentMethod, order.ID,
	).Scan(&order.UpdatedAt)
}

// UpdateStatus implements repository.OrderRepository.
func (r *orderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrOrderNotFound
	}

	return nil
}

// List implements repository.OrderRepository.
func (r *orderRepo) List(ctx context.Context, filter repository.OrderFilter) ([]*entity.Order, int, error) {
	var orders []*entity.Order

	query := `
		SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
		FROM orders
	`

	if filter.Status != "" {
		query += ` WHERE status = $1`
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", 2) + ` OFFSET $` + fmt.Sprintf("%d", 3)

	args := []interface{}{filter.Limit, filter.Offset}
	if filter.Status != "" {
		args = []interface{}{filter.Status, filter.Limit, filter.Offset}
		query = `
			SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
			FROM orders
			WHERE status = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3
		`
	} else {
		query = `
			SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
			FROM orders
			ORDER BY created_at DESC LIMIT $1 OFFSET $2
		`
	}

	err := r.db.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM orders`
	countArgs := []interface{}{}
	if filter.Status != "" {
		countQuery += ` WHERE status = $1`
		countArgs = append(countArgs, filter.Status)
	}

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	return orders, total, nil
}

// Delete implements repository.OrderRepository.
func (r *orderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		return ErrOrderNotFound
	}

	return nil
}

// GetStatistics повертає агреговані дані про продажі
func (r *orderRepo) GetStatistics(ctx context.Context) (int, float64, []repository.TopProduct, error) {
	var stats struct {
		Count   int     `db:"count"`
		Revenue float64 `db:"revenue"`
	}

	// Загальна кількість та дохід (крім скасованих)
	queryStats := `
		SELECT 
			COUNT(*) as count, 
			COALESCE(SUM(total_amount), 0) as revenue 
		FROM orders 
		WHERE status != 'cancelled'
	`
	if err := r.db.GetContext(ctx, &stats, queryStats); err != nil {
		return 0, 0, nil, fmt.Errorf("failed to get general stats: %w", err)
	}

	// Топ 5 товарів за кількістю продажів
	queryTop := `
		SELECT p.name, SUM(oi.quantity) as total_sold
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		JOIN orders o ON oi.order_id = o.id
		WHERE o.status != 'cancelled'
		GROUP BY p.name
		ORDER BY total_sold DESC
		LIMIT 5
	`
	var topProducts []repository.TopProduct
	rows, err := r.db.QueryContext(ctx, queryTop)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to get top products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tp repository.TopProduct
		if err := rows.Scan(&tp.Name, &tp.TotalSold); err != nil {
			return 0, 0, nil, err
		}
		topProducts = append(topProducts, tp)
	}

	return stats.Count, stats.Revenue, topProducts, nil
}
