package dto

// CarUpsertRequest represents the payload for creating or fully updating a car.
// When used for updates, ALL fields must be provided.
// Partial updates are not supported.
type CarUpsertRequest struct {
	Make     string  `json:"make"`
	Model    string  `json:"model"`
	Package  *string `json:"package,omitempty"`
	Color    string  `json:"color"`
	Category string  `json:"category"`
	Year     int     `json:"year"`

	Mileage *int64 `json:"mileage,omitempty"`
	Price   *int64 `json:"price,omitempty"`
}

type CreateCarRequest = CarUpsertRequest
type UpdateCarRequest = CarUpsertRequest
