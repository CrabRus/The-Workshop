# Quick Start Guide - The Workshop API

## 🚀 Starting the API

### Using Docker Compose (Recommended)
```bash
cd d:\TheWorkshop
docker-compose up -d
```

The API will be available at: `http://localhost:8080`

### Manual Start
```bash
# Install dependencies
go mod download

# Run the application
go run cmd/api/main.go
```

---

## 🌐 API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Documentation
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

---

## 🧪 Testing Endpoints

### 1. Authentication

#### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "SecurePass123!",
    "first_name": "Test",
    "last_name": "User"
  }'
```

**Response**:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 86400
}
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "SecurePass123!"
  }'
```

#### Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGc..."
  }'
```

#### Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout
```

---

### 2. Products

#### Get All Products
```bash
curl -X GET "http://localhost:8080/api/v1/products?limit=10"
```

#### Search Products
```bash
curl -X GET "http://localhost:8080/api/v1/products?search=laptop&limit=10"
```

#### Get Product by ID
```bash
curl -X GET "http://localhost:8080/api/v1/products/{PRODUCT_ID}"
```

#### Create Product (Admin)
```bash
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Product",
    "description": "Product description",
    "price": 99.99,
    "stock": 100,
    "category_id": "{CATEGORY_ID}"
  }'
```

#### Update Product (Admin)
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/products/{PRODUCT_ID}" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "price": 89.99,
    "stock": 50
  }'
```

#### Delete Product (Admin)
```bash
curl -X DELETE "http://localhost:8080/api/v1/admin/products/{PRODUCT_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

---

### 3. Categories

#### Get All Categories
```bash
curl -X GET "http://localhost:8080/api/v1/categories"
```

#### Get Category by ID
```bash
curl -X GET "http://localhost:8080/api/v1/categories/{CATEGORY_ID}"
```

---

### 4. Shopping Cart

#### View Cart
```bash
curl -X GET http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer {TOKEN}"
```

#### Add Item to Cart
```bash
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "{PRODUCT_ID}",
    "quantity": 2
  }'
```

#### Update Cart Item
```bash
curl -X PUT "http://localhost:8080/api/v1/cart/items/{CART_ITEM_ID}" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 5
  }'
```

#### Remove Item from Cart
```bash
curl -X DELETE "http://localhost:8080/api/v1/cart/items/{CART_ITEM_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

#### Clear Cart
```bash
curl -X DELETE http://localhost:8080/api/v1/cart \
  -H "Authorization: Bearer {TOKEN}"
```

---

### 5. Orders

#### Create Order
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_address": {
      "full_name": "John Doe",
      "phone_number": "+38-095-123-4567",
      "email": "john@example.com",
      "country": "Ukraine",
      "city": "Kyiv",
      "postal_code": "02000",
      "address_line": "vul. Khreshchatyk 25, apt 10"
    },
    "payment_method": "card"
  }'
```

#### Get User's Orders
```bash
curl -X GET "http://localhost:8080/api/v1/orders?limit=10" \
  -H "Authorization: Bearer {TOKEN}"
```

#### Get Specific Order
```bash
curl -X GET "http://localhost:8080/api/v1/orders/{ORDER_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

#### Cancel Order
```bash
curl -X DELETE "http://localhost:8080/api/v1/orders/{ORDER_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

---

### 6. User Profile

#### Get Current Profile
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer {TOKEN}"
```

#### Update Profile
```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "New First Name",
    "last_name": "New Last Name",
    "email": "newemail@example.com"
  }'
```

---

### 7. Admin - Users

#### List Users
```bash
curl -X GET "http://localhost:8080/api/v1/admin/users/search?limit=10" \
  -H "Authorization: Bearer {ADMIN_TOKEN}"
```

#### Search Users
```bash
curl -X GET "http://localhost:8080/api/v1/admin/users/search?search=john" \
  -H "Authorization: Bearer {ADMIN_TOKEN}"
```

#### Get User Details
```bash
curl -X GET "http://localhost:8080/api/v1/admin/users/{USER_ID}" \
  -H "Authorization: Bearer {ADMIN_TOKEN}"
```

#### Block User
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/users/{USER_ID}/block" \
  -H "Authorization: Bearer {ADMIN_TOKEN}"
```

#### Unblock User
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/users/{USER_ID}/unblock" \
  -H "Authorization: Bearer {ADMIN_TOKEN}"
```

---

### 8. Admin - Orders

#### List All Orders
```bash
curl -X GET "http://localhost:8080/api/v1/admin/orders?limit=20" \
  -H "Authorization: Bearer {ADMIN_TOKEN}"
```

#### Filter Orders by Status
```bash
curl -X GET "http://localhost:8080/api/v1/admin/orders?status=pending" \
  -H "Authorization: Bearer {ADMIN_TOKEN}"
```

#### Update Order Status
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/orders/{ORDER_ID}/status" \
  -H "Authorization: Bearer {ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "shipped"
  }'
```

---

### 9. Admin - Statistics

#### Get Platform Statistics
```bash
curl -X GET http://localhost:8080/api/v1/admin/statistics \
  -H "Authorization: Bearer {ADMIN_TOKEN}"
```

**Response**:
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

### 10. Admin - Export

#### Export Orders to CSV
```bash
curl -X POST http://localhost:8080/api/v1/admin/export/orders \
  -H "Authorization: Bearer {ADMIN_TOKEN}" \
  --output orders.csv
```

#### Export Products to CSV
```bash
curl -X POST http://localhost:8080/api/v1/admin/export/products \
  -H "Authorization: Bearer {ADMIN_TOKEN}" \
  --output products.csv
```

#### Export Users to CSV
```bash
curl -X POST http://localhost:8080/api/v1/admin/export/users \
  -H "Authorization: Bearer {ADMIN_TOKEN}" \
  --output users.csv
```

---

## 💡 Common Testing Workflow

### 1. Register and Login
```bash
# Register
REGISTER=$(curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "first_name": "Test",
    "last_name": "User"
  }')

# Extract token (use jq or similar)
TOKEN=$(echo $REGISTER | jq -r '.access_token')
```

### 2. Add Item to Cart
```bash
# Get a product first
PRODUCTS=$(curl -X GET http://localhost:8080/api/v1/products)
PRODUCT_ID=$(echo $PRODUCTS | jq -r '.products[0].id')

# Add to cart
curl -X POST http://localhost:8080/api/v1/cart/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"product_id\": \"$PRODUCT_ID\",
    \"quantity\": 1
  }"
```

### 3. Create Order
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_address": {
      "full_name": "Test User",
      "phone_number": "+38-095-123-4567",
      "email": "test@example.com",
      "country": "Ukraine",
      "city": "Kyiv",
      "postal_code": "02000",
      "address_line": "Test Street 123"
    },
    "payment_method": "card"
  }'
```

---

## 🔍 Using Postman/Insomnia

1. Import the API collection from Insomnia/Postman
2. Set up environment variables:
   - `BASE_URL`: http://localhost:8080/api/v1
   - `TOKEN`: (obtained from login response)
   - `ADMIN_TOKEN`: (for admin endpoints)

3. Use variables in requests:
   ```
   {{BASE_URL}}/products
   Header: Authorization: Bearer {{TOKEN}}
   ```

---

## 📊 Database Access

### Connect to PostgreSQL
```bash
# Using docker
docker-compose exec postgres psql -U postgres -d ecommerce_db

# Show tables
\dt

# Query users
SELECT id, email, role, created_at FROM users;

# Query products
SELECT id, name, price, stock FROM products;

# Query orders
SELECT id, user_id, status, total_amount FROM orders;
```

---

## 🐛 Troubleshooting

### Port Already in Use
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

### Database Connection Failed
```bash
# Check if containers are running
docker-compose ps

# View logs
docker-compose logs postgres
docker-compose logs api
```

### CORS Issues
- API has CORS middleware enabled
- Check `internal/handler/http/middleware.go`
- Adjust allowed origins if needed

---

## 📚 Additional Resources

- **Full API Documentation**: `/docs/HANDLERS.md`
- **README**: `/README.md`
- **Implementation Summary**: `/docs/IMPLEMENTATION_SUMMARY.md`
- **Swagger/OpenAPI**: http://localhost:8080/swagger/index.html

---

## 🎯 Testing Checklist

- [ ] Register new user
- [ ] Login with credentials
- [ ] Refresh access token
- [ ] Browse products
- [ ] Search products
- [ ] Add item to cart
- [ ] View cart
- [ ] Update cart item quantity
- [ ] Remove item from cart
- [ ] Create order
- [ ] View user orders
- [ ] Admin: View statistics
- [ ] Admin: Export orders to CSV
- [ ] Admin: List all orders
- [ ] Admin: Update order status
- [ ] Admin: Block/unblock user

---

## 🎓 Learning Resources

- Review the handler files to understand request/response patterns
- Check service files for business logic
- Study repository implementations for database queries
- Examine middleware for cross-cutting concerns

---

**Last Updated**: April 24, 2026
**Status**: Ready for Testing ✅
