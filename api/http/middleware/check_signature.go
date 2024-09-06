package middleware

import (
	"crypto/hmac"
	"encoding/hex"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
	"gitlab.com/martin.kluth1/fserve/signature"
)

type QueryParameters struct {
	Algorithm string `schema:"Fs-Algorithm"`
	Date      string `schema:"Fs-Date"`
	Expires   string `schema:"Fs-Expires"`
	Signature string `schema:"Fs-Signature"`
}

func CheckSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()

		decoder := schema.NewDecoder()

		var p QueryParameters
		err := decoder.Decode(&p, q)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to decode query parameters")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		method := req.Method
		path := req.URL.Path
		if path[0] == '/' {
			path = path[1:]
		}

		h, err := signature.CreateSignature(p.Algorithm, p.Date, p.Expires, method, path)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to create signature")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hash, err := hex.DecodeString(p.Signature)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to decode signature")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !hmac.Equal(hash, h) {
			log.Debug().
				Interface("params", p).
				Str("signature", hex.EncodeToString(h)).
				Msg("Signature mismatch")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}
