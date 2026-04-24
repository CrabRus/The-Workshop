package http

import (
	"encoding/json"
	"net/http"

	"github.com/crabrus/the-workshop/internal/service/product"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ---------- PUBLIC ----------

type productHandler struct {
	ProductService product.ProductService
}

func NewProductHandler(srv product.ProductService) *productHandler {
	return &productHandler{ProductService: srv}
}

func (h *productHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/{id}", h.GetByID)
}

// GET /products
// @Summary Get a list of products
// @Description Retrieve a paginated list of products with optional filtering and searching
// @Tags Products
// @Accept json
// @Produce json
// @Param search query string false "Search term for product name or description"
// @Param category_id query string false "Filter by category ID"
// @Param limit query int false "Number of items to return" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} product.ProductListResponse "List of products"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/products [get]
func (h *productHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := productFilterFromRequest(r)

	resp, err := h.ProductService.List(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// GET /products/{id}
// @Summary Get product by ID
// @Description Retrieve detailed information about a single product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID" format(uuid)
// @Success 200 {object} product.ProductDTO "Product details"
// @Failure 400 {object} ErrorResponse "Invalid product ID"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Router /api/v1/products/{id} [get]
func (h *productHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := h.ProductService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// ---------- ADMIN ----------

type adminProductHandler struct {
	ProductService product.ProductService
}

func NewAdminProductHandler(srv product.ProductService) *adminProductHandler {
	return &adminProductHandler{ProductService: srv}
}

func (h *adminProductHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

// POST /admin/products
// @Summary Create a new product (Admin only)
// @Description Create a new product with details. Requires admin role.
// @Tags Admin - Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body product.ProductRequest true "Product data"
// @Success 201 {object} product.ProductDTO "Product created successfully"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/products [post]
func (h *adminProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req product.ProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.ProductService.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, res)
}

// PUT /admin/products/{id}
// @Summary Update an existing product (Admin only)
// @Description Update details of an existing product by its ID. Requires admin role.
// @Tags Admin - Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID" format(uuid)
// @Param request body product.ProductRequest true "Updated product data"
// @Success 200 {object} product.ProductDTO "Product updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid input or product ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/products/{id} [put]
func (h *adminProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req product.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.ProductService.Update(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// DELETE /admin/products/{id}
// @Summary Delete a product (Admin only)
// @Description Delete a product by its ID. Requires admin role.
// @Tags Admin - Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID" format(uuid)
// @Success 200 {object} SuccessResponse "Product deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid product ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/admin/products/{id} [delete]
func (h *adminProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.ProductService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Message: "Product deleted successfully"})
}
