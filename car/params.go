package car

import "time"

// CarFilters defines optional criteria used to filter cars when listing them.
type CarFilters struct {
	Make  string
	Model string
	Year  *int
}

type CarUpdateParams struct {
	ID string // Unique identifier for the car.

	Make     *string // Manufacturer of the car (e.g., Toyota, Ford).
	Model    *string // Model name of the car (e.g., Camry, F-10).
	Color    *string // Exterior color of the car.
	Category *string // Vehicle category (e.g., SUV, Sedan, Truck).
	Year     *int    // Manufacturing year of the car.

	Package *string // Package level (e.g., SE, XSE).
	Mileage *int64  // Distance the car has traveled, measured in miles.
	Price   *int64  // Price of the car in cents.
}

func (c CarUpdateParams) Validate() error {
	if c.Make != nil && isBlank(*c.Make) {
		return ErrCarMakeRequired
	}

	if c.Model != nil && isBlank(*c.Model) {
		return ErrCarModelRequired
	}

	if c.Color != nil && isBlank(*c.Color) {
		return ErrCarColorRequired
	}

	if c.Category != nil && isBlank(*c.Category) {
		return ErrCarCategoryRequired
	}

	currentYear := time.Now().Year()
	if c.Year != nil && *c.Year < 0 && *c.Year > currentYear {
		return ErrInvalidYear
	}

	if c.Mileage != nil && *c.Mileage < 0 {
		return ErrInvalidMileage
	}

	if c.Price != nil && *c.Price < 0 {
		return ErrInvalidPrice
	}

	return nil
}
