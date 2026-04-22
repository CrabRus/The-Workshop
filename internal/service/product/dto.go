package product

import (
	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/crabrus/the-workshop/pkg/validator"
	"github.com/google/uuid"
)

// -------------------- RESPONSE --------------------

type ProductListResponse struct {
	Products []*entity.Product `json:"products"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// -------------------- REQUEST --------------------

// Create / Update DTO (один для обох, але з різною валідацією)
type ProductRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Price       *float64   `json:"price"`
	Stock       *int       `json:"stock"`
	CategoryID  *uuid.UUID `json:"category_id"`
	ImageURL    *string    `json:"image_url"`
}

// -------------------- VALIDATION (CREATE) --------------------

func ValidateProductCreate(req ProductRequest) error {
	if req.Name == nil || *req.Name == "" {
		return ErrInvalidName
	}

	if req.Price == nil || *req.Price <= 0 {
		return ErrInvalidPrice
	}

	if req.Stock == nil || *req.Stock < 0 {
		return ErrInvalidStock
	}

	if req.CategoryID == nil {
		return ErrInvalidCategory
	}

	return nil
}

// -------------------- VALIDATION (UPDATE) --------------------

func ValidateProductUpdate(req ProductRequest) error {
	if req.Name != nil {
		if err := validator.ValidateName(*req.Name); err != nil {
			return err
		}
	}

	if req.Price != nil {
		if err := validator.ValidatePrice(*req.Price); err != nil {
			return err
		}
	}

	if req.Stock != nil {
		if err := validator.ValidateStock(*req.Stock); err != nil {
			return err
		}
	}

	return nil
}
