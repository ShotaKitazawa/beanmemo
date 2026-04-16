package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/auth"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/config"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/database"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/handler"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/repository"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/ui"
	"github.com/ShotaKitazawa/beanmemo/backend/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	disableOIDC := flag.Bool("disable-oidc", false, "skip OIDC verification (dev only)")
	flag.Parse()

	cfg := config.Load()

	db, err := database.Connect(context.Background(), cfg.DSN)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	userRepo := repository.NewUserRepository(db)
	recordRepo := repository.NewRecordRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	var verifier *auth.JWTVerifier
	if *disableOIDC {
		slog.Warn("OIDC verification is disabled — do not use in production")
		if err := ensureDefaultUser(db); err != nil {
			return fmt.Errorf("seed default user: %w", err)
		}
	} else {
		if cfg.OIDCIssuerURL == "" {
			return fmt.Errorf("OIDC_ISSUER_URL is required when --disable-oidc is not set")
		}
		verifier, err = auth.NewJWTVerifier(context.Background(), cfg.OIDCIssuerURL, cfg.AuthzClaimKey, cfg.AuthzClaimValue)
		if err != nil {
			return fmt.Errorf("init JWT verifier: %w", err)
		}
	}

	secHandler := handler.NewSecurityHandler(verifier, userRepo, *disableOIDC)

	recordUC := usecase.NewRecordUsecase(recordRepo)
	statsUC := usecase.NewStatsUsecase(statsRepo, recordRepo)

	var userinfoProvider handler.UserinfoProvider
	if *disableOIDC {
		userinfoProvider = handler.NewDisabledOIDCUserinfoProvider()
	} else {
		userinfoProvider = verifier
	}

	h := handler.New(recordUC, statsUC, userinfoProvider)

	srv, err := api.NewServer(h, secHandler)
	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}

	staticFS, err := fs.Sub(ui.Static, "dist")
	if err != nil {
		return fmt.Errorf("create static fs: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", srv))
	mux.Handle("/", spaHandler(http.FS(staticFS)))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("beanmemo backend starting on %s (disableOIDC=%v)", addr, *disableOIDC)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func ensureDefaultUser(db *sql.DB) error {
	ctx := context.Background()
	_, err := db.ExecContext(ctx,
		`INSERT IGNORE INTO users (id, name, email, password_hash) VALUES (1, 'default', 'default@beanmemo.local', 'n/a')`,
	)
	return err
}

func spaHandler(fileSystem http.FileSystem) http.Handler {
	fileServer := http.FileServer(fileSystem)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fileSystem.Open(r.URL.Path)
		if err != nil {
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		defer func() { _ = f.Close() }()
		fileServer.ServeHTTP(w, r)
	})
}
