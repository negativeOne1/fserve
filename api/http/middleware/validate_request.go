package middleware

import (
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

var ISO8601 = "20060102T150405Z0700"

func checkIfLinkExpired(date, expires string) error {
	d, err := time.Parse(ISO8601, date)
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

func ValidateRequest(v signature.Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			if err := v.IsValid(p.Algorithm, p.Date, p.Expires, method, path, p.Signature); err != nil {
				log.Debug().Err(err).Msg("Failed to validate signature")
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
}
