package validator

import (
	"regexp"
	"unicode/utf8"
)

func ValidateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
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
