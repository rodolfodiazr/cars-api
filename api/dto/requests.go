package dto

// CarUpsertRequest represents the payload for creating or fully updating a car.
// When used for updates, ALL fields must be provided.
// Partial updates are not supported.
type CarUpsertRequest struct {
	Make     string `json:"make"`
	Model    string `json:"model"`
	Package  string `json:"package"`
	Color    string `json:"color"`
	Category string `json:"category"`
	Year     int    `json:"year"`

	Mileage float64 `json:"mileage"`
	Price   int64   `json:"price"`
}

type CreateCarRequest = CarUpsertRequest
type UpdateCarRequest = CarUpsertRequest
