package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID       `db:"id" json:"id"`
	UserID          uuid.UUID       `db:"user_id" json:"user_id"`
	Status          string          `db:"status" json:"status"` // pending, confirmed, shipped, delivered, cancelled
	TotalAmount     float64         `db:"total_amount" json:"total_amount"`
	ShippingAddress ShippingAddress `db:"shipping_address" json:"shipping_address"`
	PaymentMethod   string          `db:"payment_method" json:"payment_method"` // cash, card
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}

type ShippingAddress struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	AddressLine string `json:"address_line"`
}

// // Scan implements sql.Scanner
// func (sa *ShippingAddress) Scan(value interface{}) error {
// 	if value == nil {
// 		return nil
// 	}

// 	bytes, ok := value.([]byte)
// 	if !ok {
// 		return nil
// 	}

// 	return json.Unmarshal(bytes, sa)
// }

// Value implements driver.Valuer
func (sa ShippingAddress) Value() (driver.Value, error) {
	return json.Marshal(sa)
}

type OrderItem struct {
	ID        uuid.UUID `db:"id" json:"id"`
	OrderID   uuid.UUID `db:"order_id" json:"order_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	Price     float64   `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCancelled = "cancelled"

	PaymentMethodCash = "cash"
	PaymentMethodCard = "card"
)
