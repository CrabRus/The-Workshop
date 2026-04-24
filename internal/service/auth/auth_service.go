package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	repository "github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/crabrus/the-workshop/pkg/validator"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*TokenResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) bool
}

type authService struct {
	userRepo       repository.UserRepository
	jwtKey         string
	expiryDuration time.Duration
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		jwtKey = "your-secret-key-change-in-production"
	}

	return &authService{
		userRepo:       userRepo,
		jwtKey:         jwtKey,
		expiryDuration: 24 * time.Hour,
	}
}

// Login implements AuthService.
func (a *authService) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	err := validator.ValidateEmail(req.Email)
	if err != nil {
		return nil, err
	}
	// Find user by email
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Compare passwords
	if !a.ComparePassword(user.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, expiresIn, err := a.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &TokenResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   "Bearer",
	}, nil
}

// Register implements AuthService.
func (a *authService) Register(ctx context.Context, req RegisterRequest) (*TokenResponse, error) {
	err := ValidateRegister(req)
	if err != nil {
		return nil, err
	}
	// Check if user already exists
	existingUser, _ := a.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := a.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &entity.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "customer",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user to database
	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token
	token, expiresIn, err := a.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &TokenResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   "Bearer",
	}, nil
}

// ValidateToken implements AuthService.
func (a *authService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtKey), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken implements AuthService.
func (a *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Validate refresh token
	claims, err := a.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get user from database
	user, err := a.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new token
	token, expiresIn, err := a.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &TokenResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   "Bearer",
	}, nil
}

// HashPassword implements AuthService.
func (a *authService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// ComparePassword implements AuthService.
func (a *authService) ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// generateToken generates a JWT token for the user
func (a *authService) generateToken(user *entity.User) (string, int64, error) {
	expirationTime := time.Now().Add(a.expiryDuration)

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(a.expiryDuration.Seconds()), nil
}
