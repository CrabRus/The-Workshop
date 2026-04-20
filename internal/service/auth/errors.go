package auth

import "errors"

var (
	// Автентифікація
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")

	// Токени
	ErrTokenExpired        = errors.New("token expired")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrMissingAuthHeader   = errors.New("authorization header is missing")

	// Валідація полів
	ErrInvalidEmail      = errors.New("invalid email address format")
	ErrEmailRequired     = errors.New("email is required")
	ErrPasswordRequired  = errors.New("password is required")
	ErrInvalidPassword   = errors.New("invalid password") // для логіну
	ErrWeakPassword      = errors.New("password is too weak: must be at least 8 characters and include numbers")
	ErrFirstNameRequired = errors.New("first name is required")

	// Авторизація (Ролі)
	ErrPermissionDenied = errors.New("permission denied: insufficient privileges")
	ErrAdminOnly        = errors.New("this resource is restricted to administrators")

	// Системні
	ErrInternalServer = errors.New("an internal error occurred")
)
