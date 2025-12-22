package httpx

import (
	e "cars/pkg/errors"

	"encoding/json"
	"net/http"
)

// DecodeJSON decodes request body JSON into the provided struct.
func DecodeJSON[T any](r *http.Request) (*T, error) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, e.NewBadRequestError(err)
	}

	return &req, nil
}
