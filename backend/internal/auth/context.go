package auth

import "context"

type contextKey struct{}
type tokenContextKey struct{}

// WithUserID returns a new context with the given user ID stored.
func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// UserIDFromContext retrieves the user ID stored by WithUserID.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(contextKey{}).(int64)
	return id, ok
}

// WithToken returns a new context with the raw access token stored.
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenContextKey{}, token)
}

// TokenFromContext retrieves the access token stored by WithToken.
func TokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenContextKey{}).(string)
	return token, ok && token != ""
}
