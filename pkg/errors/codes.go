package errors

const (
	// General
	CodeInternalError = "INTERNAL_ERROR"
	MsgInternalError  = "Internal server error"

	CodeInvalidRequestBody = "INVALID_REQUEST_BODY"
	MsgInvalidRequestBody  = "Invalid request body"

	CodeValidationFailed = "VALIDATION_FAILED"
	MsgValidationFailed  = "Validation failed"

	// Car-Specific Errors
	CodeCarNotFound = "CAR_NOT_FOUND"
	MsgCarNotFound  = "Car not found"
)
