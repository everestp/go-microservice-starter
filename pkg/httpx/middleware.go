package httpx

import (
	"context"
	"net/http"
)

type requestIDKey struct{}

// RequestID is a simple middleware for injecting an X-Request-ID header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") == "" {
			r.Header.Set("X-Request-ID", "local-request")
		}
		w.Header().Set("X-Request-ID", r.Header.Get("X-Request-ID"))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, r.Header.Get("X-Request-ID"))))
	})
}
