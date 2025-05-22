package logger

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type ctxKey string

const loggerKey ctxKey = "logger"

// WithLogger adds a logger to the context
func WithLogger(r *http.Request) context.Context {
	prefix := fmt.Sprintf("[method:%s path:%s] ",
		r.Method,
		r.URL.Path,
	)

	logger := log.New(log.Writer(), prefix, log.LstdFlags)
	return context.WithValue(r.Context(), loggerKey, logger)
}

// FromContext retrieves the logger from the context
func FromContext(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerKey).(*log.Logger); ok {
		return logger
	}
	return log.Default()
}
