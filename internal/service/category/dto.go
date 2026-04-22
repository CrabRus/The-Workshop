package category

import (
	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/google/uuid"
)

type CategoryListResponse struct {
	Categories []*entity.Category `json:"categories"`
	Total      int                `json:"total"`
	Limit      int                `json:"limit"`
	Offset     int                `json:"offset"`
}

type CategoryRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	CategoryID  *uuid.UUID `json:"category_id"`
}

func ValidateCategoryCreate(req CategoryRequest) error {
	if req.Name == nil || *req.Name == "" {
		return ErrInvalidName
	}
	return nil
}
