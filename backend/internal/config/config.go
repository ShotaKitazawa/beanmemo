package config

import "os"

type Config struct {
	DSN             string
	Port            string
	OIDCIssuer      string // OIDC_ISSUER
	OIDCClientID    string // OIDC_CLIENT_ID
	OIDCAudience    string // OIDC_AUDIENCE
	OIDCAllowedSubs string // OIDC_ALLOWED_SUBS (comma-separated)
}

func Load() Config {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "beanmemo.db"
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
		OIDCAllowedSubs: os.Getenv("OIDC_ALLOWED_SUBS"),
	}
}
