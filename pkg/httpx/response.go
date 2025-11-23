package httpx

import (
	e "cars/pkg/errors"
	"encoding/json"
	"errors"
	"net/http"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Data any `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(successResponse{
		Data: payload,
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Code:    code,
		Message: message,
	})
}

func HandleServiceError(w http.ResponseWriter, err error) {
	var serviceError *e.ServiceError

	if errors.As(err, &serviceError) {
		writeError(w, serviceError.StatusCode, serviceError.Code, serviceError.Message)
		return
	}

	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
}
