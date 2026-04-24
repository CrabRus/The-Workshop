package order

import (
	"time"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/google/uuid"
)

// CreateOrderRequest - запрос на создание заказа
type CreateOrderRequest struct {
	ShippingAddress ShippingAddressDTO `json:"shipping_address"`
	PaymentMethod   *string            `json:"payment_method"`
}

// ShippingAddressDTO - адрес доставки
type ShippingAddressDTO struct {
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
	Country     *string `json:"country"`
	City        *string `json:"city"`
	PostalCode  *string `json:"postal_code"`
	AddressLine *string `json:"address_line"`
}

// UpdateOrderStatusRequest - запрос на изменение статуса заказа
type UpdateOrderStatusRequest struct {
	Status *string `json:"status"`
}

// OrderDTO - полная информация о заказе
type OrderDTO struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.UUID              `json:"user_id"`
	Status          string                 `json:"status"`
	TotalAmount     float64                `json:"total_amount"`
	ShippingAddress entity.ShippingAddress `json:"shipping_address"`
	PaymentMethod   string                 `json:"payment_method"`
	Items           []*OrderItemDTO        `json:"items"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// OrderItemDTO - товар в заказе
type OrderItemDTO struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderListResponse - список заказов
type OrderListResponse struct {
	Orders []*OrderDTO `json:"orders"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}
