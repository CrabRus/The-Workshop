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

// AddItemRequest - request to add an item to the cart
type AddItemRequest struct {
	ProductID *uuid.UUID `json:"product_id"`
	Quantity  *int       `json:"quantity"`
}

// UpdateItemRequest - request to update an item in the cart
type UpdateItemRequest struct {
	Quantity *int `json:"quantity"`
}

// GET /cart
// @Summary Get user's shopping cart
// @Description Retrieve the current user's shopping cart contents and total amount. Requires authentication.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} cart.CartResponse "Cart contents"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/cart [get]
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
// @Summary Add item to cart
// @Description Add a product to the current user's shopping cart. Requires authentication.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddItemRequest true "Product ID and quantity to add"
// @Success 201 {object} cart.CartItemDTO "Item added to cart"
// @Failure 400 {object} ErrorResponse "Invalid input, product_id or quantity missing, or insufficient stock"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/cart/items [post]
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
// @Summary Update item quantity in cart
// @Description Update the quantity of a specific item in the current user's shopping cart. Set quantity to 0 to remove. Requires authentication.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID" format(uuid)
// @Param request body UpdateItemRequest true "New quantity for the item"
// @Success 200 {object} cart.CartItemDTO "Item quantity updated"
// @Success 200 {object} SuccessResponse "Item removed if quantity is 0"
// @Failure 400 {object} ErrorResponse "Invalid input, quantity missing, or insufficient stock"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden: cart item does not belong to user"
// @Failure 404 {object} ErrorResponse "Cart item or product not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/cart/items/{id} [put]
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

	// Validation
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

	if res == nil {
		respondJSON(w, http.StatusOK, SuccessResponse{Message: "Item removed"})
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// DELETE /cart/items/{id}
// @Summary Remove item from cart
// @Description Remove a specific item from the current user's shopping cart. Requires authentication.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID" format(uuid)
// @Success 200 {object} SuccessResponse "Item removed from cart"
// @Failure 400 {object} ErrorResponse "Invalid cart item ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden: cart item does not belong to user"
// @Failure 404 {object} ErrorResponse "Cart item not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/cart/items/{id} [delete]
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

	respondJSON(w, http.StatusOK, SuccessResponse{Message: "Item removed from cart"})
}

// DELETE /cart
// @Summary Clear user's shopping cart
// @Description Remove all items from the current user's shopping cart. Requires authentication.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse "Cart cleared successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/cart [delete]
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
