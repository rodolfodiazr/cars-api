package httpx

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, code int, message string) {
	_ = JSON(w, code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}
