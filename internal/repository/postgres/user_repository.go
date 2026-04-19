package postgres

import (
	"context"
	"fmt"

	database "github.com/crabrus/the-workshop/internal/db"
	entity "github.com/crabrus/the-workshop/internal/domain/entity"
	repository "github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/google/uuid"
)

type userRepo struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) repository.UserRepository {
	return &userRepo{db: db}
}

// Create implements UserRepository.
func (u *userRepo) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return u.db.QueryRowContext(
		ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// Delete implements UserRepository.
func (u *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id=$1`

	result, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// GetByEmail implements UserRepository.
func (u *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetByID implements UserRepository.
func (u *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User

	query := `
		SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// List implements UserRepository.
func (r *userRepo) List(ctx context.Context, f repository.UserFilter) ([]*entity.User, int, error) {
	// 1. Формируем WHERE
	where := "WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if f.Role != nil {
		where += fmt.Sprintf(" AND role = $%d", argID)
		args = append(args, *f.Role)
		argID++
	}
	if f.Search != "" {
		where += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d)", argID, argID+1, argID+2)
		val := "%" + f.Search + "%"
		args = append(args, val, val, val)
		argID += 3
	}

	// 2. Считаем Total (без LIMIT/OFFSET)
	var total int
	countQuery := "SELECT count(*) FROM users " + where
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// 3. Получаем данные с пагинацией
	selectQuery := fmt.Sprintf(
		"SELECT id, email, first_name, last_name, role, created_at, updated_at FROM users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, argID, argID+1,
	)
	args = append(args, f.Limit, f.Offset)

	var users []*entity.User
	if err := r.db.SelectContext(ctx, &users, selectQuery, args...); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update implements UserRepository.
func (u *userRepo) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET email = $1, 
		    password_hash = $2, 
		    first_name = $3, 
		    last_name = $4, 
		    role = $5, 
		    updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := u.db.QueryRowContext(
		ctx, query,
		user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role, user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
