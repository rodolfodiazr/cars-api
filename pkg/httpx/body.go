package httpx

import (
	e "cars/pkg/errors"
	"io"

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

	var payload map[string]any
	if err := dec.Decode(&payload); err != nil {
		return nil, e.NewInvalidRequestBodyError(err)
	}

	if len(payload) == 0 {
		return nil, e.NewInvalidRequestBodyError(e.ErrEmptyBody)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, e.NewInvalidRequestBodyError(err)
	}

	var req T
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, e.NewInvalidRequestBodyError(err)
	}

	if dec.More() {
		return nil, e.NewInvalidRequestBodyError(e.ErrMultipleJSONObjects)
	}

	if err := dec.Decode(new(any)); err != io.EOF {
		return nil, e.NewInvalidRequestBodyError(e.ErrUnexpectedJSONData)
	}

	return &req, nil
}
