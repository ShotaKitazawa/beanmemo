package handler

import (
	"context"
	"fmt"
	"slices"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/auth"
)

const disabledOIDCUserID int64 = 1

type userRepository interface {
	UpsertBySub(ctx context.Context, sub, name string) error
	GetBySub(ctx context.Context, sub string) (int64, error)
}

// tokenVerifier verifies a JWT string and returns the claims.
type tokenVerifier interface {
	Verify(ctx context.Context, tokenString string) (auth.Claims, error)
}

// SecurityHandler implements api.SecurityHandler for bearer JWT authentication.
type SecurityHandler struct {
	verifier    tokenVerifier
	userRepo    userRepository
	disableOIDC bool
	allowedSubs []string
}

func NewSecurityHandler(verifier tokenVerifier, userRepo userRepository, disableOIDC bool, allowedSubs []string) *SecurityHandler {
	return &SecurityHandler{
		verifier:    verifier,
		userRepo:    userRepo,
		disableOIDC: disableOIDC,
		allowedSubs: allowedSubs,
	}
}

// HandleBearerAuth validates the JWT and stores the resolved user ID in context.
// When disableOIDC is true, all requests are treated as the default user (ID=1).
func (s *SecurityHandler) HandleBearerAuth(ctx context.Context, _ api.OperationName, t api.BearerAuth) (context.Context, error) {
	if s.disableOIDC {
		return auth.WithUserID(ctx, disabledOIDCUserID), nil
	}

	claims, err := s.verifier.Verify(ctx, t.Token)
	if err != nil {
		return ctx, fmt.Errorf("unauthorized: %w", err)
	}

	if len(s.allowedSubs) > 0 && !slices.Contains(s.allowedSubs, claims.Sub) {
		return ctx, fmt.Errorf("unauthorized: sub %q is not in OIDC_ALLOWED_SUBS", claims.Sub)
	}

	if err := s.userRepo.UpsertBySub(ctx, claims.Sub, claims.Sub); err != nil {
		return ctx, fmt.Errorf("upsert user: %w", err)
	}

	userID, err := s.userRepo.GetBySub(ctx, claims.Sub)
	if err != nil {
		return ctx, fmt.Errorf("get user: %w", err)
	}

	ctx = auth.WithToken(ctx, t.Token)
	return auth.WithUserID(ctx, userID), nil
}
