package dto

// CreateCarRequest represents the payload for creating a car.
type CreateCarRequest struct {
	Make     string `json:"make"`
	Model    string `json:"model"`
	Package  string `json:"package"`
	Color    string `json:"color"`
	Category string `json:"category"`
	Year     int    `json:"year"`

	Mileage float64 `json:"mileage"`
	Price   float64 `json:"price"`
}

// UpdateCarRequest represents the payload for updating a car.
type UpdateCarRequest struct {
	Make     string `json:"make"`
	Model    string `json:"model"`
	Package  string `json:"package"`
	Color    string `json:"color"`
	Category string `json:"category"`
	Year     int    `json:"year"`

	Mileage float64 `json:"mileage"`
	Price   float64 `json:"price"`
}
