package httpx

import (
	e "cars/pkg/errors"

	"encoding/json"
	"net/http"
)

// Decode decodes request body JSON into the provided struct.
func Decode[T any](r *http.Request) (*T, error) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, e.NewInvalidRequestBodyError(err)
	}

	return &req, nil
}
