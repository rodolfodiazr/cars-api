package dto

import (
	"errors"
	"strings"
	"time"
)

// CreateCarRequest represents the payload for creating a car.
type CreateCarRequest struct {
	Make  string `json:"make" validate:"required"`
	Model string `json:"model" validate:"required"`
	Year  int    `json:"year" validate:"required,min=1886"`
}

func (c CreateCarRequest) Validate() error {
	if strings.TrimSpace(c.Make) == "" {
		return errors.New("make is required")
	}

	if strings.TrimSpace(c.Model) == "" {
		return errors.New("model is required")
	}

	// if strings.TrimSpace(c.Color) == "" {
	// 	return errors.New("color is required")
	// }

	// if strings.TrimSpace(c.Category) == "" {
	// 	return errors.New("category is required")
	// }

	currentYear := time.Now().Year()
	if c.Year <= 0 || c.Year > currentYear {
		return errors.New("year is not valid")
	}
	return nil
}

// UpdateCarRequest represents the payload for updating a car.
type UpdateCarRequest struct {
	ID    string `json:"id,omitempty"` // optional, used only for mismatch validation
	Make  string `json:"make" validate:"required"`
	Model string `json:"model" validate:"required"`
	Year  int    `json:"year" validate:"required,min=1886"`
}
