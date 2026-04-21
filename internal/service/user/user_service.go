package user

import (
	"context"
	"fmt"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req UpdateProfileRequest, isAdmin bool) (*entity.User, error)
	List(ctx context.Context, filter repository.UserFilter) (*UserListResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) UserService {
	return &service{userRepo: userRepo}
}

func (s *service) GetProfile(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (s *service) UpdateProfile(ctx context.Context, id uuid.UUID, req UpdateProfileRequest, isAdmin bool) (*entity.User, error) {

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// ❗ централізована валідація DTO
	if err := ValidateUpdateProfile(req); err != nil {
		return nil, err
	}

	// EMAIL
	if req.Email != nil {
		user.Email = *req.Email
	}

	// PASSWORD
	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	// FIRST NAME
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	// LAST NAME
	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	// ROLE (тільки адмін)
	if req.Role != nil {
		if !isAdmin {
			return nil, fmt.Errorf("forbidden: only admin can change role")
		}
		user.Role = *req.Role
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *service) List(ctx context.Context, filter repository.UserFilter) (*UserListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return &UserListResponse{
		Users:  users,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
