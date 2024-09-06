package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/martin.kluth1/fserve/api/http/middleware"
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

func (s *Router) handleGetRequests(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}

func (s *Router) handlePutRequests(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}
