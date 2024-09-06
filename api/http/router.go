package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
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
	r.GET("/:resource", r.handleDownload)
	r.PUT("/:resource", r.handleUpload)

	c := middleware.ValidateRequest(r)
	h := middleware.Logging(c)

	return h
}
