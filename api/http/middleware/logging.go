package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		ww := &responseStats{w: w}
		next.ServeHTTP(ww, r)

		log.Info().
			Time("received_time", start).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("agent", r.UserAgent()).
			Dur("latency", time.Since(start)).
			Int("status", ww.code).
			Msg("")
	})
}
