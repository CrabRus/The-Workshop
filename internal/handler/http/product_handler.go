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
	// r.Get("/search", h.Search)
}

// GET /products
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

// // GET /products/search
// func (h *productHandler) Search(w http.ResponseWriter, r *http.Request) {
// 	filter := productFilterFromRequest(r)

// 	if filter.Search == "" {
// 		respondError(w, http.StatusBadRequest, "search query is required")
// 		return
// 	}

// 	resp, err := h.ProductService.List(r.Context(), filter)
// 	if err != nil {
// 		respondError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	respondJSON(w, http.StatusOK, resp)
// }

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

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
	})
}
