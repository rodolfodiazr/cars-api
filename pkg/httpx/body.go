package httpx

import (
	e "cars/pkg/errors"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// DecodeJSON decodes request body JSON into the provided struct.
func DecodeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}

type Validatable interface {
	Validate() error
}

// DecodeAndValidate decodes JSON and then validates the struct.
func DecodeAndValidate[T Validatable](r *http.Request) (*T, error) {
	var obj T
	if err := DecodeJSON(r, &obj); err != nil {
		return nil, e.NewInvalidRequestBodyError(err)
	}

	if err := obj.Validate(); err != nil {
		return nil, e.NewInvalidRequestBodyError(err)
	}

	return &obj, nil
}
