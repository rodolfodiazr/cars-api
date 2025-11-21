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
}

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, code int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(Response{
		Data: payload,
	})
}

func errorMessage(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func HandleServiceError(w http.ResponseWriter, err error) {
	var serviceError *e.ServiceError

	if errors.As(err, &serviceError) {
		errorMessage(w, serviceError.StatusCode, serviceError.Code, serviceError.Message)
		return
	}

	errorMessage(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
}
