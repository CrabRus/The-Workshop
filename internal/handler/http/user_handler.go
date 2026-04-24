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
	r.Get("/me", h.GetProfile)
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
	r.Get("/search", h.List)
	r.Get("/{id}", h.GetUserByID)
	r.Put("/{id}", h.UpdateUser)
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

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
	})
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
