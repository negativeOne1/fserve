package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// CreateChain is a helper function to chain middlewares.
func CreateChain(ms ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			m := ms[i]
			next = m(next)
		}
		return next
	}
}
