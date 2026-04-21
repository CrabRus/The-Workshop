package http

import (
	"encoding/json"
	"net/http"

	userSrv "github.com/crabrus/the-workshop/internal/service/user"
	"github.com/crabrus/the-workshop/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	UserSrv userSrv.UserService
}

func NewUserHandler(srv userSrv.UserService) *userHandler {
	return &userHandler{UserSrv: srv}
}

func (h *userHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/me", h.GetProfile)
		r.Put("/me", h.UpdateUser)
	})
}

// GET /api/v1/users/me
func (h *userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	user, err := h.UserSrv.GetProfile(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// PUT /api/v1/users/me
func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	var req userSrv.UpdateProfileRequest

	// ❗ FIX: обов’язково перевірка decode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedUser, err := h.UserSrv.UpdateProfile(r.Context(), userID, req, false)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, updatedUser)
}
