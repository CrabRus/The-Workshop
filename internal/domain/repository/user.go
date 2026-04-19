package repository

import (
	"context"

	entity "github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter UserFilter) ([]*entity.User, int, error)
}

type UserFilter struct {
	Search  string
	Role    *string
	Limit   int
	Offset  int
	OrderBy string
}
