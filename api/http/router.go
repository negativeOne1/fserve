package http

import (
	"crypto/hmac"
	"encoding/hex"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/martin.kluth1/fserve/api/http/middleware"
	"gitlab.com/martin.kluth1/fserve/signature"
)

const (
	FsAlgorithmHeader = "Fs-Algorithm"
	FsDateHeader      = "Fs-Date"
	FsExpiresHeader   = "Fs-Expires"
	FsSignatureHeader = "Fs-Signature"
)

type Router struct {
	httprouter.Router
}

func NewRouter() http.Handler {
	r := &Router{Router: *httprouter.New()}
	r.GET("/:resource", r.handleGetRequests)
	r.PUT("/:resource", r.handlePutRequests)

	h := middleware.Logging(r)

	return h
}

func (r *Router) handleGetRequests(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	q := req.URL.Query()

	date := q.Get(FsDateHeader)
	expires := q.Get(FsExpiresHeader)
	algo := q.Get(FsAlgorithmHeader)

	sign := q.Get(FsSignatureHeader)
	s, err := hex.DecodeString(sign)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	method := req.Method
	path := req.URL.Path

	h, err := signature.CreateSignature(algo, date, expires, method, path)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !hmac.Equal(s, h) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, err := w.Write([]byte("Hello, World!")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Router) handlePutRequests(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}
