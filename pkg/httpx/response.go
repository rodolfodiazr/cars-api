package httpx

import (
	e "cars/pkg/errors"
	"encoding/json"
	"errors"
	"net/http"
)

// ErrorResponse represents the standard JSON structure returned
// when an API request fails.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// JSON writes a JSON response with the provided HTTP status code.
//
// It automatically sets the Content-Type header to application/json.
//
// For responses that must not include a body, such as 204 No Content
// and 304 Not Modified, the payload is ignored.
func JSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusNoContent || status == http.StatusNotModified {
		return nil
	}

	return json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err *e.ServiceError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    err.Code,
		Message: err.Message,
		Details: err.Details(),
	})
}

func HandleServiceError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var serviceError *e.ServiceError
	if errors.As(err, &serviceError) {
		writeError(w, serviceError)
		return
	}

	writeError(w, e.NewInternalError(err))
}
