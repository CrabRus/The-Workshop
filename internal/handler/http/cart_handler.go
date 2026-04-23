package http

import (
	"encoding/json"
	"net/http"

	"github.com/crabrus/the-workshop/internal/service/cart"
	"github.com/crabrus/the-workshop/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type cartHandler struct {
	CartService cart.CartService
}

func NewCartHandler(srv cart.CartService) *cartHandler {
	return &cartHandler{CartService: srv}
}

func (h *cartHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetCart)
	r.Post("/items", h.AddItem)
	r.Put("/items/{id}", h.UpdateItem)
	r.Delete("/items/{id}", h.RemoveItem)
	r.Delete("/", h.ClearCart)
}

// AddItemRequest - запрос на добавление товара в корзину
type AddItemRequest struct {
	ProductID *uuid.UUID `json:"product_id"`
	Quantity  *int       `json:"quantity"`
}

// UpdateItemRequest - запрос на обновление количества товара
type UpdateItemRequest struct {
	Quantity *int `json:"quantity"`
}

// GET /cart
func (h *cartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	resp, err := h.CartService.GetCart(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// POST /cart/items
func (h *cartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	// Валидация
	if req.ProductID == nil || req.Quantity == nil {
		respondError(w, http.StatusBadRequest, "product_id and quantity are required")
		return
	}

	if req.Quantity == nil || *req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "quantity must be greater than 0")
		return
	}

	res, err := h.CartService.AddItem(r.Context(), userID, *req.ProductID, *req.Quantity)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "product not found" {
			status = http.StatusNotFound
		}
		respondError(w, status, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, res)
}

// PUT /cart/items/{id}
func (h *cartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	cartItemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	// Валидация
	if req.Quantity == nil {
		respondError(w, http.StatusBadRequest, "quantity is required")
		return
	}

	if *req.Quantity < 0 {
		respondError(w, http.StatusBadRequest, "quantity cannot be negative")
		return
	}

	res, err := h.CartService.UpdateItem(r.Context(), userID, cartItemID, *req.Quantity)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "cart item not found":
			status = http.StatusNotFound
		case "unauthorized: cart item does not belong to this user":
			status = http.StatusForbidden
		}
		respondError(w, status, err.Error())
		return
	}

	// Если quantity было 0, то элемент удалён и мы возвращаем nil
	if res == nil {
		respondJSON(w, http.StatusOK, map[string]string{"status": "item removed"})
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// DELETE /cart/items/{id}
func (h *cartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	cartItemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.CartService.RemoveItem(r.Context(), userID, cartItemID)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "cart item not found":
			status = http.StatusNotFound
		case "unauthorized: cart item does not belong to this user":
			status = http.StatusForbidden
		}
		respondError(w, status, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "item removed"})
}

// DELETE /cart
func (h *cartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	err := h.CartService.Clear(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "cart cleared"})
}
