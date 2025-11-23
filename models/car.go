package models

import (
	"errors"
	"strings"
	"time"
)

// Car represents a vehicle with various attributes such as make, model, package, color, category, year, mileage, and price.
type Car struct {
	ID     string `json:"id"` // Unique identifier for the car.
	BodyID string `json:"-"`

	Make     string `json:"make"`              // Manufacturer of the car (e.g., Toyota, Ford).
	Model    string `json:"model"`             // Model name of the car (e.g., Camry, F-10).
	Package  string `json:"package,omitempty"` // Package level (e.g., SE, XSE).
	Color    string `json:"color"`             // Exterior color of the car.
	Category string `json:"category"`          // Vehicle category (e.g., SUV, Sedan, Truck).
	Year     int    `json:"year"`              // Manufacturing year of the car.

	Mileage float64 `json:"mileage,omitempty"` // Distance the car has traveled, measured in miles.
	Price   float64 `json:"price,omitempty"`   // Price of the car in cents.
}

// Cars represents a collection of Car objects.
type Cars []Car

// Validate checks that all required fields in the Car struct are present and valid.
func (c Car) Validate() error {
	if strings.TrimSpace(c.Make) == "" {
		return errors.New("make is required")
	}

	if strings.TrimSpace(c.Model) == "" {
		return errors.New("model is required")
	}

	if strings.TrimSpace(c.Color) == "" {
		return errors.New("color is required")
	}

	if strings.TrimSpace(c.Category) == "" {
		return errors.New("category is required")
	}

	currentYear := time.Now().Year()
	if c.Year <= 0 || c.Year > currentYear {
		return errors.New("year is not valid")
	}
	return nil
}
