package http

import (
	"net/http"

	authService "github.com/crabrus/the-workshop/internal/service/auth"
	userService "github.com/crabrus/the-workshop/internal/service/user"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type RouterConfig struct {
	AuthService authService.AuthService
	UserService userService.UserService
}

func NewRouter(config RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	// Global middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(CORS)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		authHandler := NewAuthHandler(config.AuthService)
		authHandler.RegisterRoutes(r)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))

			// User routes
			userHandler := NewUserHandler(config.UserService)
			userHandler.RegisterRoutes(r)
		})

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))
			r.Use(RequireAdmin)

			// 👉 тут будуть admin handlers (products, orders, users...)
			// приклад:
			// adminHandler := NewAdminHandler(...)
			// adminHandler.RegisterRoutes(r)
		})
	})

	return r
}
