package validator

import "errors"

var (
	//user validator
	ErrInvalidEmail   = errors.New("invalid email format")
	ErrWeakPassword   = errors.New("password must be at least 8 characters")
	ErrEmptyFirstName = errors.New("first name cannot be empty")
	ErrEmptyLastName  = errors.New("last name cannot be empty")

	//product validator
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrEmptyDescription = errors.New("description cannot be empty")
	ErrNegValue         = errors.New("negative value")
)
