package httpx

import (
	e "cars/pkg/errors"
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

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
