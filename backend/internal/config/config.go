package config

import "os"

type Config struct {
	DSN             string
	Port            string
	OIDCIssuer      string // OIDC_ISSUER
	OIDCClientID    string // OIDC_CLIENT_ID
	OIDCAudience    string // OIDC_AUDIENCE
	AuthzClaimKey   string // AUTHZ_CLAIM_KEY
	AuthzClaimValue string // AUTHZ_CLAIM_VALUE
}

func Load() Config {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "beanmemo:beanmemo@tcp(localhost:3306)/beanmemo?parseTime=true"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		DSN:             dsn,
		Port:            port,
		OIDCIssuer:      os.Getenv("OIDC_ISSUER"),
		OIDCClientID:    os.Getenv("OIDC_CLIENT_ID"),
		OIDCAudience:    os.Getenv("OIDC_AUDIENCE"),
		AuthzClaimKey:   os.Getenv("AUTHZ_CLAIM_KEY"),
		AuthzClaimValue: os.Getenv("AUTHZ_CLAIM_VALUE"),
	}
}
