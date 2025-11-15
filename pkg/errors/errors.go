package errors

import "errors"

// General Errors
var (
	ErrInvalidRequestBody   = errors.New("invalid request body")
	ErrEncodingResponse     = errors.New("failed to encode response")
	ErrMethodNotAllowed     = errors.New("method not allowed")
	ErrBodyIDMismatch       = errors.New("ID in body does not match URL ID")
	ErrIDNotAllowedOnCreate = errors.New("ID must not be provided when creating")
	ErrInvalidCarPathFormat = errors.New("invalid path format: expected /cars/{id}")
)

// Repository Errors
var (
	ErrCarNotFound = errors.New("car not found")
)
