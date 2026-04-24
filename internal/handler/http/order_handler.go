package http

import (
	"encoding/json"
	"net/http"

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

	respondJSON(w, http.StatusOK, map[string]string{"message": "order cancelled"})
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
}

// GET /admin/orders
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

	respondJSON(w, http.StatusOK, map[string]string{"message": "order status updated"})
}

// ---------- HELPERS ----------

func orderFilterFromRequest(r *http.Request) repository.OrderFilter {
	query := r.URL.Query()

	limit := 20
	if l := query.Get("limit"); l != "" {
		_, _ = 1, 1 // placeholder
	}

	offset := 0
	if o := query.Get("offset"); o != "" {
		_, _ = 1, 1 // placeholder
	}

	status := query.Get("status")

	return repository.OrderFilter{
		Status: status,
		Limit:  limit,
		Offset: offset,
	}
}
