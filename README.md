# 🛒 The Workshop - Professional E-commerce API

[![Go Version](https://img.shields.io/github/go-mod/go-version/crabrus/the-workshop)](https://golang.org)
[![Status](https://img.shields.io/badge/status-fully_completed-success)](https://github.com/crabrus/the-workshop)

A robust, production-ready RESTful API for an e-commerce platform built with **Go**, **PostgreSQL**, and **Clean Architecture**. This project provides a full-featured backend including user authentication, product catalog management, shopping cart functionality, order processing, and a comprehensive admin dashboard with data export capabilities.

## 🚀 Key Features

### 🔐 Authentication & Security
*   **JWT Auth**: Secure authentication using Access and Refresh tokens.
*   **RBAC**: Role-Based Access Control (Customer vs. Admin).
*   **Secure Hashing**: Password protection using Bcrypt.
*   **Validation**: Strict input validation for all requests.

### 📦 Product & Category Management
*   **Catalog**: Paginated listing with advanced search and filtering.
*   **Inventory**: Atomic stock control (Increase/Decrease) to prevent race conditions.
*   **Admin CRUD**: Full management of products and categories.

### 🛒 Shopping Experience
*   **Persistent Cart**: Saved cart items for authorized users.
*   **Checkout**: Seamless order creation from cart contents.
*   **Order Tracking**: Complete lifecycle from `pending` to `delivered` or `cancelled`.

### 📊 Admin Dashboard & Utilities
*   **Live Statistics**: Total revenue, average order value, and top-selling products.
*   **CSV Export**: One-click data export for Orders, Products, and Users.
*   **User Management**: Admin ability to search, update, and manage user accounts.
*   **Structured Logging**: Production-ready JSON logging using Go 1.21 `slog`.

## 🏗️ Architecture

The project follows **Clean Architecture** principles to ensure maintainability, scalability, and testability:

*   **Domain**: Business entities and repository/service interfaces (Framework independent).
*   **Service**: Implementation of business logic and DTO transformations.
*   **Repository**: Data access layer (PostgreSQL implementation using `sqlx`).
*   **Handler**: HTTP layer managing routing, middleware, and request/response parsing.

## 🛠️ Tech Stack

*   **Language**: Go 1.21+
*   **Database**: PostgreSQL 15+
*   **Router**: Chi (Lightweight and fast)
*   **Documentation**: Swagger / OpenAPI 2.0
*   **Containerization**: Docker & Docker Compose
*   **Config**: Environment-based configuration

## 🚦 Getting Started

### Prerequisites
*   Go 1.21 or higher
*   Docker and Docker Compose
*   Git

### Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/crabrus/the-workshop.git
   cd the-workshop
   ```

2. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env to set your DB credentials and JWT secret
   ```

3. **Run via Docker Compose**:
   ```bash
   docker-compose up -d
   ```

4. **Seed the database** (Optional):
   ```bash
   # Run the seed script to populate 20+ products and test users
   docker exec -i postgres_container psql -U postgres -d ecommerce_db < scripts/seed.sql
   ```

## 📖 API Documentation

Once the server is running, you can access the interactive API documentation:

*   **Swagger UI**: `http://localhost:8080/swagger/index.html`
*   **Health Check**: `http://localhost:8080/health`

### Primary Endpoints Summary

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Create a new account | Public |
| `POST` | `/api/v1/auth/login` | Get JWT tokens | Public |
| `GET` | `/api/v1/products` | List products (search/filter) | Public |
| `GET` | `/api/v1/cart` | View personal cart | User |
| `POST` | `/api/v1/orders` | Checkout from cart | User |
| `GET` | `/api/v1/admin/statistics` | Platform metrics | Admin |
| `POST` | `/api/v1/admin/export/orders`| Download Orders CSV | Admin |

## 📊 Admin Statistics & CSV Export

The API provides a powerful administration suite. Administrators can retrieve a real-time overview of the platform's performance, including:
*   **Total Revenue** (calculated automatically excluding cancelled orders).
*   **Top 5 Products** by sales volume.
*   **Orders Breakdown** by status.

**Example CSV Export:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/export/products \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  --output products_list.csv
```

## 🔍 Project Structure

```text
/cmd/api               # Application entry point
/internal
  /domain/entity       # Domain models
  /domain/repository   # Repository interfaces
  /service             # Business logic (Auth, Order, Product...)
  /repository/postgres # Database implementation
  /handler/http        # REST Controllers & Middlewares
  /config              # App configuration
/pkg                   # Shared utilities (Logger, Validator)
/scripts               # SQL seeds and helper scripts
/migrations            # SQL migration files
```

## 🧪 Testing

Run the full test suite including unit and repository tests:
```bash
go test -v ./...
```

## 📝 Future Roadmap
*   [ ] Redis implementation for product caching.
*   [ ] WebSocket notifications for order status updates.
*   [ ] Integration with a real Payment Gateway (Stripe/PayPal).
*   [ ] Prometheus metrics for monitoring.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---
*Developed as the final project for the "Advanced Go Development" module.*