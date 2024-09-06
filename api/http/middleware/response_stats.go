package middleware

import (
	"net/http"
)

type responseStats struct {
	w    http.ResponseWriter
	code int
	buf  []byte
}

func (r *responseStats) Header() http.Header {
	return r.w.Header()
}

func (r *responseStats) WriteHeader(statusCode int) {
	if r.code != 0 {
		return
	}
	r.w.WriteHeader(statusCode)
	r.code = statusCode
}

func (r *responseStats) Write(p []byte) (n int, err error) {
	if r.code == 0 {
		r.WriteHeader(http.StatusOK)
	}

	if r.buf == nil {
		r.buf = append(r.buf, p...)
	}

	return r.w.Write(p)
}
