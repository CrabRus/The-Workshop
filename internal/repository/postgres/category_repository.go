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
	CategoryNotFount = errors.New("category not found")
)

type categoryRepo struct {
	db *database.DB
}

// GetByName implements repository.CategoryRepository.
func (c *categoryRepo) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	var category entity.Category

	query := `
		SELECT id, name, description, created_at,
		FROM categories
		WHERE name = $1
	`

	err := c.db.GetContext(ctx, &category, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by name: %w", err)
	}

	return &category, nil
}

func NewCategoryRepo(db *database.DB) repository.CategoryRepository {
	return &categoryRepo{db: db}
}

// Create implements repository.CategoryRepository.
func (c *categoryRepo) Create(ctx context.Context, category *entity.Category) error {
	query := `
		INSERT INTO categories (id, name, description)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`

	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}

	return c.db.QueryRowContext(
		ctx, query,
		category.ID, category.Name, category.Description,
	).Scan(&category.CreatedAt)
}

// Delete implements repository.CategoryRepository.
func (c *categoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id=$1`

	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		return CategoryNotFount
	}
	return nil
}

// GetByID implements repository.CategoryRepository.
func (c *categoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var category entity.Category

	query := `
		SELECT id, name, description, created_at
		FROM categories
		WHERE id = $1
	`

	err := c.db.GetContext(ctx, &category, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}

	return &category, nil
}

// List implements repository.CategoryRepository.
func (c *categoryRepo) List(ctx context.Context, filter repository.CategoryFilter) ([]*entity.Category, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if filter.Search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argID, argID+1)
		val := "%" + filter.Search + "%"
		args = append(args, val, val)
		argID += 2
	}

	// 2. Считаем Total (без LIMIT/OFFSET)
	var total int
	countQuery := "SELECT count(*) FROM categories " + where
	if err := c.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// 3. Получаем данные с пагинацией
	selectQuery := fmt.Sprintf(
		"SELECT id, name, description, created_at FROM categories %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, argID, argID+1,
	)
	args = append(args, filter.Limit, filter.Offset)

	var categories []*entity.Category
	if err := c.db.SelectContext(ctx, &categories, selectQuery, args...); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
