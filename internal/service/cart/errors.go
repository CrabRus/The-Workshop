package cart

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock available")
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrInvalidQuantity   = errors.New("invalid quantity: must be greater than 0")
	ErrUnauthorized      = errors.New("unauthorized: cart item does not belong to this user")
)
