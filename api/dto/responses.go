package dto

// CarResponse represents a car returned to the client.
type CarResponse struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}
