package httpx

import (
	e "cars/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"
)

// DecodeJSON decodes request body JSON into the provided struct.
func DecodeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}

// DecodeAndValidate decodes JSON and calls Validate() if the struct has it.
type Validatable interface {
	Validate() error
}

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
