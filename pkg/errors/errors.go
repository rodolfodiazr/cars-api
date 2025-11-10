package errors

import "errors"

// General Errors
var (
	ErrIDIsRequired       = errors.New("id is required")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrEncodingResponse   = errors.New("failed to encode response")
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrBodyIDMismatch     = errors.New("ID in body does not match URL ID")
)

// Repository Errors
var (
	ErrCarNotFound = errors.New("car not found")
)
