package middleware

import (
	"cars/pkg/contextkeys"
	"cars/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.statusCode == 0 {
		rw.statusCode = code
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := uuid.NewString()

		// 1) Put request ID into context
		ctx := context.WithValue(r.Context(), contextkeys.RequestIDKey, reqID)
		r = r.WithContext(ctx)

		// 2) logger.WithLogger RETURNS context.Context, NOT *http.Request
		ctx = logger.WithLogger(r)

		// 3) Rebuild request with updated context
		r = r.WithContext(ctx)

		// 4) Wrap response writer
		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log := logger.FromContext(ctx)

		log.Printf("%d %s %s (%v)",
			rw.statusCode,
			r.Method,
			r.URL.Path,
			duration,
		)
	})
}
