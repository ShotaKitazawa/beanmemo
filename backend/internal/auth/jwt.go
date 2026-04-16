package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Claims holds the verified claims from a JWT.
type Claims struct {
	Sub string
}

// UserinfoResult holds user info returned by the OIDC userinfo endpoint.
type UserinfoResult struct {
	Sub     string
	Name    string
	Email   string
	Picture string
}

// JWTVerifier validates JWT tokens using OIDC JWKS.
type JWTVerifier struct {
	issuerURL        string
	claimKey         string
	claimValue       string
	jwksCache        *jwk.Cache
	jwksURI          string
	userinfoEndpoint string
}

type oidcDiscovery struct {
	JWKSURI          string `json:"jwks_uri"`
	Issuer           string `json:"issuer"`
	UserinfoEndpoint string `json:"userinfo_endpoint"`
}

// NewJWTVerifier creates a JWTVerifier by fetching the OIDC discovery document
// and setting up a JWKS cache.
func NewJWTVerifier(ctx context.Context, issuerURL, claimKey, claimValue string) (*JWTVerifier, error) {
	discoveryURL := strings.TrimRight(issuerURL, "/") + "/.well-known/openid-configuration"
	resp, err := http.Get(discoveryURL) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("fetch OIDC discovery: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var disc oidcDiscovery
	if err := json.NewDecoder(resp.Body).Decode(&disc); err != nil {
		return nil, fmt.Errorf("decode OIDC discovery: %w", err)
	}

	cache := jwk.NewCache(ctx)
	if err := cache.Register(disc.JWKSURI, jwk.WithMinRefreshInterval(15*time.Minute)); err != nil {
		return nil, fmt.Errorf("register JWKS: %w", err)
	}
	if _, err := cache.Refresh(ctx, disc.JWKSURI); err != nil {
		return nil, fmt.Errorf("initial JWKS fetch: %w", err)
	}

	return &JWTVerifier{
		issuerURL:        issuerURL,
		claimKey:         claimKey,
		claimValue:       claimValue,
		jwksCache:        cache,
		jwksURI:          disc.JWKSURI,
		userinfoEndpoint: disc.UserinfoEndpoint,
	}, nil
}

// FetchUserinfo calls the OIDC userinfo endpoint with the given access token.
func (v *JWTVerifier) FetchUserinfo(ctx context.Context, accessToken string) (UserinfoResult, error) {
	if v.userinfoEndpoint == "" {
		return UserinfoResult{}, fmt.Errorf("userinfo_endpoint not present in OIDC discovery document")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.userinfoEndpoint, nil)
	if err != nil {
		return UserinfoResult{}, fmt.Errorf("create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UserinfoResult{}, fmt.Errorf("fetch userinfo: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return UserinfoResult{}, fmt.Errorf("userinfo endpoint returned status %d", resp.StatusCode)
	}
	var raw struct {
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return UserinfoResult{}, fmt.Errorf("decode userinfo response: %w", err)
	}
	return UserinfoResult{
		Sub:     raw.Sub,
		Name:    raw.Name,
		Email:   raw.Email,
		Picture: raw.Picture,
	}, nil
}

// Verify parses and validates a JWT string, then checks authorization claim.
func (v *JWTVerifier) Verify(ctx context.Context, tokenString string) (Claims, error) {
	keySet, err := v.jwksCache.Get(ctx, v.jwksURI)
	if err != nil {
		return Claims{}, fmt.Errorf("get JWKS: %w", err)
	}

	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("parse token: %w", err)
	}

	if v.claimKey != "" && v.claimValue != "" {
		if !hasClaimValue(token, v.claimKey, v.claimValue) {
			return Claims{}, fmt.Errorf("missing required claim %q=%q", v.claimKey, v.claimValue)
		}
	}

	return Claims{Sub: token.Subject()}, nil
}

// hasClaimValue checks whether the JWT token contains a claim that matches the expected value.
// The claim can be a string or an array of strings.
func hasClaimValue(token jwt.Token, key, expected string) bool {
	val, ok := token.Get(key)
	if !ok {
		return false
	}
	switch v := val.(type) {
	case string:
		return v == expected
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok && s == expected {
				return true
			}
		}
	case []string:
		for _, s := range v {
			if s == expected {
				return true
			}
		}
	}
	return false
}
