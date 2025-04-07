package middleware

import (
	"cars/pkg/logger"
	"net/http"
)

// Logging injects a logger into the request context
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogger(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
