package auth

import "context"

type contextKey struct{}

// WithUserID returns a new context with the given user ID stored.
func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// UserIDFromContext retrieves the user ID stored by WithUserID.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(contextKey{}).(int64)
	return id, ok
}
