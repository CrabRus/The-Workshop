package http

import (
	"net/http"

	authService "github.com/crabrus/the-workshop/internal/service/auth"
	productService "github.com/crabrus/the-workshop/internal/service/product"
	userService "github.com/crabrus/the-workshop/internal/service/user"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type RouterConfig struct {
	AuthService    authService.AuthService
	UserService    userService.UserService
	ProductService productService.ProductService
}

func NewRouter(config RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	// ---------- GLOBAL MIDDLEWARES ----------
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(CORS)

	// ---------- HEALTH ----------
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// ---------- API v1 ----------
	r.Route("/api/v1", func(r chi.Router) {

		// ---------- PUBLIC ----------

		// AUTH
		authHandler := NewAuthHandler(config.AuthService)
		authHandler.RegisterRoutes(r)

		// PRODUCTS (public)
		productHandler := NewProductHandler(config.ProductService)
		r.Route("/products", productHandler.RegisterRoutes)

		// ---------- AUTH (USER) ----------

		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))

			// USERS (/me)
			userHandler := NewUserHandler(config.UserService)
			r.Route("/users", userHandler.RegisterRoutes)
		})

		// ---------- ADMIN ----------

		r.Route("/admin", func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))
			r.Use(RequireAdmin)

			// ----- ADMIN PRODUCTS -----
			adminProductHandler := NewAdminProductHandler(config.ProductService)
			r.Route("/products", adminProductHandler.RegisterRoutes)

			// ----- ADMIN USERS 🔥 -----
			adminUserHandler := NewAdminUserHandler(config.UserService)
			r.Route("/users", adminUserHandler.RegisterRoutes)

			// (майбутнє)
			// /admin/orders
			// /admin/categories
		})
	})

	return r
}
