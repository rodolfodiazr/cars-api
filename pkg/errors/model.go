package errors

type ServiceError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (e *ServiceError) Details() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}
