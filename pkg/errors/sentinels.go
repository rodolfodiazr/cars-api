package errors

import "errors"

var (
	ErrCarNotFound          = errors.New("car not found")
	ErrInvalidCarPathFormat = errors.New("invalid car path format")
	ErrIDNotAllowedOnCreate = errors.New("id not allowed on create")
)
