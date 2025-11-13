package middleware

import (
	"cars/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ctxKey string

const requestIDKey ctxKey = "requestID"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := uuid.NewString()

		ctx := context.WithValue(r.Context(), requestIDKey, id)
		ctx = logger.WithLogger(r.WithContext(ctx))

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r.WithContext(ctx))

		duration := time.Since(start)
		log := logger.FromContext(ctx)
		log.Printf("[req:%s] %d %s %s (%v)", id, rw.statusCode, r.Method, r.URL.Path, duration)
	})
}
