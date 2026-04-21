package utils

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const ContextKeyUserID contextKey = "user_id"
const ContextKeyUserRole contextKey = "user_role"
const ContextKeyUserEmail contextKey = "user_email"

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	return userID, ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyUserRole).(string)
	return role, ok
}

func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ContextKeyUserEmail).(string)
	return email, ok
}
