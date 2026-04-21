package product

import "errors"

var (
	ErrAdminOnly            = errors.New("forbidden: only admin")
	ErrProductAlreadyExists = errors.New("product with this name already exists")
	ErrInvalidName          = errors.New("invalid name")
	ErrInvalidPrice         = errors.New("invalid price amount")
	ErrInvalidStock         = errors.New("invalid stock amount")
	ErrInvalidCategory      = errors.New("invalid category")
)
