package errors

import "errors"

var (
	ErrEmptyBody           = errors.New("empty body")
	ErrMultipleJSONObjects = errors.New("multiple JSON objects in request body")
	ErrUnexpectedJSONData  = errors.New("unexpected data after JSON object")

	ErrCarNotFound = errors.New("car not found")
)
