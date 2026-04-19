package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/auth"
)

// setupTestServer creates a mock OIDC server and returns a JWTVerifier along with the private key.
func setupTestServer(t *testing.T) (*auth.JWTVerifier, *rsa.PrivateKey, *httptest.Server) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}

	pubKey, err := jwk.FromRaw(privateKey.Public())
	if err != nil {
		t.Fatalf("create JWK: %v", err)
	}
	if err := pubKey.Set(jwk.KeyIDKey, "test-kid"); err != nil {
		t.Fatalf("set key ID: %v", err)
	}
	if err := pubKey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		t.Fatalf("set algorithm: %v", err)
	}

	keySet := jwk.NewSet()
	if err := keySet.AddKey(pubKey); err != nil {
		t.Fatalf("add key to set: %v", err)
	}
	jwksBytes, err := json.Marshal(keySet)
	if err != nil {
		t.Fatalf("marshal JWKS: %v", err)
	}

	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"issuer":            srvURL,
				"jwks_uri":          srvURL + "/.well-known/jwks.json",
				"userinfo_endpoint": srvURL + "/userinfo",
			})
		case "/.well-known/jwks.json":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jwksBytes)
		case "/userinfo":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"sub":     "user|123",
				"name":    "Test User",
				"email":   "test@example.com",
				"picture": "https://example.com/pic.jpg",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	srvURL = srv.URL

	verifier, err := auth.NewJWTVerifier(context.Background(), srv.URL)
	if err != nil {
		srv.Close()
		t.Fatalf("create JWTVerifier: %v", err)
	}

	return verifier, privateKey, srv
}

func buildToken(t *testing.T, privateKey *rsa.PrivateKey, issuer, sub string, extraClaims map[string]any) string {
	t.Helper()

	b := jwt.NewBuilder().
		Issuer(issuer).
		Subject(sub).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(time.Hour))

	token, err := b.Build()
	if err != nil {
		t.Fatalf("build JWT: %v", err)
	}
	for k, v := range extraClaims {
		if err := token.Set(k, v); err != nil {
			t.Fatalf("set claim %q: %v", k, err)
		}
	}

	privKey, err := jwk.FromRaw(privateKey)
	if err != nil {
		t.Fatalf("create private JWK: %v", err)
	}
	if err := privKey.Set(jwk.KeyIDKey, "test-kid"); err != nil {
		t.Fatalf("set key ID: %v", err)
	}
	if err := privKey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		t.Fatalf("set algorithm: %v", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privKey))
	if err != nil {
		t.Fatalf("sign JWT: %v", err)
	}
	return string(signed)
}

func TestJWTVerifier_ValidToken(t *testing.T) {
	verifier, key, srv := setupTestServer(t)
	defer srv.Close()

	tokenStr := buildToken(t, key, srv.URL, "user|123", nil)

	claims, err := verifier.Verify(context.Background(), tokenStr)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if claims.Sub != "user|123" {
		t.Errorf("expected sub 'user|123', got %q", claims.Sub)
	}
}

func TestJWTVerifier_InvalidSignature(t *testing.T) {
	verifier, _, srv := setupTestServer(t)
	defer srv.Close()

	otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tokenStr := buildToken(t, otherKey, srv.URL, "user|456", nil)

	_, err = verifier.Verify(context.Background(), tokenStr)
	if err == nil {
		t.Fatal("expected error for invalid signature")
	}
}

func TestJWTVerifier_FetchUserinfo_Success(t *testing.T) {
	verifier, key, srv := setupTestServer(t)
	defer srv.Close()

	tokenStr := buildToken(t, key, srv.URL, "user|123", nil)

	result, err := verifier.FetchUserinfo(context.Background(), tokenStr)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Sub != "user|123" {
		t.Errorf("expected sub 'user|123', got %q", result.Sub)
	}
	if result.Name != "Test User" {
		t.Errorf("expected name 'Test User', got %q", result.Name)
	}
	if result.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", result.Email)
	}
}

func TestJWTVerifier_ExpiredToken(t *testing.T) {
	verifier, key, srv := setupTestServer(t)
	defer srv.Close()

	privKey, err := jwk.FromRaw(key)
	if err != nil {
		t.Fatalf("create JWK: %v", err)
	}
	if err := privKey.Set(jwk.KeyIDKey, "test-kid"); err != nil {
		t.Fatalf("set key ID: %v", err)
	}
	if err := privKey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		t.Fatalf("set alg: %v", err)
	}

	token, err := jwt.NewBuilder().
		Issuer(srv.URL).
		Subject("user|expired").
		IssuedAt(time.Now().Add(-2 * time.Hour)).
		Expiration(time.Now().Add(-1 * time.Hour)).
		Build()
	if err != nil {
		t.Fatalf("build token: %v", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, privKey))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	_, err = verifier.Verify(context.Background(), string(signed))
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}
