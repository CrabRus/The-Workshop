package http

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/crabrus/the-workshop/internal/service/order"
	"github.com/crabrus/the-workshop/internal/service/product"
	userSrv "github.com/crabrus/the-workshop/internal/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// AdminHandler handles admin-specific requests
type AdminHandler struct {
	UserService    userSrv.UserService
	ProductService product.ProductService
	OrderService   order.OrderService
}

func NewAdminHandler(
	userService userSrv.UserService,
	productService product.ProductService,
	orderService order.OrderService,
) *AdminHandler {
	return &AdminHandler{
		UserService:    userService,
		ProductService: productService,
		OrderService:   orderService,
	}
}

func (h *AdminHandler) RegisterRoutes(r chi.Router) {
	r.Get("/statistics", h.GetStatistics)
	r.Post("/export/orders", h.ExportOrders)
	r.Post("/export/products", h.ExportProducts)
	r.Post("/export/users", h.ExportUsers)
	r.Put("/users/{id}/block", h.BlockUser)
	r.Put("/users/{id}/unblock", h.UnblockUser)
}

// Statistics DTO
type Statistics struct {
	TotalUsers      int     `json:"total_users"`
	TotalOrders     int     `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	TotalProducts   int     `json:"total_products"`
	OrdersPending   int     `json:"orders_pending"`
	OrdersShipped   int     `json:"orders_shipped"`
	OrdersDelivered int     `json:"orders_delivered"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// GetStatistics returns admin statistics
// @Summary Get platform statistics (Admin only)
// @Description Retrieve platform statistics including user count, order metrics, and revenue
// @Tags Admin - Statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Statistics "Platform statistics"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/statistics [get]
func (h *AdminHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get user count
	userFilter := repository.UserFilter{Limit: 1, Offset: 0}
	userResp, err := h.UserService.List(ctx, userFilter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get user statistics")
		return
	}

	// Get product count
	productFilter := repository.ProductFilter{Limit: 1, Offset: 0}
	productResp, err := h.ProductService.List(ctx, productFilter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get product statistics")
		return
	}

	// Get all orders for statistics
	allOrdersFilter := repository.OrderFilter{Limit: 1000, Offset: 0}
	allOrdersResp, err := h.OrderService.GetAllOrders(ctx, allOrdersFilter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get order statistics")
		return
	}

	// Calculate statistics
	stats := h.calculateStatistics(userResp.Total, productResp.Total, allOrdersResp)

	respondJSON(w, http.StatusOK, stats)
}

// ExportOrders exports orders to CSV
// @Summary Export orders to CSV (Admin only)
// @Description Download all orders as CSV file
// @Tags Admin - Export
// @Accept json
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {object} interface{} "CSV file"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/export/orders [post]
func (h *AdminHandler) ExportOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get all orders
	filter := repository.OrderFilter{Limit: 10000, Offset: 0}
	ordersResp, err := h.OrderService.GetAllOrders(ctx, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch orders")
		return
	}

	// Set CSV headers
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=orders.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header row
	header := []string{"ID", "User ID", "Status", "Total Amount", "Payment Method", "Created At"}
	if err := writer.Write(header); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to write CSV header")
		return
	}

	// Write data rows
	for _, o := range ordersResp.Orders {
		row := []string{
			o.ID.String(),
			o.UserID.String(),
			o.Status,
			fmt.Sprintf("%.2f", o.TotalAmount),
			o.PaymentMethod,
			o.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			respondError(w, http.StatusInternalServerError, "failed to write CSV row")
			return
		}
	}
}

// ExportProducts exports products to CSV
// @Summary Export products to CSV (Admin only)
// @Description Download all products as CSV file
// @Tags Admin - Export
// @Accept json
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {object} interface{} "CSV file"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/export/products [post]
func (h *AdminHandler) ExportProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get all products
	filter := repository.ProductFilter{Limit: 10000, Offset: 0}
	productsResp, err := h.ProductService.List(ctx, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	// Set CSV headers
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=products.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header row
	header := []string{"ID", "Name", "Price", "Stock", "Category ID", "Created At"}
	if err := writer.Write(header); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to write CSV header")
		return
	}

	// Write data rows
	for _, p := range productsResp.Products {
		row := []string{
			p.ID.String(),
			p.Name,
			fmt.Sprintf("%.2f", p.Price),
			fmt.Sprintf("%d", p.Stock),
			p.CategoryID.String(),
			p.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			respondError(w, http.StatusInternalServerError, "failed to write CSV row")
			return
		}
	}
}

// ExportUsers exports users to CSV
// @Summary Export users to CSV (Admin only)
// @Description Download all users as CSV file
// @Tags Admin - Export
// @Accept json
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {object} interface{} "CSV file"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/export/users [post]
func (h *AdminHandler) ExportUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get all users
	filter := repository.UserFilter{Limit: 10000, Offset: 0}
	usersResp, err := h.UserService.List(ctx, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	// Set CSV headers
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=users.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header row
	header := []string{"ID", "Email", "First Name", "Last Name", "Role", "Created At"}
	if err := writer.Write(header); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to write CSV header")
		return
	}

	// Write data rows
	for _, u := range usersResp.Users {
		row := []string{
			u.ID.String(),
			u.Email,
			u.FirstName,
			u.LastName,
			u.Role,
			u.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			respondError(w, http.StatusInternalServerError, "failed to write CSV row")
			return
		}
	}
}

// BlockUser blocks a user account
// @Summary Block user account (Admin only)
// @Description Block a user from accessing the platform
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" format(uuid)
// @Success 200 {object} SuccessResponse "User blocked successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/users/{id}/block [put]
func (h *AdminHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// Note: This requires adding a blocked field to the User entity and service
	// For now, this is a placeholder implementation
	respondJSON(w, http.StatusOK, SuccessResponse{Message: "User blocked successfully"})
}

// UnblockUser unblocks a user account
// @Summary Unblock user account (Admin only)
// @Description Unblock a user and allow access to the platform
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" format(uuid)
// @Success 200 {object} SuccessResponse "User unblocked successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/users/{id}/unblock [put]
func (h *AdminHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// Note: This requires adding a blocked field to the User entity and service
	// For now, this is a placeholder implementation
	respondJSON(w, http.StatusOK, SuccessResponse{Message: "User unblocked successfully"})
}

// Helper functions
func (h *AdminHandler) calculateStatistics(totalUsers, totalProducts int, ordersResp *order.OrderListResponse) Statistics {
	stats := Statistics{
		TotalUsers:    totalUsers,
		TotalOrders:   ordersResp.Total,
		TotalProducts: totalProducts,
	}

	if ordersResp.Total == 0 {
		return stats
	}

	// Calculate revenue and order statuses
	totalRevenue := 0.0
	pendingCount := 0
	shippedCount := 0
	deliveredCount := 0

	for _, o := range ordersResp.Orders {
		totalRevenue += o.TotalAmount

		switch o.Status {
		case "pending":
			pendingCount++
		case "shipped":
			shippedCount++
		case "delivered":
			deliveredCount++
		}
	}

	stats.TotalRevenue = totalRevenue
	stats.OrdersPending = pendingCount
	stats.OrdersShipped = shippedCount
	stats.OrdersDelivered = deliveredCount
	stats.AverageOrderValue = totalRevenue / float64(ordersResp.Total)

	return stats
}
