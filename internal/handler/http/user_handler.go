package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/crabrus/the-workshop/internal/domain/repository"
	userSrv "github.com/crabrus/the-workshop/internal/service/user"
	"github.com/crabrus/the-workshop/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ---------- PUBLIC ----------

type userHandler struct {
	UserSrv userSrv.UserService
}

func NewUserHandler(srv userSrv.UserService) *userHandler {
	return &userHandler{UserSrv: srv}
}

func (h *userHandler) RegisterRoutes(r chi.Router) {
	// @Summary Get current user profile
	// @Description Retrieve the profile information of the authenticated user.
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} entity.User "User profile"
	// @Failure 401 {object} ErrorResponse "Unauthorized"
	// @Failure 404 {object} ErrorResponse "User not found"
	// @Router /api/v1/users/me [get]
	r.Get("/me", h.GetProfile)

	// @Summary Update current user profile
	// @Description Update the profile information of the authenticated user.
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body userSrv.UpdateProfileRequest true "Updated user profile data"
	// @Success 200 {object} entity.User "User profile updated successfully"
	// @Failure 400 {object} ErrorResponse "Invalid input"
	// @Failure 401 {object} ErrorResponse "Unauthorized"
	// @Failure 404 {object} ErrorResponse "User not found"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /api/v1/users/me [put]
	r.Put("/me", h.UpdateMe)
}

func (h *userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.UserSrv.GetProfile(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *userHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req userSrv.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	user, err := h.UserSrv.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// ---------- ADMIN ----------

type adminUserHandler struct {
	UserSrv userSrv.UserService
}

func NewAdminUserHandler(srv userSrv.UserService) *adminUserHandler {
	return &adminUserHandler{UserSrv: srv}
}

func (h *adminUserHandler) RegisterRoutes(r chi.Router) {
	// @Summary Search and list users (Admin only)
	// @Description Retrieve a paginated list of users with optional searching and filtering. Requires admin role.
	// @Tags Admin - Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param search query string false "Search term for user first name, last name, or email"
	// @Param limit query int false "Number of items to return" default(20)
	// @Param offset query int false "Number of items to skip" default(0)
	// @Success 200 {object} userSrv.UserListResponse "List of users"
	// @Failure 401 {object} ErrorResponse "Unauthorized"
	// @Failure 403 {object} ErrorResponse "Forbidden"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	r.Get("/search", h.List)

	// @Summary Get user by ID (Admin only)
	// @Description Retrieve detailed information about a user by their ID. Requires admin role.
	// @Tags Admin - Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID" format(uuid)
	// @Success 200 {object} entity.User "User details"
	// @Failure 400 {object} ErrorResponse "Invalid user ID"
	// @Failure 401 {object} ErrorResponse "Unauthorized"
	// @Failure 403 {object} ErrorResponse "Forbidden"
	// @Failure 404 {object} ErrorResponse "User not found"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	r.Get("/{id}", h.GetUserByID)

	// @Summary Update user by ID (Admin only)
	// @Description Update the profile information of a user by their ID. Requires admin role.
	// @Tags Admin - Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID" format(uuid)
	// @Param request body userSrv.UpdateProfileRequest true "Updated user data"
	// @Success 200 {object} entity.User "User updated successfully"
	// @Failure 400 {object} ErrorResponse "Invalid input or user ID"
	// @Failure 401 {object} ErrorResponse "Unauthorized"
	// @Failure 403 {object} ErrorResponse "Forbidden"
	// @Failure 404 {object} ErrorResponse "User not found"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /api/v1/admin/users/{id} [put]
	r.Put("/{id}", h.UpdateUser)

	// @Summary Delete user by ID (Admin only)
	// @Description Delete a user by their ID. Requires admin role.
	// @Tags Admin - Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID" format(uuid)
	// @Success 200 {object} SuccessResponse "User deleted successfully"
	// @Failure 400 {object} ErrorResponse "Invalid user ID"
	// @Failure 401 {object} ErrorResponse "Unauthorized"
	// @Failure 403 {object} ErrorResponse "Forbidden"
	// @Failure 404 {object} ErrorResponse "User not found"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /api/v1/admin/users/{id} [delete]
	r.Delete("/{id}", h.DeleteUser)
}

func (h *adminUserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.UserSrv.GetProfile(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *adminUserHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := userFilterFromRequest(r)

	resp, err := h.UserSrv.List(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *adminUserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req userSrv.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	user, err := h.UserSrv.UpdateByAdmin(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *adminUserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.UserSrv.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Message: "User deleted successfully"})
}

// ---------- FILTER ----------

func userFilterFromRequest(r *http.Request) repository.UserFilter {
	q := r.URL.Query()

	filter := repository.UserFilter{}
	if search := q.Get("search"); search != "" {
		filter.Search = search
	}

	if limit, err := strconv.Atoi(q.Get("limit")); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.Atoi(q.Get("offset")); err == nil {
		filter.Offset = offset
	}

	return filter
}
