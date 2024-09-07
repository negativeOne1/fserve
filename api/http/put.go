package http

import (
	"net/http"
)

func (s *Router) handleUpload(w http.ResponseWriter, req *http.Request) {
	resource := req.PathValue("resource")

	file, _, err := req.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := s.storage.Save(resource, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
