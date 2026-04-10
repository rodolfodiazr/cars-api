package middleware

import (
	"cars/pkg/contextkeys"
	"cars/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// responseWriter wraps http.ResponseWriter to capture the response
// status code and the total number of bytes written.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

// WriteHeader records the response status code before delegating
// to the underlying ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	if rw.statusCode == 0 {
		rw.statusCode = code
	}
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes the response body, ensuring a default 200 OK status
// is recorded if WriteHeader was not called explicitly.
//
// It also tracks the total number of bytes written.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

// Logging is an HTTP middleware that enriches each request with a unique
// request ID and a request-scoped logger stored in the context.
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
