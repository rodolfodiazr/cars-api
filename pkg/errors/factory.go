package errors

import "net/http"

func wrap(code string, status int, message string, err error) *ServiceError {
	return &ServiceError{
		Code:       code,
		Message:    message,
		StatusCode: status,
		Err:        err,
	}
}

// General Errors
func NewInvalidRequestBodyError(err error) *ServiceError {
	return wrap(CodeInvalidRequestBody, http.StatusBadRequest, MsgInvalidRequestBody, err)
}

func NewEncodingResponseError(err error) *ServiceError {
	return wrap(CodeEncodingResponse, http.StatusInternalServerError, MsgEncodingResponse, err)
}

func NewMethodNotAllowedError(err error) *ServiceError {
	return wrap(CodeMethodNotAllowed, http.StatusMethodNotAllowed, MsgMethodNotAllowed, err)
}

func NewBodyIDMismatchError(err error) *ServiceError {
	return wrap(CodeBodyIDMismatch, http.StatusBadRequest, MsgBodyIDMismatch, err)
}

func NewIDNotAllowedOnCreateError(err error) *ServiceError {
	return wrap(CodeIDNotAllowedOnCreate, http.StatusBadRequest, MsgIDNotAllowedOnCreate, err)
}

func NewInternalServiceError(err error) *ServiceError {
	return wrap(CodeInternalError, http.StatusInternalServerError, MsgInternalError, err)
}

// Car-Specific Errors
func NewCarNotFoundError(err error) *ServiceError {
	return wrap(CodeCarNotFound, http.StatusNotFound, MsgCarNotFound, err)
}

func NewInvalidCarPathError(err error) *ServiceError {
	return wrap(CodeInvalidCarPath, http.StatusBadRequest, MsgInvalidCarPath, err)
}
