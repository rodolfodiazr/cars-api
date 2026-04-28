package errors

// ServiceError represents an application-level error with structured
// metadata suitable for HTTP responses.
//
// It contains a machine-readable Code, a human-readable Message,
// the associated HTTP StatusCode, and an optional underlying error (Err)
// for internal debugging purposes.
type ServiceError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

// Error implements the error interface.
//
// It returns the human-readable message intended for clients.
func (e *ServiceError) Error() string {
	return e.Message
}

func (e *ServiceError) Details() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}
