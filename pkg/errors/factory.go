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
func NewInternalError(err error) *ServiceError {
	return wrap(CodeInternalError, http.StatusInternalServerError, MsgInternalError, err)
}

func NewInvalidRequestBodyError(err error) *ServiceError {
	return wrap(CodeInvalidRequestBody, http.StatusBadRequest, MsgInvalidRequestBody, err)
}

func NewValidationError(err error) *ServiceError {
	return wrap(CodeValidationFailed, http.StatusBadRequest, MsgValidationFailed, err)
}

// Car-Specific Errors
func NewCarNotFoundError(err error) *ServiceError {
	return wrap(CodeCarNotFound, http.StatusNotFound, MsgCarNotFound, err)
}
