package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/ShotaKitazawa/beanmemo/backend/internal/api"
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
	cfg := config.Load()

	db, err := database.Connect(context.Background(), cfg.DSN)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	// Ensure default user exists (Phase 1: single user)
	if err := ensureDefaultUser(db); err != nil {
		return fmt.Errorf("failed to seed default user: %w", err)
	}

	recordRepo := repository.NewRecordRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	recordUC := usecase.NewRecordUsecase(recordRepo)
	statsUC := usecase.NewStatsUsecase(statsRepo, recordRepo)

	h := handler.New(recordUC, statsUC)

	srv, err := api.NewServer(h)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	staticFS, err := fs.Sub(ui.Static, "dist")
	if err != nil {
		return fmt.Errorf("failed to create static fs: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", srv))
	mux.Handle("/", spaHandler(http.FS(staticFS)))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("beanmemo backend starting on %s", addr)
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
