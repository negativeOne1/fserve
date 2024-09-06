package http

import (
	"bufio"
	"net/http"
)

func (s *Router) handleUpload(w http.ResponseWriter, req *http.Request) {
	resource := req.PathValue("resource")
	file, _, err := req.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := bufio.NewReader(file)
	defer req.Body.Close()

	if err := s.storage.PutFile(resource, body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
