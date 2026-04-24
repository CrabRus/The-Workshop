package http

import (
	"net/http"

	authService "github.com/crabrus/the-workshop/internal/service/auth"
	cartService "github.com/crabrus/the-workshop/internal/service/cart"
	categoryService "github.com/crabrus/the-workshop/internal/service/category"
	orderService "github.com/crabrus/the-workshop/internal/service/order"
	productService "github.com/crabrus/the-workshop/internal/service/product"
	userService "github.com/crabrus/the-workshop/internal/service/user"

	_ "github.com/crabrus/the-workshop/docs"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type RouterConfig struct {
	AuthService     authService.AuthService
	UserService     userService.UserService
	ProductService  productService.ProductService
	CartService     cartService.CartService
	CategoryService categoryService.CategoryService
	OrderService    orderService.OrderService
}

func NewRouter(config RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	// ---------- GLOBAL MIDDLEWARES ----------
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(CORS)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

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

		//CATEGORIES (public)
		categoryHandler := NewCategoryHandler(config.CategoryService)
		r.Route("/categories", categoryHandler.RegisterRoutes)

		// ---------- AUTH (USER) ----------

		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))

			// USERS (/me)
			userHandler := NewUserHandler(config.UserService)
			r.Route("/users", userHandler.RegisterRoutes)

			// CART
			cartHandler := NewCartHandler(config.CartService)
			r.Route("/cart", cartHandler.RegisterRoutes)

			// ORDERS
			orderHandler := NewOrderHandler(config.OrderService)
			r.Route("/orders", orderHandler.RegisterRoutes)
		})

		// ---------- ADMIN ----------

		r.Route("/admin", func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))
			r.Use(RequireAdmin)

			// ----- ADMIN PRODUCTS -----
			adminProductHandler := NewAdminProductHandler(config.ProductService)
			r.Route("/products", adminProductHandler.RegisterRoutes)

			// ----- ADMIN USERS -----
			adminUserHandler := NewAdminUserHandler(config.UserService)
			r.Route("/users", adminUserHandler.RegisterRoutes)

			// ----- ADMIN CATEGORIES -----
			adminCategoryHandler := NewAdminCategoryHandler(config.CategoryService)
			r.Route("/categories", adminCategoryHandler.RegisterRoutes)

			// ----- ADMIN ORDERS -----
			adminOrderHandler := NewAdminOrderHandler(config.OrderService)
			r.Route("/orders", adminOrderHandler.RegisterRoutes)

			// ----- ADMIN STATISTICS & EXPORT -----
			adminHandler := NewAdminHandler(config.UserService, config.ProductService, config.OrderService)
			adminHandler.RegisterRoutes(r)
		})
	})

	return r
}
