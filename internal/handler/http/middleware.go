package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/crabrus/the-workshop/internal/service/auth"
	"github.com/crabrus/the-workshop/pkg/utils"
)

func RequireAuth(authSrv auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "Invalid auth header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, http.StatusUnauthorized, "Invalid auth header format")
				return
			}

			tokenString := parts[1]

			claims, err := authSrv.ValidateToken(tokenString)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, utils.ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, utils.ContextKeyUserRole, claims.Role)
			ctx = context.WithValue(ctx, utils.ContextKeyUserEmail, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(utils.ContextKeyUserRole).(string)
		if !ok {
			respondError(w, http.StatusUnauthorized, "User not authorized")
			return
		}

		if role != "admin" {
			respondError(w, http.StatusForbidden, "Admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
