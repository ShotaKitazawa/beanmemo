package handler

import (
	"context"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/auth"
)

// disabledOIDCUserinfoProvider returns a fixed dummy response when --disable-oidc is set.
type disabledOIDCUserinfoProvider struct{}

// NewDisabledOIDCUserinfoProvider returns a UserinfoProvider that always returns dummy data.
func NewDisabledOIDCUserinfoProvider() UserinfoProvider {
	return &disabledOIDCUserinfoProvider{}
}

func (d *disabledOIDCUserinfoProvider) FetchUserinfo(_ context.Context, _ string) (auth.UserinfoResult, error) {
	return auth.UserinfoResult{
		Sub:   "1",
		Name:  "dev",
		Email: "dev@beanmemo.local",
	}, nil
}
