package postgres

import (
	"context"
	"errors"
	"fmt"

	database "github.com/crabrus/the-workshop/internal/db"
	"github.com/crabrus/the-workshop/internal/domain/entity"
	repository "github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrProductNotFount = errors.New("product not found")
)

type productRepo struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) repository.ProductRepository {
	return &productRepo{db: db}
}

// Create implements repository.ProductRepository.
func (p *productRepo) Create(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, stock, category_id, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	return p.db.QueryRowContext(
		ctx, query,
		product.ID, product.Name, product.Description, product.Price, product.Stock, product.CategoryID, nil, // Without images
	).Scan(&product.CreatedAt, &product.UpdatedAt)
}

// Delete implements repository.ProductRepository.
func (p *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id=$1`

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		return ErrProductNotFount
	}
	return nil
}

// GetByCategoryID implements repository.ProductRepository.
func (p *productRepo) GetByCategoryID(ctx context.Context, category_id uuid.UUID) (*entity.Product, error) {
	var product entity.Product

	query := `
		SELECT id, name, description, price, stock, category_id, image_url, created_at, updated_at
		FROM products
		WHERE category_id = $1
	`

	err := p.db.GetContext(ctx, &product, query, category_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by category_id: %w", err)
	}

	return &product, nil
}

// GetByID implements repository.ProductRepository.
func (p *productRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var product entity.Product

	query := `
		SELECT id, name, description, price, stock, category_id, image_url, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	err := p.db.GetContext(ctx, &product, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	return &product, nil
}

// GetByName implements repository.ProductRepository.
func (p *productRepo) GetByName(ctx context.Context, name string) (*entity.Product, error) {
	var product entity.Product

	query := `
		SELECT id, name, description, price, stock, category_id, image_url, created_at, updated_at
		FROM products
		WHERE name = $1
	`

	err := p.db.GetContext(ctx, &product, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by name: %w", err)
	}

	return &product, nil
}

// List implements repository.ProductRepository.
func (p *productRepo) List(ctx context.Context, filter repository.ProductFilter) ([]*entity.Product, int, error) {
	// 1. Формируем WHERE
	where := "WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if filter.CategoryID != nil {
		where += fmt.Sprintf(" AND category_id = $%d", argID)
		args = append(args, *filter.CategoryID)
		argID++
	}
	if filter.Search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argID, argID+1)
		val := "%" + filter.Search + "%"
		args = append(args, val, val)
		argID += 2
	}

	// 2. Считаем Total (без LIMIT/OFFSET)
	var total int
	countQuery := "SELECT count(*) FROM products " + where
	if err := p.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// 3. Получаем данные с пагинацией
	selectQuery := fmt.Sprintf(
		"SELECT id, name, description, price, stock, category_id, image_url, created_at, updated_at FROM products %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, argID, argID+1,
	)
	args = append(args, filter.Limit, filter.Offset)

	var products []*entity.Product
	if err := p.db.SelectContext(ctx, &products, selectQuery, args...); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update implements repository.ProductRepository.
func (p *productRepo) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE products
		SET name = $1,
			description = $2,
			price = $3,
			stock = $4,
			category_id = $5,
			image_url = $6,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`

	err := p.db.QueryRowContext(
		ctx, query,
		product.Name, product.Description, product.Price, product.Stock, product.CategoryID, product.Image_url, product.ID,
	).Scan(&product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// DecreaseStock implements repository.ProductRepository.
func (p *productRepo) DecreaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	query := `
		UPDATE products
		SET stock = stock - $1,
			updated_at = NOW()
		WHERE id = $2 AND stock >= $1
		RETURNING id
	`

	var productID uuid.UUID
	err := p.db.QueryRowContext(ctx, query, quantity, id).Scan(&productID)

	if err != nil {
		return fmt.Errorf("failed to decrease stock: %w", err)
	}

	return nil
}
