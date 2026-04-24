package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/crabrus/the-workshop/internal/domain/repository"
)

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, ErrorResponse{
		Error:  msg,
		Status: status,
	})
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func productFilterFromRequest(r *http.Request) repository.ProductFilter {
	q := r.URL.Query()

	filter := repository.ProductFilter{
		Search:  q.Get("search"),
		OrderBy: q.Get("order_by"),
	}

	// limit
	if limitStr := q.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// offset
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	// category_id (у тебе string*)
	if cat := q.Get("category_id"); cat != "" {
		filter.CategoryID = &cat
	}

	return filter
}

func categoryFilterFromRequest(r *http.Request) repository.CategoryFilter {
	q := r.URL.Query()

	filter := repository.CategoryFilter{
		Search:  q.Get("search"),
		OrderBy: q.Get("order_by"),
	}

	// limit
	if limitStr := q.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// offset
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}
	return filter
}
