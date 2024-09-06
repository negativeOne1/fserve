package http

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"gitlab.com/martin.kluth1/fserve/api/http/middleware"
	"gitlab.com/martin.kluth1/fserve/storage"
)

type QueryParameters struct {
	Algorithm string `schema:"Fs-Algorithm"`
	Date      string `schema:"Fs-Date"`
	Expires   string `schema:"Fs-Expires"`
	Signature string `schema:"Fs-Signature"`
}

type Router struct {
	httprouter.Router
	storage storage.Storage
}

func NewRouter(s storage.Storage) http.Handler {
	r := &Router{
		Router:  *httprouter.New(),
		storage: s,
	}
	r.GET("/:resource", r.handleGetRequests)
	r.PUT("/:resource", r.handlePutRequests)

	c := middleware.CheckSignature(r)
	h := middleware.Logging(c)

	return h
}

func checkIfLinkExpired(date, expires string) error {
	d, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return err
	}

	e, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		return err
	}

	if d.Add(time.Duration(e) * time.Second).Before(time.Now().UTC()) {
		return errors.New("Link expired")
	}

	return nil
}

func (r *Router) handleGetRequests(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	q := req.URL.Query()

	decoder := schema.NewDecoder()
	var p QueryParameters
	if err := decoder.Decode(&p, q); err != nil {
		log.Debug().Err(err).Msg("Failed to decode query parameters")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := checkIfLinkExpired(p.Date, p.Expires); err != nil {
		log.Error().Err(err).Msg("Link expired")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reader, err := r.storage.GetFile(ps.ByName("resource"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get file")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ct := "application/octet-stream"
	ext := filepath.Ext(ps.ByName("resource"))
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

func (s *Router) handlePutRequests(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}
