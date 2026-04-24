# The Workshop - E-Commerce API

A comprehensive RESTful API for an e-commerce platform built with Go, PostgreSQL, and following Clean Architecture principles.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Database Schema](#database-schema)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)

---

## 🎯 Overview

The Workshop is a production-ready e-commerce API that provides complete functionality for managing:
- User authentication and authorization
- Product catalog with categories and filtering
- Shopping cart operations
- Order management with status tracking
- Admin dashboard with statistics and data export

Built with enterprise-grade patterns and best practices.

---

## ✨ Features

### User Management
✅ User registration with email validation
✅ JWT-based authentication with refresh tokens
✅ Role-based access control (customer, admin)
✅ User profile management
✅ Admin user blocking/unblocking

### Product Management
✅ Browse products with pagination
✅ Search and filter by category, price, availability
✅ Product details with inventory tracking
✅ Admin CRUD operations
✅ Category management

### Shopping Cart
✅ Add/remove items from cart
✅ Update quantities with stock validation
✅ Cart persistence for authenticated users
✅ Automatic total calculation

### Order Management
✅ Order creation from shopping cart
✅ Multiple payment methods (cash, card mock)
✅ Order status tracking (pending → delivered)
✅ Order history for users
✅ Order cancellation with stock restoration
✅ Admin order management and status updates

### Admin Dashboard
✅ Platform statistics (users, orders, revenue)
✅ Export data to CSV (orders, products, users)
✅ Order management and status updates
✅ User management and account control

---

## 🛠️ Tech Stack

- **Language**: Go 1.25+
- **Database**: PostgreSQL 15+
- **Router**: Chi v5
- **Authentication**: JWT (golang-jwt/jwt)
- **Database Driver**: sqlx with lib/pq
- **Validation**: go-playground/validator
- **Migrations**: golang-migrate
- **Documentation**: Swagger/OpenAPI 3.0
- **Containerization**: Docker & Docker Compose
- **Logging**: Go's standard slog package

---

## 📁 Project Structure

```
the-workshop/
├── cmd/
│   └── api/                    # Application entry point
│       └── main.go
│
├── internal/
│   ├── config/                 # Configuration management
│   │   └── config.go
│   │
│   ├── db/                     # Database initialization
│   │   └── db.go
│   │
│   ├── domain/                 # Business logic & entities
│   │   ├── entity/             # Domain models
│   │   ├── repository/         # Repository interfaces
│   │   └── service/            # Service interfaces (optional)
│   │
│   ├── repository/             # Data access layer
│   │   └── postgres/           # PostgreSQL implementation
│   │
│   ├── service/                # Business logic
│   │   ├── auth/
│   │   ├── user/
│   │   ├── product/
│   │   ├── category/
│   │   ├── cart/
│   │   └── order/
│   │
│   └── handler/                # HTTP handlers
│       └── http/               # REST endpoints
│
├── pkg/                        # Reusable packages
│   ├── logger/                 # Logging
│   ├── validator/              # Validation
│   └── utils/                  # Utilities
│
├── migrations/                 # Database migrations
├── docs/                       # Documentation
├── docker-compose.yml          # Docker Compose configuration
├── Dockerfile                  # Production Dockerfile
├── Dockerfile.dev              # Development Dockerfile
├── go.mod & go.sum             # Go dependencies
└── README.md                   # This file
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose
- PostgreSQL 15 (or use Docker)
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/crabrus/the-workshop.git
cd the-workshop
```

2. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Start Docker containers**
```bash
docker-compose up -d
```

4. **Apply database migrations**
```bash
docker-compose exec api migrate -path migrations -database "postgres://user:password@db:5432/dbname?sslmode=disable" up
```

5. **Seed database (optional)**
```bash
psql -U your_db_user -d your_db_name < scripts/seed.sql
```

6. **Run the application**
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

---

## ⚙️ Configuration

Configuration is managed through environment variables (`.env` file).

### Required Variables

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_DAYS=7

# Logging
LOG_LEVEL=info
LOG_JSON=false
```

### Optional Variables

```env
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis_password_123

ENVIRONMENT=production
DEBUG=false
```

---

## 📚 API Endpoints

### Authentication
```
POST   /api/v1/auth/register     # Register new user
POST   /api/v1/auth/login        # Login user
POST   /api/v1/auth/refresh      # Refresh access token
POST   /api/v1/auth/logout       # Logout user
```

### Products (Public)
```
GET    /api/v1/products          # List products
GET    /api/v1/products/{id}     # Get product details
GET    /api/v1/products?search=  # Search products
```

### Categories (Public)
```
GET    /api/v1/categories        # List categories
GET    /api/v1/categories/{id}   # Get category details
```

### Cart (Auth Required)
```
GET    /api/v1/cart              # View cart
POST   /api/v1/cart/items        # Add to cart
PUT    /api/v1/cart/items/{id}   # Update quantity
DELETE /api/v1/cart/items/{id}   # Remove from cart
DELETE /api/v1/cart              # Clear cart
```

### Orders (Auth Required)
```
POST   /api/v1/orders            # Create order
GET    /api/v1/orders            # Get user's orders
GET    /api/v1/orders/{id}       # Get order details
DELETE /api/v1/orders/{id}       # Cancel order
```

### Users (Auth Required)
```
GET    /api/v1/users/me          # Get current profile
PUT    /api/v1/users/me          # Update profile
```

### Admin
```
GET    /api/v1/admin/statistics  # Platform statistics
GET    /api/v1/admin/users/search  # List users
PUT    /api/v1/admin/users/{id}  # Update user
DELETE /api/v1/admin/users/{id}  # Delete user
PUT    /api/v1/admin/users/{id}/block   # Block user
PUT    /api/v1/admin/users/{id}/unblock # Unblock user
GET    /api/v1/admin/orders      # List all orders
PUT    /api/v1/admin/orders/{id}/status # Update order status
POST   /api/v1/admin/export/orders   # Export orders to CSV
POST   /api/v1/admin/export/products # Export products to CSV
POST   /api/v1/admin/export/users    # Export users to CSV
```

**Full API documentation**: See [HANDLERS.md](docs/HANDLERS.md)

---

## 🗄️ Database Schema

### Tables

- `users` - User accounts and authentication
- `categories` - Product categories
- `products` - Product catalog with inventory
- `cart_items` - User shopping carts
- `orders` - Order records
- `order_items` - Items in each order

### Key Migrations

```
000001_create_users_table.sql
000002_create_categories_table.sql
000003_create_products_table.sql
000004_create_cart_items_table.sql
000005_create_orders_tables.sql
```

### ER Diagram
```
Users
  ├─ Carts (cart_items)
  └─ Orders
     └─ OrderItems
        └─ Products
           └─ Categories
```

---

## 👨‍💻 Development

### Project Structure Rationale

- **Domain Layer** (`internal/domain/`): Contains business logic and is independent
- **Service Layer** (`internal/service/`): Implements use cases
- **Repository Layer** (`internal/repository/`): Data access abstraction
- **Handler Layer** (`internal/handler/`): HTTP request/response handling

This follows **Clean Architecture** principles for maintainability and testability.

### Adding New Features

1. Create entity in `internal/domain/entity/`
2. Create repository interface in `internal/domain/repository/`
3. Implement repository in `internal/repository/postgres/`
4. Create service in `internal/service/`
5. Create handler in `internal/handler/http/`
6. Register routes in `internal/handler/http/router.go`

### Code Style

- Follow Go conventions
- Use `gofmt` for formatting
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused

---

## 🧪 Testing

### Unit Tests

```bash
go test ./... -v
```

### Test Coverage

```bash
go test -cover ./...
```

### Integration Tests

```bash
go test -tags=integration ./...
```

### Running Specific Tests

```bash
go test ./internal/service/product -v -run TestListProducts
```

---

## 🐳 Docker

### Build Docker Image

```bash
docker build -t the-workshop:latest .
```

### Run with Docker Compose

```bash
docker-compose up -d
docker-compose logs -f api
```

### Stop Containers

```bash
docker-compose down
```

### Remove All Data

```bash
docker-compose down -v
```

---

## 📊 Database Migrations

### Create New Migration

```bash
migrate create -ext sql -dir migrations create_new_table
```

### Apply Migrations

```bash
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up
```

### Rollback Migrations

```bash
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" down
```

### Check Migration Status

```bash
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" version
```

---

## 📖 API Documentation

### Swagger/OpenAPI

1. Install swag:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate documentation:
```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

3. View docs at: `http://localhost:8080/swagger/index.html`

---

## 🔐 Security Considerations

- Passwords are hashed with bcrypt
- JWT tokens with expiration
- CORS headers configured
- SQL injection prevention via parameterized queries
- Input validation on all endpoints
- Admin endpoints require role verification

### Production Checklist

- [ ] Change JWT_SECRET to a strong random value
- [ ] Set ENVIRONMENT=production
- [ ] Enable HTTPS/SSL
- [ ] Use strong database password
- [ ] Enable database backups
- [ ] Set up monitoring and logging
- [ ] Configure rate limiting
- [ ] Enable CORS restrictions

---

## 📈 Performance

- Database indexing on frequently queried fields
- Pagination support for large datasets
- Connection pooling via sqlx
- Efficient queries with selective field retrieval

### Optimization Tips

- Use pagination for list endpoints
- Filter results by relevant fields
- Cache frequently accessed data (consider Redis)
- Monitor slow queries

---

## 🐛 Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose ps

# View logs
docker-compose logs postgres

# Verify connection string in .env
```

### Port Already in Use

```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or use different port
export SERVER_PORT=8081
```

### Swagger Not Generating

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Clean and regenerate
rm -rf docs
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

---

## 📝 Logging

The API uses Go's standard `slog` package for structured logging.

### Log Levels

- `debug` - Detailed diagnostic information
- `info` - General information about application flow
- `warn` - Warning messages for potentially harmful situations
- `error` - Error messages for failed operations

### Viewing Logs

```bash
# Live logs
docker-compose logs -f api

# Filtered logs
docker-compose logs api | grep ERROR
```

---

## 🚢 Deployment

### Cloud Platforms

#### AWS
- ECS for container orchestration
- RDS for PostgreSQL
- CloudFront for CDN

#### Heroku
```bash
heroku create your-app-name
heroku addons:create heroku-postgresql:hobby-dev
git push heroku main
```

#### Docker Hub
```bash
docker build -t your-username/the-workshop:latest .
docker push your-username/the-workshop:latest
```

---

## 📚 Additional Resources

- [Go Best Practices](https://golang.org/doc/effective_go)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT Authentication](https://tools.ietf.org/html/rfc7519)
- [REST API Design](https://restfulapi.net/)

---

## 👥 Team

Developed as a final project for "Advanced Go Development" course.

---

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

---

## 📞 Support

For issues, questions, or contributions:
- Open an issue on GitHub
- Contact the development team
- Check the documentation in `/docs`

---

## 🎉 Acknowledgments

- Go community and libraries
- Database design best practices
- Clean Architecture principles
- Course instructors and mentors

---

**Last Updated**: April 24, 2026

**Current Version**: 1.0.0

**Status**: Production Ready ✅
