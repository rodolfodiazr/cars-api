package errors

import "net/http"

// wrap constructs a ServiceError with the provided metadata and underlying error.
//
// It is an internal helper used to standardize error creation across the package.
func wrap(code string, status int, message string, err error) *ServiceError {
	return &ServiceError{
		Code:       code,
		StatusCode: status,
		Message:    message,
		Err:        err,
	}
}

// NewInternalError returns a ServiceError representing an unexpected internal failure.
//
// It should be used when an error cannot be classified or safely exposed to the client.
func NewInternalError(err error) *ServiceError {
	return wrap(CodeInternalError, http.StatusInternalServerError, MsgInternalError, err)
}

// NewInvalidRequestBodyError returns a ServiceError indicating that the request
// body is malformed or cannot be parsed.
func NewInvalidRequestBodyError(err error) *ServiceError {
	return wrap(CodeInvalidRequestBody, http.StatusBadRequest, MsgInvalidRequestBody, err)
}

// NewValidationError returns a ServiceError indicating that request validation failed.
//
// It should be used when input data does not meet required constraints.
func NewValidationError(err error) *ServiceError {
	return wrap(CodeValidationFailed, http.StatusBadRequest, MsgValidationFailed, err)
}

// Car-Specific Errors
func NewCarNotFoundError(err error) *ServiceError {
	return wrap(CodeCarNotFound, http.StatusNotFound, MsgCarNotFound, err)
}
