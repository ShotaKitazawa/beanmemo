package config

import "os"

type Config struct {
	DSN             string
	Port            string
	OIDCIssuerURL   string // OIDC_ISSUER_URL
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
		OIDCIssuerURL:   os.Getenv("OIDC_ISSUER_URL"),
		AuthzClaimKey:   os.Getenv("AUTHZ_CLAIM_KEY"),
		AuthzClaimValue: os.Getenv("AUTHZ_CLAIM_VALUE"),
	}
}
