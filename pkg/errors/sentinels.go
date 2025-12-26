package errors

import "errors"

var (
	ErrEmptyBody   = errors.New("empty body")
	ErrCarNotFound = errors.New("car not found")
)
