package http

import (
	"net/http"

	"gitlab.com/martin.kluth1/fserve/api/http/middleware"
	"gitlab.com/martin.kluth1/fserve/cache"
	"gitlab.com/martin.kluth1/fserve/signature"
	"gitlab.com/martin.kluth1/fserve/storage"
)

type QueryParameters struct {
	Algorithm string `schema:"Fs-Algorithm"`
	Date      string `schema:"Fs-Date"`
	Expires   string `schema:"Fs-Expires"`
	Signature string `schema:"Fs-Signature"`
}

type Router struct {
	http.ServeMux
	storage storage.Storage
}

func NewRouter(s storage.Storage, c cache.Cache, v signature.Validator) http.Handler {
	r := &Router{
		ServeMux: *http.NewServeMux(),
		storage:  s,
	}

	r.HandleFunc("GET /{resource}", r.handleDownload)
	r.HandleFunc("PUT /{resource}", r.handleUpload)

	m := middleware.CreateChain(
		middleware.Logging,
		middleware.ValidateRequest(v),
		middleware.Caching(c),
	)

	return m(r)
}
