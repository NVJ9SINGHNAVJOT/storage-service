package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/navjot/storage-service/internal/api"
	"github.com/navjot/storage-service/internal/storage"
	"github.com/navjot/storage-service/pkg"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := pkg.LoadEnv(".env"); err != nil {
		slog.Error("failed to load .env", "error", err)
		os.Exit(1)
	}

	fs := storage.New()
	if err := fs.Init(); err != nil {
		slog.Error("failed to initialize storage directories", "error", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", api.Health)

	mux.HandleFunc("POST /api/media/upload", api.Upload(fs))
	mux.HandleFunc("GET /api/media", api.List(fs))
	mux.HandleFunc("GET /api/media/{id}", api.Metadata(fs))
	mux.HandleFunc("GET /api/media/{id}/file", api.Download(fs))
	mux.HandleFunc("DELETE /api/media/{id}", api.Delete(fs))

	slog.Info("storage service starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
