package http

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func (r *Router) handleDownload(w http.ResponseWriter, req *http.Request) {
	resource := req.PathValue("resource")

	reader, err := r.storage.GetFile(resource)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get file")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ct := "application/octet-stream"
	ext := filepath.Ext(resource)
	if m := mime.TypeByExtension(ext); m != "" {
		ct = m
	}

	w.Header().Set("Content-Type", ct)
	w.WriteHeader(http.StatusOK)
	written, err := io.Copy(w, reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if written == 0 {
		log.Error().Msg("No bytes written")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
