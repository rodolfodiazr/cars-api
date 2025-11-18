package httpx

import (
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
