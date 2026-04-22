package category

import "errors"

var (
	ErrAdminOnly             = errors.New("forbidden: only admin")
	ErrCategoryAlreadyExists = errors.New("category with this name already exists")
	ErrInvalidName           = errors.New("invalid name")
)
