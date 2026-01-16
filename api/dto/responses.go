package dto

// CarResponse represents a car returned to the client.
type CarResponse struct {
	ID string `json:"id"`

	Make     string `json:"make"`
	Model    string `json:"model"`
	Package  string `json:"package,omitempty"`
	Color    string `json:"color"`
	Category string `json:"category"`
	Year     int    `json:"year"`

	Mileage int64 `json:"mileage,omitempty"`
	Price   int64 `json:"price,omitempty"`
}
