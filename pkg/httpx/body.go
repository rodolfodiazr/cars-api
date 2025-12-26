package httpx

import (
	e "cars/pkg/errors"

	"encoding/json"
	"net/http"
)

// Decode decodes request body JSON into the provided struct.
func Decode[T any](r *http.Request) (*T, error) {
	if r.Body == nil {
		return nil, e.NewInvalidRequestBodyError(e.ErrEmptyBody)
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req T
	if err := dec.Decode(&req); err != nil {
		return nil, e.NewInvalidRequestBodyError(err)
	}

	return &req, nil
}
