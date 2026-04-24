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
// @Summary Get a list of categories
// @Description Retrieve a paginated list of product categories with optional searching
// @Tags categories
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} category.CategoryListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/categories [get]
func (h *categoryHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := categoryFilterFromRequest(r)

	resp, err := h.CategoryService.List(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// GET /categories/{id}
// @Summary Get category by ID
// @Description Retrieve category by UUID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID" format(uuid)
// @Success 200 {object} entity.Category
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/categories/{id} [get]
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
	r.Delete("/{id}", h.Delete)
}

// POST /admin/categories
// @Summary Create category
// @Description Create new category (admin only)
// @Tags admin-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body category.CategoryRequest true "Category data"
// @Success 201 {object} entity.Category
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/v1/admin/categories [post]
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

// DELETE /admin/categories/{id}
// @Summary Delete category
// @Description Delete category by ID (admin only)
// @Tags admin-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID" format(uuid)
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/v1/admin/categories/{id} [delete]
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

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Category deleted successfully",
	})
}
