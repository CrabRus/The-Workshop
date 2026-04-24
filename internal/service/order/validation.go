package order

import (
	"fmt"

	"github.com/crabrus/the-workshop/internal/domain/entity"
)

// validateShippingAddressField перевіряє, чи є покажчик на рядок nil або вказує на порожній рядок.
// Вона повертає більш конкретне повідомлення про помилку.
func validateShippingAddressField(field *string, fieldName string) error {
	if field == nil || *field == "" {
		return fmt.Errorf("%w: %s cannot be empty", ErrInvalidAddressData, fieldName)
	}
	return nil
}

func ValidateCreateOrder(req CreateOrderRequest) error {
	// Валідація полів адреси доставки
	if err := validateShippingAddressField(req.ShippingAddress.FullName, "Full Name"); err != nil {
		return err
	}
	if err := validateShippingAddressField(req.ShippingAddress.PhoneNumber, "Phone Number"); err != nil {
		return err
	}
	if err := validateShippingAddressField(req.ShippingAddress.Email, "Email"); err != nil {
		return err
	}
	if err := validateShippingAddressField(req.ShippingAddress.Country, "Country"); err != nil {
		return err
	}
	if err := validateShippingAddressField(req.ShippingAddress.City, "City"); err != nil {
		return err
	}
	if err := validateShippingAddressField(req.ShippingAddress.PostalCode, "Postal Code"); err != nil {
		return err
	}
	if err := validateShippingAddressField(req.ShippingAddress.AddressLine, "Address Line"); err != nil {
		return err
	}
	if req.PaymentMethod == nil || *req.PaymentMethod == "" {
		return ErrInvalidPaymentMethod
	}

	if *req.PaymentMethod != entity.PaymentMethodCash && *req.PaymentMethod != entity.PaymentMethodCard {
		return ErrInvalidPaymentMethod
	}

	return nil
}
