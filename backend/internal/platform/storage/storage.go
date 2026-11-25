package storage

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Presigner produces short-lived upload targets for object storage.
type Presigner interface {
	// PresignPut returns a URL the client can PUT the file to, plus the public URL of the file.
	PresignPut(ctx context.Context, key, contentType string) (uploadURL, publicURL string, err error)
}

// DevPresigner points uploads at the local dev upload route instead of object storage.
type DevPresigner struct {
	BaseURL string
}

func NewDevPresigner(baseURL string) *DevPresigner {
	return &DevPresigner{BaseURL: strings.TrimRight(baseURL, "/")}
}

func (d *DevPresigner) PresignPut(_ context.Context, key, _ string) (string, string, error) {
	public := d.BaseURL + "/uploads/" + key
	return public + "?signature=dev", public, nil
}

// DevUploadHandler stores PUT bodies under dir. It stands in for object storage in development only.
func DevUploadHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/uploads/")
		key = filepath.ToSlash(filepath.Clean("/" + key))
		key = strings.TrimPrefix(key, "/")
		if r.URL.Query().Get("signature") != "dev" || key == "" || strings.Contains(key, "..") {
			http.Error(w, "invalid upload target", http.StatusForbidden)
			return
		}

		path := filepath.Join(dir, filepath.FromSlash(key))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			http.Error(w, "storage unavailable", http.StatusInternalServerError)
			return
		}
		f, err := os.Create(path)
		if err != nil {
			http.Error(w, "storage unavailable", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		if _, err := io.Copy(f, http.MaxBytesReader(w, r.Body, 10<<20)); err != nil {
			http.Error(w, "upload failed", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
