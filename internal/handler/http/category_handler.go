package http

import (
	"encoding/json"
	"net/http"

	"github.com/crabrus/the-workshop/internal/service/category"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ---------- PUBLIC ----------

type categoryHandler struct {
	CategoryService category.CategoryService
}

func NewCategoryHandler(srv category.CategoryService) *categoryHandler {
	return &categoryHandler{CategoryService: srv}
}

func (h *categoryHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/{id}", h.GetByID)
}

// GET /categories
func (h *categoryHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := categoryFilterFromRequest(r)

	resp, err := h.CategoryService.List(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// GET /products/{id}
func (h *categoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := h.CategoryService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// ---------- ADMIN ----------

type adminCategoryHandler struct {
	CategoryService category.CategoryService
}

func NewAdminCategoryHandler(srv category.CategoryService) *adminCategoryHandler {
	return &adminCategoryHandler{CategoryService: srv}
}

func (h *adminCategoryHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	// r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

// POST /admin/categories
func (h *adminCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req category.CategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.CategoryService.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, res)
}

// // PUT /admin/products/{id}
// func (h *adminCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	id, err := uuid.Parse(chi.URLParam(r, "id"))
// 	if err != nil {
// 		respondError(w, http.StatusBadRequest, "invalid id")
// 		return
// 	}

// 	var req category.CategoryRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		respondError(w, http.StatusBadRequest, "invalid body")
// 		return
// 	}

// 	res, err := h.CategoryService.Update(r.Context(), id, req)
// 	if err != nil {
// 		respondError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	respondJSON(w, http.StatusOK, res)
// }

// DELETE /admin/products/{id}
func (h *adminCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.CategoryService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
	})
}
