package auth

import (
	"context"

	"github.com/google/uuid"
)

type UserContextKey struct{}

// ContextWithUser добавляет ID пользователя в контекст.
func ContextWithUser(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserContextKey{}, userID)
}

// UserIDFromContext извлекает ID пользователя из контекста.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserContextKey{}).(uuid.UUID)
	return userID, ok
}
