package handler

import (
	"context"
	"fmt"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/auth"
)

const disabledOIDCUserID int64 = 1

type userRepository interface {
	UpsertBySub(ctx context.Context, sub, name string) error
	GetBySub(ctx context.Context, sub string) (int64, error)
}

// SecurityHandler implements api.SecurityHandler for bearer JWT authentication.
type SecurityHandler struct {
	verifier    *auth.JWTVerifier
	userRepo    userRepository
	disableOIDC bool
}

func NewSecurityHandler(verifier *auth.JWTVerifier, userRepo userRepository, disableOIDC bool) *SecurityHandler {
	return &SecurityHandler{
		verifier:    verifier,
		userRepo:    userRepo,
		disableOIDC: disableOIDC,
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

	if err := s.userRepo.UpsertBySub(ctx, claims.Sub, claims.Sub); err != nil {
		return ctx, fmt.Errorf("upsert user: %w", err)
	}

	userID, err := s.userRepo.GetBySub(ctx, claims.Sub)
	if err != nil {
		return ctx, fmt.Errorf("get user: %w", err)
	}

	return auth.WithUserID(ctx, userID), nil
}
