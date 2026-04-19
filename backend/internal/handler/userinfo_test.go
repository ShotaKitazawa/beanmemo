package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/auth"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/handler"
)

type mockUserinfoProvider struct {
	result auth.UserinfoResult
	err    error
}

func (m *mockUserinfoProvider) FetchUserinfo(_ context.Context, _ string) (auth.UserinfoResult, error) {
	return m.result, m.err
}

// tokenCtx returns a context with only the Bearer token (no userID required for security:[]).
func tokenCtx(token string) context.Context {
	return auth.WithToken(context.Background(), token)
}

var enabledOIDC = handler.OIDCConfig{Enabled: true, Issuer: "https://idp.example.com", ClientID: "client", Audience: "aud"}
var disabledOIDC = handler.OIDCConfig{Enabled: false}

func TestGetUserinfo_Success(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{
		result: auth.UserinfoResult{
			Sub:     "user|123",
			Name:    "Alice",
			Email:   "alice@example.com",
			Picture: "https://example.com/alice.jpg",
		},
	}, enabledOIDC)

	res, err := h.GetUserinfo(tokenCtx("tok"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	typed, ok := res.(*api.UserinfoResponse)
	if !ok {
		t.Fatalf("expected *api.UserinfoResponse, got %T", res)
	}
	if typed.Sub != "user|123" {
		t.Errorf("expected sub 'user|123', got %q", typed.Sub)
	}
	if !typed.Name.Set || typed.Name.Null || typed.Name.Value != "Alice" {
		t.Errorf("unexpected Name: %+v", typed.Name)
	}
	if !typed.Email.Set || typed.Email.Null || typed.Email.Value != "alice@example.com" {
		t.Errorf("unexpected Email: %+v", typed.Email)
	}
}

func TestGetUserinfo_Unauthorized_WhenOIDCEnabledAndNoToken(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{}, enabledOIDC)

	res, err := h.GetUserinfo(context.Background()) // no token in context
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(*api.GetUserinfoUnauthorized); !ok {
		t.Fatalf("expected *api.GetUserinfoUnauthorized, got %T", res)
	}
}

func TestGetUserinfo_DisabledOIDC_NoTokenRequired(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{
		result: auth.UserinfoResult{Sub: "1", Name: "dev"},
	}, disabledOIDC)

	// With OIDC disabled, no token is needed → should return 200
	res, err := h.GetUserinfo(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(*api.UserinfoResponse); !ok {
		t.Fatalf("expected *api.UserinfoResponse, got %T", res)
	}
}

func TestGetUserinfo_ProviderError(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{
		err: errors.New("upstream failed"),
	}, enabledOIDC)

	res, err := h.GetUserinfo(tokenCtx("tok"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(*api.GetUserinfoInternalServerError); !ok {
		t.Fatalf("expected *api.GetUserinfoInternalServerError, got %T", res)
	}
}

func TestGetOidcConfig_Enabled(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{}, enabledOIDC)

	res, err := h.GetOidcConfig(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Enabled {
		t.Error("expected Enabled=true")
	}
	if res.Issuer.Value != "https://idp.example.com" {
		t.Errorf("unexpected Issuer: %+v", res.Issuer)
	}
}

func TestGetOidcConfig_Disabled(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{}, disabledOIDC)

	res, err := h.GetOidcConfig(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Enabled {
		t.Error("expected Enabled=false")
	}
	if !res.Issuer.Null {
		t.Errorf("expected Issuer to be null: %+v", res.Issuer)
	}
}
