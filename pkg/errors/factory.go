package errors

import "net/http"

func wrap(code string, status int, message string, err error) *ServiceError {
	return &ServiceError{
		Code:       code,
		StatusCode: status,
		Message:    message,
		Err:        err,
	}
}

// General Errors
func NewInvalidRequestBodyError(err error) *ServiceError {
	return wrap(CodeInvalidRequestBody, http.StatusBadRequest, MsgInvalidRequestBody, err)
}

func NewMethodNotAllowedError() *ServiceError {
	return wrap(CodeMethodNotAllowed, http.StatusMethodNotAllowed, MsgMethodNotAllowed, nil)
}

func NewBodyIDMismatchError() *ServiceError {
	return wrap(CodeBodyIDMismatch, http.StatusBadRequest, MsgBodyIDMismatch, nil)
}

func NewIDNotAllowedOnCreateError() *ServiceError {
	return wrap(CodeIDNotAllowedOnCreate, http.StatusBadRequest, MsgIDNotAllowedOnCreate, nil)
}

// Car-Specific Errors
func NewCarNotFoundError(err error) *ServiceError {
	return wrap(CodeCarNotFound, http.StatusNotFound, MsgCarNotFound, err)
}

func NewInvalidCarPathError() *ServiceError {
	return wrap(CodeInvalidCarPath, http.StatusBadRequest, MsgInvalidCarPath, nil)
}
