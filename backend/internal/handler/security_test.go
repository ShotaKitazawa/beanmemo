package handler_test

import (
	"context"
	"testing"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/auth"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/handler"
)

type mockTokenVerifier struct {
	sub string
	err error
}

func (m *mockTokenVerifier) Verify(_ context.Context, _ string) (auth.Claims, error) {
	return auth.Claims{Sub: m.sub}, m.err
}

type mockUserRepo struct{}

func (m *mockUserRepo) UpsertBySub(_ context.Context, _, _ string) error { return nil }
func (m *mockUserRepo) GetBySub(_ context.Context, _ string) (int64, error) {
	return 1, nil
}

func TestHandleBearerAuth_AllowedSubs_Match(t *testing.T) {
	h := handler.NewSecurityHandler(
		&mockTokenVerifier{sub: "user|123"},
		&mockUserRepo{},
		false,
		[]string{"user|123", "user|456"},
	)

	ctx, err := h.HandleBearerAuth(context.Background(), api.GetUserinfoOperation, api.BearerAuth{Token: "tok"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if _, ok := auth.UserIDFromContext(ctx); !ok {
		t.Error("expected userID in context")
	}
}

func TestHandleBearerAuth_AllowedSubs_NoMatch(t *testing.T) {
	h := handler.NewSecurityHandler(
		&mockTokenVerifier{sub: "user|999"},
		&mockUserRepo{},
		false,
		[]string{"user|123"},
	)

	_, err := h.HandleBearerAuth(context.Background(), api.GetUserinfoOperation, api.BearerAuth{Token: "tok"})
	if err == nil {
		t.Fatal("expected error for sub not in allowed list")
	}
}

func TestHandleBearerAuth_AllowedSubs_Empty_AllowsAll(t *testing.T) {
	h := handler.NewSecurityHandler(
		&mockTokenVerifier{sub: "any-user"},
		&mockUserRepo{},
		false,
		nil, // no restriction
	)

	_, err := h.HandleBearerAuth(context.Background(), api.GetUserinfoOperation, api.BearerAuth{Token: "tok"})
	if err != nil {
		t.Fatalf("expected no error when allowedSubs is empty, got: %v", err)
	}
}

func TestHandleBearerAuth_DisableOIDC_SkipsCheck(t *testing.T) {
	h := handler.NewSecurityHandler(
		nil,
		&mockUserRepo{},
		true,
		[]string{"only-this-sub"},
	)

	ctx, err := h.HandleBearerAuth(context.Background(), api.GetUserinfoOperation, api.BearerAuth{Token: "tok"})
	if err != nil {
		t.Fatalf("expected no error in disable-oidc mode, got: %v", err)
	}
	if _, ok := auth.UserIDFromContext(ctx); !ok {
		t.Error("expected userID in context")
	}
}
