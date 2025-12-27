package dto

// CarUpsertRequest represents the payload for creating/updating a car.
type CarUpsertRequest struct {
	Make     string `json:"make"`
	Model    string `json:"model"`
	Package  string `json:"package"`
	Color    string `json:"color"`
	Category string `json:"category"`
	Year     int    `json:"year"`

	Mileage float64 `json:"mileage"`
	Price   float64 `json:"price"`
}

type CreateCarRequest = CarUpsertRequest
type UpdateCarRequest = CarUpsertRequest
