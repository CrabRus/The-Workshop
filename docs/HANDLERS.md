# API Handlers Documentation

This document describes all HTTP handlers and endpoints for The Workshop e-commerce API.

## Table of Contents

1. [Authentication Handler](#authentication-handler)
2. [Product Handler](#product-handler)
3. [Category Handler](#category-handler)
4. [Cart Handler](#cart-handler)
5. [Order Handler](#order-handler)
6. [User Handler](#user-handler)
7. [Admin Handler](#admin-handler)

---

## Authentication Handler

**File**: `internal/handler/http/auth_handler.go`

Handles user registration, login, token refresh, and logout operations.

### Endpoints

#### POST /api/v1/auth/register
Register a new user account.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response** (201 Created):
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 86400
}
```

**Errors**:
- `400 Bad Request` - Invalid input
- `409 Conflict` - Email already registered
- `500 Internal Server Error` - Server error

---

#### POST /api/v1/auth/login
Authenticate user and return JWT tokens.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response** (200 OK):
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 86400
}
```

**Errors**:
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Wrong credentials
- `500 Internal Server Error` - Server error

---

#### POST /api/v1/auth/refresh
Get a new access token using refresh token.

**Request Body**:
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response** (200 OK):
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 86400
}
```

**Errors**:
- `400 Bad Request` - Missing refresh token
- `401 Unauthorized` - Invalid or expired token

---

#### POST /api/v1/auth/logout
Logout user (client-side token deletion).

**Response** (200 OK):
```json
{
  "message": "Logged out successfully"
}
```

---

## Product Handler

**File**: `internal/handler/http/product_handler.go`

Manages product viewing and searching. Includes admin CRUD operations.

### Public Endpoints

#### GET /api/v1/products
List all products with pagination and filtering.

**Query Parameters**:
- `search` (string) - Search in product name/description
- `category_id` (uuid) - Filter by category
- `limit` (int, default: 20) - Items per page
- `offset` (int, default: 0) - Pagination offset

**Response** (200 OK):
```json
{
  "products": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Product Name",
      "description": "Product description",
      "price": 99.99,
      "stock": 100,
      "category_id": "uuid",
      "created_at": "2026-04-24T10:00:00Z"
    }
  ],
  "total": 150,
  "limit": 20,
  "offset": 0
}
```

**Example with search**:
```bash
GET /api/v1/products?search=laptop&limit=10&offset=0
```

---

#### GET /api/v1/products/{id}
Get detailed information about a product.

**Parameters**:
- `id` (path, uuid) - Product ID

**Response** (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Product Name",
  "description": "Detailed description",
  "price": 99.99,
  "stock": 100,
  "category_id": "uuid",
  "created_at": "2026-04-24T10:00:00Z",
  "updated_at": "2026-04-24T10:00:00Z"
}
```

**Errors**:
- `400 Bad Request` - Invalid ID format
- `404 Not Found` - Product not found

---

### Admin Endpoints

#### POST /api/v1/admin/products
Create a new product (Admin only).

**Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Request Body**:
```json
{
  "name": "New Product",
  "description": "Product description",
  "price": 99.99,
  "stock": 100,
  "category_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response** (201 Created): Product object

**Errors**:
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Not authenticated
- `403 Forbidden` - Not admin
- `500 Internal Server Error`

---

#### PUT /api/v1/admin/products/{id}
Update product details (Admin only).

**Parameters**:
- `id` (path, uuid) - Product ID

**Request Body**: Same as POST

**Response** (200 OK): Updated product object

---

#### DELETE /api/v1/admin/products/{id}
Delete a product (Admin only).

**Parameters**:
- `id` (path, uuid) - Product ID

**Response** (200 OK):
```json
{
  "message": "Product deleted successfully"
}
```

---

## Category Handler

**File**: `internal/handler/http/category_handler.go`

Manages product categories.

### Public Endpoints

#### GET /api/v1/categories
List all categories with pagination.

**Query Parameters**:
- `search` (string) - Search in category name
- `limit` (int, default: 20)
- `offset` (int, default: 0)

**Response** (200 OK):
```json
{
  "categories": [
    {
      "id": "uuid",
      "name": "Electronics",
      "description": "Electronic products",
      "created_at": "2026-04-24T10:00:00Z"
    }
  ],
  "total": 5,
  "limit": 20,
  "offset": 0
}
```

---

#### GET /api/v1/categories/{id}
Get category details.

**Response** (200 OK): Category object

---

### Admin Endpoints

#### POST /api/v1/admin/categories
Create a new category (Admin only).

**Request Body**:
```json
{
  "name": "New Category",
  "description": "Category description"
}
```

---

#### DELETE /api/v1/admin/categories/{id}
Delete a category (Admin only).

---

## Cart Handler

**File**: `internal/handler/http/cart_handler.go`

Manages shopping cart operations.

**All endpoints require authentication** (Bearer token).

### Endpoints

#### GET /api/v1/cart
View current user's cart.

**Response** (200 OK):
```json
{
  "items": [
    {
      "id": "uuid",
      "product_id": "uuid",
      "product_name": "Product Name",
      "product_price": 99.99,
      "quantity": 2,
      "sum": 199.98
    }
  ],
  "total_amount": 199.98,
  "total_items": 2
}
```

---

#### POST /api/v1/cart/items
Add item to cart.

**Request Body**:
```json
{
  "product_id": "550e8400-e29b-41d4-a716-446655440000",
  "quantity": 2
}
```

**Response** (201 Created): Updated cart item

**Errors**:
- `400 Bad Request` - Invalid data or insufficient stock
- `404 Not Found` - Product not found

---

#### PUT /api/v1/cart/items/{id}
Update cart item quantity.

**Parameters**:
- `id` (path, uuid) - Cart item ID

**Request Body**:
```json
{
  "quantity": 5
}
```

**Response** (200 OK): Updated cart item

---

#### DELETE /api/v1/cart/items/{id}
Remove item from cart.

**Response** (200 OK):
```json
{
  "message": "Item removed from cart"
}
```

---

#### DELETE /api/v1/cart
Clear entire cart.

**Response** (200 OK):
```json
{
  "message": "Cart cleared"
}
```

---

## Order Handler

**File**: `internal/handler/http/order_handler.go`

Manages order creation, viewing, and cancellation.

**All endpoints require authentication** (Bearer token).

### Public Endpoints

#### POST /api/v1/orders
Create a new order from cart.

**Request Body**:
```json
{
  "shipping_address": {
    "full_name": "John Doe",
    "phone_number": "+38-095-123-45-67",
    "email": "john@example.com",
    "country": "Ukraine",
    "city": "Kyiv",
    "postal_code": "02000",
    "address_line": "vul. Khreshchatyk 25, apt 10"
  },
  "payment_method": "card"
}
```

**Response** (201 Created):
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "status": "pending",
  "total_amount": 1500.50,
  "shipping_address": { ... },
  "payment_method": "card",
  "items": [ ... ],
  "created_at": "2026-04-24T10:30:00Z",
  "updated_at": "2026-04-24T10:30:00Z"
}
```

**Errors**:
- `400 Bad Request` - Invalid data or empty cart
- `401 Unauthorized` - Not authenticated

---

#### GET /api/v1/orders
Get user's orders with filtering.

**Query Parameters**:
- `status` (string) - Filter: pending, confirmed, shipped, delivered, cancelled
- `limit` (int, default: 20)
- `offset` (int, default: 0)

**Response** (200 OK): List of orders

---

#### GET /api/v1/orders/{id}
Get specific order details.

**Parameters**:
- `id` (path, uuid) - Order ID

**Response** (200 OK): Order object

---

#### DELETE /api/v1/orders/{id}
Cancel an order (return items to stock).

**Response** (200 OK):
```json
{
  "message": "order cancelled"
}
```

**Errors**:
- `400 Bad Request` - Order cannot be cancelled
- `404 Not Found` - Order not found

---

### Admin Endpoints

#### GET /api/v1/admin/orders
Get all orders in system (Admin only).

**Query Parameters**:
- `status` (string) - Filter by status
- `limit` (int)
- `offset` (int)

**Response** (200 OK): List of all orders

---

#### PUT /api/v1/admin/orders/{id}/status
Update order status (Admin only).

**Request Body**:
```json
{
  "status": "shipped"
}
```

**Valid statuses**:
- `pending` - Awaiting confirmation
- `confirmed` - Confirmed
- `shipped` - Shipped
- `delivered` - Delivered
- `cancelled` - Cancelled

**Response** (200 OK):
```json
{
  "message": "order status updated"
}
```

---

## User Handler

**File**: `internal/handler/http/user_handler.go`

Manages user profile and admin user management.

### Public Endpoints (Auth Required)

#### GET /api/v1/users/me
Get current user's profile.

**Response** (200 OK):
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "customer",
  "created_at": "2026-04-24T10:00:00Z",
  "updated_at": "2026-04-24T10:00:00Z"
}
```

---

#### PUT /api/v1/users/me
Update current user's profile.

**Request Body**:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "newemail@example.com"
}
```

**Response** (200 OK): Updated user object

---

### Admin Endpoints

#### GET /api/v1/admin/users/search
Search and list users (Admin only).

**Query Parameters**:
- `search` (string) - Search in name or email
- `limit` (int)
- `offset` (int)

**Response** (200 OK):
```json
{
  "users": [ ... ],
  "total": 150,
  "limit": 20,
  "offset": 0
}
```

---

#### GET /api/v1/admin/users/{id}
Get user details (Admin only).

**Response** (200 OK): User object

---

#### PUT /api/v1/admin/users/{id}
Update user details (Admin only).

**Request Body**: User update data

**Response** (200 OK): Updated user object

---

#### DELETE /api/v1/admin/users/{id}
Delete user account (Admin only).

**Response** (200 OK):
```json
{
  "message": "User deleted successfully"
}
```

---

#### PUT /api/v1/admin/users/{id}/block
Block user account (Admin only).

**Response** (200 OK):
```json
{
  "message": "User blocked successfully"
}
```

---

#### PUT /api/v1/admin/users/{id}/unblock
Unblock user account (Admin only).

**Response** (200 OK):
```json
{
  "message": "User unblocked successfully"
}
```

---

## Admin Handler

**File**: `internal/handler/http/admin_handler.go`

Provides admin statistics, data export, and user management features.

**All endpoints require Admin role** (Bearer token with admin role).

### Endpoints

#### GET /api/v1/admin/statistics
Get platform statistics and metrics.

**Response** (200 OK):
```json
{
  "total_users": 150,
  "total_orders": 500,
  "total_revenue": 50000.00,
  "total_products": 200,
  "orders_pending": 50,
  "orders_shipped": 150,
  "orders_delivered": 300,
  "average_order_value": 100.00
}
```

---

#### POST /api/v1/admin/export/orders
Export all orders as CSV file (Admin only).

**Response** (200 OK): CSV file download
```csv
ID,User ID,Status,Total Amount,Payment Method,Created At
uuid,uuid,pending,1500.50,card,2026-04-24 10:30:00
...
```

---

#### POST /api/v1/admin/export/products
Export all products as CSV file (Admin only).

**Response** (200 OK): CSV file download
```csv
ID,Name,Price,Stock,Category ID,Created At
uuid,Product Name,99.99,100,uuid,2026-04-24 10:00:00
...
```

---

#### POST /api/v1/admin/export/users
Export all users as CSV file (Admin only).

**Response** (200 OK): CSV file download
```csv
ID,Email,First Name,Last Name,Role,Created At
uuid,user@example.com,John,Doe,customer,2026-04-24 10:00:00
...
```

---

## Authentication

Most endpoints require authentication using JWT Bearer tokens.

### Header Format
```
Authorization: Bearer <ACCESS_TOKEN>
```

### Token Types
- **Access Token**: Short-lived (24 hours), used for API requests
- **Refresh Token**: Long-lived (7 days), used to get new access tokens

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message",
  "status": 400
}
```

### Common HTTP Status Codes
- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Conflict (e.g., duplicate email)
- `500 Internal Server Error` - Server error

---

## Rate Limiting

Currently, there is no rate limiting implemented. This is recommended for production deployments.

---

## CORS Support

The API supports Cross-Origin Resource Sharing (CORS) for browser-based requests. CORS headers are sent with all responses.

---

## Best Practices

### 1. Always Include Bearer Token
```bash
curl -H "Authorization: Bearer <TOKEN>" https://api.example.com/api/v1/orders
```

### 2. Handle Pagination
```bash
# Get first 20 items
curl https://api.example.com/api/v1/products?limit=20&offset=0

# Get next 20
curl https://api.example.com/api/v1/products?limit=20&offset=20
```

### 3. Search Operations
```bash
# Search products
curl https://api.example.com/api/v1/products?search=laptop

# Filter by category
curl https://api.example.com/api/v1/products?category_id=uuid
```

### 4. Error Handling
Always check the HTTP status code and error message:
```json
{
  "error": "Email already registered",
  "status": 409
}
```

---

## Development

### Running the API
```bash
docker-compose up -d
```

### API Documentation
- Swagger UI: `http://localhost:8080/swagger/index.html`
- API Health: `http://localhost:8080/health`

### Testing Endpoints
Use Postman, Insomnia, or curl to test endpoints. Example:
```bash
curl -X GET http://localhost:8080/api/v1/products
```

---

## Support

For issues or questions, please contact the development team.
