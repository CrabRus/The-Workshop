package user

import (
	entity "github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/crabrus/the-workshop/pkg/validator"
)

type UserListResponse struct {
	Users  []*entity.User
	Total  int
	Limit  int
	Offset int
}

// ❗ FIX: додані JSON теги
type UpdateProfileRequest struct {
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      *string `json:"role"`
}

// централізована валідація
func ValidateUpdateProfile(req UpdateProfileRequest) error {
	if req.Email != nil {
		if err := validator.ValidateEmail(*req.Email); err != nil {
			return err
		}
	}

	if req.Password != nil && *req.Password != "" {
		if err := validator.ValidatePassword(*req.Password); err != nil {
			return err
		}
	}

	if req.FirstName != nil {
		if err := validator.ValidateFirstName(*req.FirstName); err != nil {
			return err
		}
	}

	if req.LastName != nil {
		if err := validator.ValidateLastName(*req.LastName); err != nil {
			return err
		}
	}

	return nil
}
