package errors

import "errors"

// General Errors
var (
	ErrIDIsRequired       = errors.New("id is required")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrEncodingResponse   = errors.New("failed to encode response")
	ErrMethodNotAllowed   = errors.New("Method Not Allowed")
)

// Repository Errors
var (
	ErrCarNotFound = errors.New("car not found")
)
