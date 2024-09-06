package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"gitlab.com/martin.kluth1/fserve/cache"
)

func Caching(c cache.Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			key := r.URL.Path
			data, found := c.Get(key)
			if found {
				log.Debug().Msgf("Cache hit for key %s %d", key, len(data))
				if _, err := w.Write(data); err != nil {
					log.Error().Err(err).Msg("Failed to write data")
				}
				return
			}

			ww := &responseStats{w: w}
			next.ServeHTTP(ww, r)

			if ww.code == http.StatusOK {
				log.Debug().Msgf("Caching data for key %s %d", key, len(ww.buf))
				if err := c.Set(key, ww.buf); err != nil {
					log.Error().Err(err).Msg("Failed to cache data")
				}
			}
		})
	}
}
