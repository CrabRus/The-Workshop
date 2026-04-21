package validator

func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyFirstName
	}
	return nil
}

func ValidatePrice(price float64) error {
	if price <= 0 {
		return ErrNegValue
	}
	return nil
}

func ValidateStock(stock int) error {
	if stock < 0 {
		return ErrNegValue
	}
	return nil
}
