package repositories

import (
	"cars/models"
	"testing"
)

func TestDefaultCarRepository_List(t *testing.T) {
	// Initial seed data
	cars := map[string]models.Car{
		"1": {
			ID:   "1",
			Make: "Ford", Model: "F10", Package: "Base",
			Color: "Silver", Year: 2010, Category: "Truck",
			Mileage: 120123, Price: 1999900,
		},
		"2": {
			ID:   "2",
			Make: "Toyota", Model: "Camry", Package: "SE",
			Color: "White", Year: 2019, Category: "Sedan",
			Mileage: 3999, Price: 2899000,
		},
		"3": {
			ID:   "3",
			Make: "Toyota", Model: "Rav4", Package: "XSE",
			Color: "Red", Year: 2018, Category: "SUV",
			Mileage: 24001, Price: 2275000,
		},
		"4": {
			ID:   "4",
			Make: "Ford", Model: "Bronco", Package: "Badlands",
			Color: "Burnt Orange", Year: 2022, Category: "SUV",
			Mileage: 0, Price: 4499000,
		},
	}

	repo := NewCarRepository(cars)

	tests := []struct {
		name     string
		filters  models.CarFilters
		expected []string
	}{
		{
			name:     "No filters, returns all",
			filters:  models.CarFilters{},
			expected: []string{"1", "2", "3", "4"},
		},
		{
			name:     "Filter by make",
			filters:  models.CarFilters{Make: "Toyota"},
			expected: []string{"2", "3"},
		},
		{
			name:     "Filter by model",
			filters:  models.CarFilters{Model: "Bronco"},
			expected: []string{"4"},
		},
		{
			name:     "Filter by year",
			filters:  models.CarFilters{Year: 2010},
			expected: []string{"1"},
		},
		{
			name: "Filter by make and model",
			filters: models.CarFilters{
				Make:  "Ford",
				Model: "F10",
			},
			expected: []string{"1"},
		},
		{
			name: "Filter by make and year",
			filters: models.CarFilters{
				Make: "Toyota",
				Year: 2019,
			},
			expected: []string{"2"},
		},
		{
			name: "Filter by model and year",
			filters: models.CarFilters{
				Model: "Rav4",
				Year:  2018,
			},
			expected: []string{"3"},
		},
		{
			name: "Filter by make, model and year",
			filters: models.CarFilters{
				Make:  "Ford",
				Model: "Bronco",
				Year:  2022,
			},
			expected: []string{"4"},
		},
		{
			name: "Case-insensitive filtering",
			filters: models.CarFilters{
				Make:  "tOyOtA",
				Model: "cAmRy",
			},
			expected: []string{"2"},
		},
		{
			name:     "No match, returns empty list",
			filters:  models.CarFilters{Make: "BMW"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(tt.filters)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d cars, got %d", len(tt.expected), len(got))
			}

			// Verify IDs
			for _, id := range tt.expected {
				found := false
				for _, car := range got {
					if car.ID == id {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected car ID %s not found in result", id)
				}
			}
		})
	}
}
