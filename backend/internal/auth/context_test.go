package auth

import (
	"context"
	"testing"
)

func TestWithToken_RoundTrip(t *testing.T) {
	ctx := WithToken(context.Background(), "my-access-token")
	token, ok := TokenFromContext(ctx)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if token != "my-access-token" {
		t.Errorf("expected %q, got %q", "my-access-token", token)
	}
}

func TestTokenFromContext_Missing(t *testing.T) {
	_, ok := TokenFromContext(context.Background())
	if ok {
		t.Fatal("expected ok=false when no token is stored")
	}
}
