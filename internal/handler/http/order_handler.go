package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/crabrus/the-workshop/internal/service/order"
	"github.com/crabrus/the-workshop/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ---------- PUBLIC ----------

type orderHandler struct {
	OrderService order.OrderService
}

func NewOrderHandler(srv order.OrderService) *orderHandler {
	return &orderHandler{OrderService: srv}
}

func (h *orderHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.CreateOrder)
	r.Get("/", h.GetUserOrders)
	r.Get("/{id}", h.GetOrder)
	r.Delete("/{id}", h.CancelOrder)
}

// POST /orders
// @Summary Create a new order
// @Description Create a new order from the user's shopping cart. Requires authentication.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body order.CreateOrderRequest true "Order creation data"
// @Success 201 {object} order.OrderDTO "Order created successfully"
// @Failure 400 {object} ErrorResponse "Invalid input, empty cart, or invalid payment method"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Product not found (if stock check fails)"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/orders [post]
func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	var req order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.OrderService.CreateOrder(r.Context(), userID, req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == order.ErrOrderNotFound.Error() {
			status = http.StatusNotFound
		} else if err.Error() == order.ErrEmptyCart.Error() {
			status = http.StatusBadRequest
		} else if err.Error() == order.ErrInvalidPaymentMethod.Error() {
			status = http.StatusBadRequest
		}
		respondError(w, status, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, res)
}

// GET /orders
// @Summary Get user's orders
// @Description Retrieve a paginated list of orders for the current user. Requires authentication.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter orders by status (pending, confirmed, shipped, delivered, cancelled)"
// @Param limit query int false "Number of items to return" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} order.OrderListResponse "List of user's orders"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
func (h *orderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	filter := orderFilterFromRequest(r)

	resp, err := h.OrderService.GetUserOrders(r.Context(), userID, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// GET /orders/{id}
// @Summary Get order by ID
// @Description Retrieve detailed information about a specific order for the current user. Requires authentication.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID" format(uuid)
// @Success 200 {object} order.OrderDTO "Order details"
// @Failure 400 {object} ErrorResponse "Invalid order ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden: order does not belong to user"
// @Failure 404 {object} ErrorResponse "Order not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
func (h *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := h.OrderService.GetOrder(r.Context(), userID, orderID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == order.ErrOrderNotFound.Error() {
			status = http.StatusNotFound
		} else if err.Error() == order.ErrUnauthorized.Error() {
			status = http.StatusForbidden
		}
		respondError(w, status, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// DELETE /orders/{id}
// @Summary Cancel an order
// @Description Cancel a specific order for the current user and return items to stock. Requires authentication.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID" format(uuid)
// @Success 200 {object} SuccessResponse "Order cancelled successfully"
// @Failure 400 {object} ErrorResponse "Invalid order ID, order already cancelled, or cannot cancel delivered/shipped order"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden: order does not belong to user"
// @Failure 404 {object} ErrorResponse "Order not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/orders/{id} [delete]
func (h *orderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.OrderService.CancelOrder(r.Context(), userID, orderID); err != nil {
		status := http.StatusBadRequest
		if err.Error() == order.ErrOrderNotFound.Error() {
			status = http.StatusNotFound
		} else if err.Error() == order.ErrUnauthorized.Error() {
			status = http.StatusForbidden
		}
		respondError(w, status, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Message: "Order cancelled successfully"})
}

// ---------- ADMIN ----------

type adminOrderHandler struct {
	OrderService order.OrderService
}

func NewAdminOrderHandler(srv order.OrderService) *adminOrderHandler {
	return &adminOrderHandler{OrderService: srv}
}

func (h *adminOrderHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.ListAllOrders)
	r.Put("/{id}/status", h.UpdateOrderStatus)
	r.Get("/statistics", h.GetStatistics)
}

// GET /admin/orders
// @Summary Get all orders (Admin only)
// @Description Retrieve a paginated list of all orders in the system. Requires admin role.
// @Tags Admin - Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter orders by status (pending, confirmed, shipped, delivered, cancelled)"
// @Param limit query int false "Number of items to return" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} order.OrderListResponse "List of all orders"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
func (h *adminOrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	filter := orderFilterFromRequest(r)

	resp, err := h.OrderService.GetAllOrders(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// PUT /admin/orders/{id}/status
// @Summary Update order status (Admin only)
// @Description Update the status of a specific order by its ID. Requires admin role.
// @Tags Admin - Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID" format(uuid)
// @Param request body order.UpdateOrderStatusRequest true "New order status"
// @Success 200 {object} SuccessResponse "Order status updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid order ID, status missing, or invalid status"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Order not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/orders/{id}/status [put]
func (h *adminOrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req order.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if req.Status == nil || *req.Status == "" {
		respondError(w, http.StatusBadRequest, "status is required")
		return
	}

	if err := h.OrderService.UpdateOrderStatus(r.Context(), orderID, *req.Status); err != nil {
		status := http.StatusBadRequest
		if err.Error() == order.ErrOrderNotFound.Error() {
			status = http.StatusNotFound
		} else if err.Error() == order.ErrInvalidStatus.Error() {
			status = http.StatusBadRequest
		}
		respondError(w, status, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Message: "Order status updated successfully"})
}

// GET /admin/orders/statistics
// @Summary Get e-commerce statistics (Admin only)
// @Description Retrieve aggregated statistics such as total orders, total revenue, and top-selling products. Requires admin role.
// @Tags Admin - Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} order.StatisticsDTO "E-commerce statistics"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
func (h *adminOrderHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.OrderService.GetStatistics(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// ---------- HELPERS ----------

func orderFilterFromRequest(r *http.Request) repository.OrderFilter {
	query := r.URL.Query()

	limit := 20
	if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
		limit = l
	}

	offset := 0
	if o, err := strconv.Atoi(query.Get("offset")); err == nil && o >= 0 {
		offset = o
	}

	status := query.Get("status")

	return repository.OrderFilter{
		Status: status,
		Limit:  limit,
		Offset: offset,
	}
}
