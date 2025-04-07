package logger

import (
	"context"
	"log"
)

type ctxKey string

const loggerKey ctxKey = "logger"

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, log.Default())
}

// FromContext retrieves the logger from the context
func FromContext(ctx context.Context) *log.Logger {
	if logger, ok := ctx.Value(loggerKey).(*log.Logger); ok {
		return logger
	}
	return log.Default()
}
