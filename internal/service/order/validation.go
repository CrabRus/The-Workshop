package order

import (
	"github.com/crabrus/the-workshop/internal/domain/entity"
)

func ValidateCreateOrder(req CreateOrderRequest) error {
	if req.ShippingAddress.FullName == nil || *req.ShippingAddress.FullName == "" {
		return ErrInvalidAddressData
	}

	if req.ShippingAddress.PhoneNumber == nil || *req.ShippingAddress.PhoneNumber == "" {
		return ErrInvalidAddressData
	}

	if req.ShippingAddress.Email == nil || *req.ShippingAddress.Email == "" {
		return ErrInvalidAddressData
	}

	if req.ShippingAddress.Country == nil || *req.ShippingAddress.Country == "" {
		return ErrInvalidAddressData
	}

	if req.ShippingAddress.City == nil || *req.ShippingAddress.City == "" {
		return ErrInvalidAddressData
	}

	if req.ShippingAddress.PostalCode == nil || *req.ShippingAddress.PostalCode == "" {
		return ErrInvalidAddressData
	}

	if req.ShippingAddress.AddressLine == nil || *req.ShippingAddress.AddressLine == "" {
		return ErrInvalidAddressData
	}

	if req.PaymentMethod == nil || *req.PaymentMethod == "" {
		return ErrInvalidPaymentMethod
	}

	if *req.PaymentMethod != entity.PaymentMethodCash && *req.PaymentMethod != entity.PaymentMethodCard {
		return ErrInvalidPaymentMethod
	}

	return nil
}
