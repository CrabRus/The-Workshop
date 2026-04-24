package order

import "errors"

var (
	ErrInvalidOrderData      = errors.New("invalid order data")
	ErrOrderNotFound         = errors.New("order not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrEmptyCart             = errors.New("cart is empty")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrInvalidStatus         = errors.New("invalid status")
	ErrInsufficientStock     = errors.New("insufficient stock")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrInvalidAddressData    = errors.New("invalid address data")
	ErrAdminOnly             = errors.New("admin only")
)
