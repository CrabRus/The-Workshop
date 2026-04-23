package cart

import (
	"github.com/google/uuid"
)

// CartItemDTO - полная информация о товаре в корзине
type CartItemDTO struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductPrice float64   `json:"product_price"`
	Quantity     int       `json:"quantity"`
	Sum          float64   `json:"sum"` // price × quantity
}

// CartResponse - полная информация о корзине пользователя
type CartResponse struct {
	Items       []*CartItemDTO `json:"items"`
	TotalAmount float64        `json:"total_amount"`
	TotalItems  int            `json:"total_items"`
}
