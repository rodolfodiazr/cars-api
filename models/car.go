package models

import (
	"errors"
	"strings"
	"time"
)

// Car represents a vehicle with various attributes such as make, model, package, color, category, year, mileage, and price.
type Car struct {
	ID string // Unique identifier for the car.

	Make     string // Manufacturer of the car (e.g., Toyota, Ford).
	Model    string // Model name of the car (e.g., Camry, F-10).
	Package  string // Package level (e.g., SE, XSE).
	Color    string // Exterior color of the car.
	Category string // Vehicle category (e.g., SUV, Sedan, Truck).
	Year     int    // Manufacturing year of the car.

	Mileage float64 // Distance the car has traveled, measured in miles.
	Price   float64 // Price of the car in cents.
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
