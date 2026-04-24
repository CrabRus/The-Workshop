package validator

import (
	"net/mail"
	"unicode"
	"unicode/utf8"
)

func ValidateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	if addr.Address != email {
		return ErrInvalidEmail
	}
	if len(email) > 254 {
		return ErrInvalidEmail
	}

	return nil
}

func ValidatePassword(password string) error {
	var (
		hasMinLen  = utf8.RuneCountInString(password) >= 8
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if !hasMinLen {
		return ErrPasswordTooShort
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrWeakPassword
	}

	return nil
}

func ValidateFirstName(name string) error {
	if name == "" {
		return ErrEmptyFirstName
	}
	return nil
}

func ValidateLastName(name string) error {
	if name == "" {
		return ErrEmptyLastName
	}
	return nil
}
