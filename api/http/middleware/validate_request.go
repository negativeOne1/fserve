package middleware

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"time"

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

func ValidateRequest(next http.Handler) http.Handler {
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
		if len(path) > 0 && path[0] == '/' {
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

		if err := checkIfLinkExpired(p.Date, p.Expires); err != nil {
			log.Error().Err(err).Msg("Link expired")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}
