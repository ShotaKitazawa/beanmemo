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

func authenticatedCtx(userID int64, token string) context.Context {
	ctx := auth.WithUserID(context.Background(), userID)
	if token != "" {
		ctx = auth.WithToken(ctx, token)
	}
	return ctx
}

func TestGetUserinfo_Success(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{
		result: auth.UserinfoResult{
			Sub:     "user|123",
			Name:    "Alice",
			Email:   "alice@example.com",
			Picture: "https://example.com/alice.jpg",
		},
	})

	ctx := authenticatedCtx(1, "tok")
	res, err := h.GetUserinfo(ctx)
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

func TestGetUserinfo_Unauthorized(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{})

	res, err := h.GetUserinfo(context.Background()) // no userID in context
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(*api.GetUserinfoUnauthorized); !ok {
		t.Fatalf("expected *api.GetUserinfoUnauthorized, got %T", res)
	}
}

func TestGetUserinfo_ProviderError(t *testing.T) {
	h := handler.New(nil, nil, &mockUserinfoProvider{
		err: errors.New("upstream failed"),
	})

	ctx := authenticatedCtx(1, "tok")
	res, err := h.GetUserinfo(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.(*api.GetUserinfoInternalServerError); !ok {
		t.Fatalf("expected *api.GetUserinfoInternalServerError, got %T", res)
	}
}
