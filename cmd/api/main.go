package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crabrus/the-workshop/internal/db"
	"github.com/crabrus/the-workshop/internal/repository/postgres"

	httpHandler "github.com/crabrus/the-workshop/internal/handler/http"
	"github.com/crabrus/the-workshop/internal/service/auth"
	"github.com/crabrus/the-workshop/internal/service/cart"
	"github.com/crabrus/the-workshop/internal/service/category"
	"github.com/crabrus/the-workshop/internal/service/order"
	"github.com/crabrus/the-workshop/internal/service/product"
	"github.com/crabrus/the-workshop/internal/service/user"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment varibles")
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "ecommerce_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	serverPort := getEnv("SERVER_PORT", "8080")
	//jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")

	dbConfig := db.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  dbSSLMode,
	}

	database, err := db.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Database connection established")
	userRepo := postgres.NewUserRepository(database)
	productRepo := postgres.NewProductRepository(database)
	categoryRepo := postgres.NewCategoryRepo(database)
	cartItemsRepo := postgres.NewCartItemRepository(database)
	orderRepo := postgres.NewOrderRepository(database)
	orderItemRepo := postgres.NewOrderItemRepository(database)

	authService := auth.NewAuthService(userRepo)
	userService := user.NewService(userRepo)
	productService := product.NewService(productRepo)
	categoryService := category.NewService(categoryRepo)
	cartService := cart.NewService(cartItemsRepo, productService)
	orderService := order.NewService(orderRepo, orderItemRepo, cartService, productService, productRepo)

	router := httpHandler.NewRouter(httpHandler.RouterConfig{
		AuthService:     authService,
		UserService:     userService,
		ProductService:  productService,
		CartService:     cartService,
		CategoryService: categoryService,
		OrderService:    orderService,
	})

	server := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on http://localhost:%s", serverPort)
		log.Printf("API documentation: http://localhost:%s/api/v1", serverPort)
		log.Printf("Health check: http://localhost:%s/health", serverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server....")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
