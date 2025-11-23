package errors

const (
	// General
	CodeInvalidRequestBody = "INVALID_REQUEST_BODY"
	MsgInvalidRequestBody  = "Invalid request body"

	CodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	MsgMethodNotAllowed  = "Method not allowed"

	CodeBodyIDMismatch = "BODY_ID_MISMATCH"
	MsgBodyIDMismatch  = "ID in body does not match URL ID"

	CodeIDNotAllowedOnCreate = "ID_NOT_ALLOWED_ON_CREATE"
	MsgIDNotAllowedOnCreate  = "ID must not be provided when creating a new record"

	// Car-Specific Errors
	CodeCarNotFound = "CAR_NOT_FOUND"
	MsgCarNotFound  = "Car not found"

	CodeInvalidCarPath = "INVALID_CAR_PATH"
	MsgInvalidCarPath  = "Invalid path format: expected /cars/{id}"
)
