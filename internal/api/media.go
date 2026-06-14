package api

import (
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/navjot/storage-service/helper"
	"github.com/navjot/storage-service/internal/models"
	"github.com/navjot/storage-service/internal/storage"
)

func List(fs *storage.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		media, err := fs.List()
		if err != nil {
			slog.Error("failed to list media", "error", err)
			helper.WriteError(w, http.StatusInternalServerError, "failed to list media")
			return
		}
		if media == nil {
			media = []*models.Media{}
		}
		helper.WriteJSON(w, http.StatusOK, media)
	}
}

func Metadata(fs *storage.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		m, err := fs.Get(id)
		if err != nil {
			helper.WriteError(w, http.StatusNotFound, "media not found")
			return
		}
		helper.WriteJSON(w, http.StatusOK, m)
	}
}

func Download(fs *storage.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		dir, err := fs.DirForID(id)
		if err != nil {
			helper.WriteError(w, http.StatusNotFound, "media not found")
			return
		}
		m, err := fs.Get(id)
		if err != nil {
			slog.Error("failed to read metadata for download", "id", id, "error", err)
			helper.WriteError(w, http.StatusInternalServerError, "failed to read metadata")
			return
		}
		http.ServeFile(w, r, filepath.Join(dir, m.Filename))
	}
}
